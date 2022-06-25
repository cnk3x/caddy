// Copyright 2020 Paul Greenberg (greenpau@outlook.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package versioned

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// PackageManager stores metadata about a package.
type PackageManager struct {
	Name          string        `json:"name" xml:"name" yaml:"name"`
	Version       string        `json:"version" xml:"version" yaml:"version"`
	Description   string        `json:"description" xml:"description" yaml:"description"`
	Documentation string        `json:"documentation" xml:"documentation" yaml:"documentation"`
	Git           gitMetadata   `json:"git" xml:"git" yaml:"git"`
	Build         buildMetadata `json:"build" xml:"build" yaml:"build"`
}

// NewPackageManager return an instance of PackageManager.
func NewPackageManager(s string) *PackageManager {
	return &PackageManager{
		Name: s,
	}
}

// gitMetadata stores Git-related metadata.
type gitMetadata struct {
	Branch string `json:"branch" xml:"branch" yaml:"branch"`
	Commit string `json:"commit" xml:"commit" yaml:"commit"`
}

// buildInfo stores build-related metadata.
type buildMetadata struct {
	OperatingSystem string `json:"os" xml:"os" yaml:"os"`
	Architecture    string `json:"arch" xml:"arch" yaml:"arch"`
	User            string `json:"user" xml:"user" yaml:"user"`
	Date            string `json:"date" xml:"date" yaml:"date"`
}

// Banner returns package
func (p *PackageManager) Banner() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s", p.Name, p.Version))
	if p.Git.Branch != "" {
		sb.WriteString(fmt.Sprintf(", branch: %s", p.Git.Branch))
	}
	if p.Git.Commit != "" {
		sb.WriteString(fmt.Sprintf(", commit: %s", p.Git.Commit))
	}
	if p.Build.User != "" && p.Build.Date != "" {
		sb.WriteString(fmt.Sprintf(", build on %s by %s",
			p.Build.Date, p.Build.User,
		))
		if p.Build.OperatingSystem != "" && p.Build.Architecture != "" {
			sb.WriteString(
				fmt.Sprintf(" for %s/%s",
					p.Build.OperatingSystem, p.Build.Architecture,
				))
		}
		sb.WriteString(fmt.Sprintf(
			" (%s/%s %s)",
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
		))
	}
	return sb.String()
}

// ShortBanner returns one-line information about a package.
func (p *PackageManager) ShortBanner() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s", p.Name, p.Version))
	return sb.String()
}

// SetVersion sets Version attribute of PackageManager.
func (p *PackageManager) SetVersion(v, d string) {
	if v != "" {
		p.Version = v
		return
	}
	p.Version = d
}

// SetGitBranch sets Git.Branch attribute of PackageManager.
func (p *PackageManager) SetGitBranch(v, d string) {
	if v != "" {
		p.Git.Branch = v
		return
	}
	p.Git.Branch = d
}

// SetGitCommit sets Git.Commit attribute of PackageManager.
func (p *PackageManager) SetGitCommit(v, d string) {
	if v != "" {
		p.Git.Commit = v
		return
	}
	p.Git.Commit = d
}

// SetBuildUser sets Build.User attribute of PackageManager.
func (p *PackageManager) SetBuildUser(v, d string) {
	if v != "" {
		p.Build.User = v
		return
	}
	p.Build.User = d
}

// SetBuildDate sets Build.Date attribute of PackageManager.
func (p *PackageManager) SetBuildDate(v, d string) {
	if v != "" {
		p.Build.Date = v
		return
	}
	p.Build.Date = d
}

func (p *PackageManager) String() string {
	return p.Banner()
}

// Version represents a software version.
// The version format is `major.minor.patch`.
type Version struct {
	Major    uint64
	Minor    uint64
	Patch    uint64
	FilePath string
	FileName string
	FileType string
	FileDir  string
}

func parseVersion(s string) (uint64, uint64, uint64, error) {
	var major, minor, patch uint64
	var err error
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"")
	s = strings.Trim(s, "'")
	s = strings.TrimSpace(s)
	parts := strings.Split(s, ".")
	if s == "" {
		return major, minor, patch, fmt.Errorf("empty string")
	}
	if len(parts) != 3 {
		return major, minor, patch, fmt.Errorf("version must be in major.minor.patch format, version string: %s", s)
	}
	if major, err = strconv.ParseUint(parts[0], 10, 64); err != nil {
		return major, minor, patch, fmt.Errorf("failed to parse major version, version string: %s", s)
	}
	if minor, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
		return major, minor, patch, fmt.Errorf("failed to parse minor version, version string: %s", s)
	}
	if patch, err = strconv.ParseUint(parts[2], 10, 64); err != nil {
		return major, minor, patch, fmt.Errorf("failed to parse patch version, version string: %s", s)
	}
	return major, minor, patch, nil
}

