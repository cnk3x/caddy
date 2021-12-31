// Copyright 2020 Paul Greenberg greenpau@outlook.com
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

package identity

// Handle is the name associated with online services, e.g. Github, Twitter, etc.
type Handle struct {
	Github  string `json:"github,omitempty" xml:"github,omitempty" yaml:"github,omitempty"`
	Twitter string `json:"twitter,omitempty" xml:"twitter,omitempty" yaml:"twitter,omitempty"`
}

// NewHandle returns an instance of Handle
func NewHandle() *Handle {
	return &Handle{}
}
