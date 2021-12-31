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

package errors

// Public key errors.
const (
	ErrPublicKeyEmptyPayload         StandardError = "public key payload is empty"
	ErrPublicKeyInvalidUsage         StandardError = "public key usage %q is invalid"
	ErrPublicKeyUsagePayloadMismatch StandardError = "public key usage %q does not match its payload"
	ErrPublicKeyBlockType            StandardError = "public key block type %q is invalid"
	ErrPublicKeyParse                StandardError = "public key parse failed: %v"
	ErrPublicKeyUsageUnsupported     StandardError = "public key usage %q is unsupported"
	ErrPublicKeyTypeUnsupported      StandardError = "public key type %q is unsupported"
)
