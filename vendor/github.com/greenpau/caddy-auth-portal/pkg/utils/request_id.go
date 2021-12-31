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

package utils

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// GetContentType returns requested content type.
func GetContentType(r *http.Request) string {
	ct := r.Header.Get("Accept")
	if ct == "" {
		ct = "text/html"
	}
	return ct
}

// GetRequestID returns request ID.
func GetRequestID(r *http.Request) string {
	return uuid.NewV4().String()
}