// NewVersion returns an instance of Version.
func NewVersion(s string) (*Version, error) {
	major, minor, patch, err := parseVersion(s)
	if err != nil {
		return nil, err
	}
	version := &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
	if err := version.SetFile("VERSION"); err != nil {
		return nil, err
	}
	return version, nil
}

// SetFile sets file details of Version.
func (v *Version) SetFile(fp string) error {
	fileDir, fileName := filepath.Split(fp)
	// fileExt := filepath.Ext(fileName)
	v.FilePath = fp
	v.FileDir = fileDir
	v.FileName = fileName
	if fileName == "package.json" {
		v.FileType = "npm-package"
		return nil
	}
	if fileName == "setup.py" {
		v.FileType = "python-package"
		return nil
	}
	v.FileType = "version-file"
	return nil
}

// String returns string representation of Version.
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Bytes returns byte representation of Version string.
func (v *Version) Bytes() []byte {
	return []byte(v.String())
}

// IncrementMajor increments major version
func (v *Version) IncrementMajor(i uint64) {
	v.Major++
	v.Minor = 0
	v.Patch = 0
}

// IncrementMinor increments minor version
func (v *Version) IncrementMinor(i uint64) {
	v.Minor++
	v.Patch = 0
}

// IncrementPatch increments patch version
func (v *Version) IncrementPatch(i uint64) {
	v.Patch++
}

func (v *Version) readVersionFromFile() error {
	var versionFound bool
	switch v.FileType {
	case "version-file", "python-package":
		fh, err := os.Open(v.FilePath)
		if err != nil {
			return err
		}
		defer fh.Close()
		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			line := scanner.Text()
			if v.FileType == "python-package" && !strings.HasPrefix(line, "__version__") {
				continue
			}
			if v.FileType == "python-package" {
				line = strings.SplitN(line, "=", 2)[1]
			}
			major, minor, patch, err := parseVersion(line)
			if err != nil {
				return fmt.Errorf("%s: %s", v.FileType, err)
			}
			v.Major = major
			v.Minor = minor
			v.Patch = patch
			versionFound = true
			break
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		if !versionFound {
			return fmt.Errorf("version not found")
		}
	case "npm-package":
		fc, err := ioutil.ReadFile(v.FilePath)
		if err != nil {
			return err
		}
		var fd map[string]interface{}
		if err := json.Unmarshal(fc, &fd); err != nil {
			return err
		}
		if fd == nil {
			return fmt.Errorf("version not found in %s", v.FilePath)
		}
		versionStr, exists := fd["version"]
		if !exists {
			return fmt.Errorf("version not found in %s", v.FilePath)
		}
		major, minor, patch, err := parseVersion(versionStr.(string))
		if err != nil {
			return err
		}
		v.Major = major
		v.Minor = minor
		v.Patch = patch
		versionFound = true
	default:
		return fmt.Errorf("read error, file type %s is unsupported", v.FileType)
	}

	if !versionFound {
		return fmt.Errorf("version string not found")
	}

	return nil
}

// NewVersionFromFile return Version instance by
// reading VERSION file in a current directory.
func NewVersionFromFile(fp string) (*Version, error) {
	if fp == "" {
		fp = "VERSION"
	}
	version := &Version{}
	if err := version.SetFile(fp); err != nil {
		return nil, err
	}
	if err := version.readVersionFromFile(); err != nil {
		return nil, err
	}
	return version, nil
}

// UpdateFile updates version information in the file associated
// with the version.
func (v *Version) UpdateFile() error {
	fi, err := os.Stat(v.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create version file.
			f, err := os.OpenFile(v.FilePath, os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				return fmt.Errorf("error creating %q file: %v", v.FilePath, err)
			}
			if err := f.Close(); err != nil {
				return fmt.Errorf("error closing %q file: %v", v.FilePath, err)
			}
			fi, err = os.Stat(v.FilePath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("path %s is not a file", v.FilePath)
	}
	mode := fi.Mode()

	switch v.FileType {
	case "version-file":
		return ioutil.WriteFile(v.FilePath, v.Bytes(), mode.Perm())
	case "python-package", "npm-package":
		var buffer bytes.Buffer
		fh, err := os.Open(v.FilePath)
		if err != nil {
			return err
		}
		defer fh.Close()
		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			line := scanner.Text()
			if v.FileType == "python-package" && strings.HasPrefix(line, "__version__") {
				buffer.WriteString("__version__ = '" + v.String() + "'\n")
				continue
			}
			if v.FileType == "npm-package" && strings.Contains(line, "\"version\":") {
				buffer.WriteString("  \"version\": \"" + v.String() + "\",\n")
				continue
			}
			buffer.WriteString(line + "\n")
		}
		return ioutil.WriteFile(v.FilePath, buffer.Bytes(), mode.Perm())
	default:
		return fmt.Errorf("update error, file type %s is unsupported", v.FileType)
	}
}
