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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const allowedLinkChars = "0123456789abcdefghijklmnopqrstuvwxyz-"

// TableOfContents represent Markdown Table of Contents section.
type TableOfContents struct {
	FilePath  string
	entries   []*tocEntry
	maxDepth  int
	minDepth  int
	lastDepth int
	sep       string
	linkRef   map[string]int
}

type tocEntry struct {
	title string
	link  string
	depth int
}

// NewTableOfContents return a new instance of TableOfContents.
func NewTableOfContents() *TableOfContents {
	return &TableOfContents{
		FilePath: "README.md",
		entries:  []*tocEntry{},
		minDepth: 1000,
		maxDepth: 0,
		sep:      "*",
		linkRef:  make(map[string]int),
	}
}

// AddFilePath adds markdown file path.
func (toc *TableOfContents) AddFilePath(s string) {
	if s == "" {
		return
	}
	toc.FilePath = s
}

// AddHeading adds an entry to TableOfContents.
func (toc *TableOfContents) AddHeading(s string) error {
	if s == "" {
		return fmt.Errorf("cannot add an empty string")
	}
	if !strings.HasPrefix(strings.TrimSpace(s), "#") {
		return fmt.Errorf("heading must start with a pound")
	}
	arr := strings.SplitN(s, " ", 2)
	h := &tocEntry{
		depth: len(arr[0]),
		title: strings.TrimSpace(arr[1]),
	}
	if h.depth > toc.maxDepth {
		toc.maxDepth = h.depth
	}
	if h.depth < toc.minDepth {
		toc.minDepth = h.depth
	}
	depthDiff := h.depth - toc.lastDepth
	if (depthDiff) > 1 && toc.lastDepth > 0 {
		return fmt.Errorf(
			"heading hopped more than one level: %d, %d (current) vs. %d (previous)",
			depthDiff, h.depth, toc.lastDepth,
		)
	}
	toc.lastDepth = h.depth
	toc.entries = append(toc.entries, h)
	return nil
}

func (toc *TableOfContents) getLink(s string) string {
	s = strings.ToLower(s)
	link := "#"
	for _, c := range s {
		if string(c) == " " {
			link += "-"
			continue
		}
		if !strings.Contains(allowedLinkChars, string(c)) {
			continue
		}
		link += string(c)
	}
	i, exists := toc.linkRef[link]
	if exists {
		toc.linkRef[link]++
		link = fmt.Sprintf("%s-%d", link, i)
	} else {
		toc.linkRef[link] = 1
	}

	return link
}

// ToString return string representation of TableOfContents.
func (toc *TableOfContents) ToString() string {
	var tocBuffer bytes.Buffer
	for _, h := range toc.entries {
		offsetDepth := h.depth - toc.minDepth
		tocBuffer.WriteString(strings.Repeat("  ", offsetDepth))
		tocBuffer.WriteString(fmt.Sprintf("%s [%s](%s)", toc.sep, h.title, toc.getLink(h.title)))
		tocBuffer.WriteString("\n")
	}
	return tocBuffer.String()
}

// UpdateToc updates table of contents of the provided file.
func UpdateToc(toc *TableOfContents) error {
	fi, err := os.Stat(toc.FilePath)
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("path %q is not a file", toc.FilePath)
	}

	var fileBuffer bytes.Buffer
	var fileLines []string
	var tocBuffer bytes.Buffer
	var tocBeginMarker = "<!-- begin-markdown-toc -->"
	var tocEndMarker = "<!-- end-markdown-toc -->"
	var isTocOutdated = true
	var isTocFound bool
	var isInsideToc bool
	var tocIndex int
	var firstHeadingIndex int

	fh, err := os.Open(toc.FilePath)
	if err != nil {
		return err
	}
	defer fh.Close()

	// Discovery Scan
	var i int
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		if !isTocFound && firstHeadingIndex == 0 {
			if strings.HasPrefix(line, "##") {
				firstHeadingIndex = i
			}
		}

		if strings.HasPrefix(line, tocEndMarker) {
			isInsideToc = false
			continue
		}
		if strings.HasPrefix(line, tocBeginMarker) {
			isInsideToc = true
			isTocFound = true
			tocIndex = i
			firstHeadingIndex = 0
			continue
		}

		if isInsideToc {
			tocBuffer.WriteString(line + "\n")
			continue
		}

		if !isInsideToc {
			if strings.HasPrefix(line, "##") {
				if firstHeadingIndex == 0 {
					firstHeadingIndex = i
				}
				if err := toc.AddHeading(line); err != nil {
					return fmt.Errorf("toc error: %s", err.Error())
				}
			}
		}

		fileLines = append(fileLines, line)
		i++
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if isTocFound && isInsideToc {
		return fmt.Errorf("toc error: failed to find end marker")
	}

	if !isTocOutdated {
		return nil
	}

	// Found outdated Table of Contents
	if isTocFound {
		fileBuffer.WriteString(strings.Join(fileLines[:tocIndex+1], "\n"))
	} else {
		fileBuffer.WriteString(strings.Join(fileLines[:firstHeadingIndex], "\n"))
		fileBuffer.WriteString("\n")
	}
	fileBuffer.WriteString(tocBeginMarker + "\n")
	fileBuffer.WriteString("## Table of Contents" + "\n\n")
	fileBuffer.WriteString(toc.ToString() + "\n")
	fileBuffer.WriteString(tocEndMarker + "\n")
	if isTocFound {
		fileBuffer.WriteString(strings.Join(fileLines[tocIndex:], "\n") + "\n")
	} else {
		fileBuffer.WriteString("\n")
		fileBuffer.WriteString(strings.Join(fileLines[firstHeadingIndex:], "\n") + "\n")
	}
	mode := fi.Mode()

	return ioutil.WriteFile(toc.FilePath, fileBuffer.Bytes(), mode.Perm())
}
