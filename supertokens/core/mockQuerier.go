/*
 * Copyright (c) 2020, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package core

import (
	"flag"
	"net/http"
)

// MockedHTTPClient mocked http client
type MockedHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var idMap = map[string]MockedHTTPClient{}

// AddMockedHTTPHandler used during testing
func AddMockedHTTPHandler(requestID string, handler MockedHTTPClient) {
	if flag.Lookup("test.v") != nil {
		idMap[requestID] = handler
	}
}

// GetMockedHTTPClient is used during testing
func GetMockedHTTPClient(requestID string) MockedHTTPClient {
	if flag.Lookup("test.v") == nil {
		return nil
	}
	value := idMap[requestID]
	if value == nil {
		return nil
	}
	return value
}

// ResetHTTPMocking sets idMap to an empty map
func ResetHTTPMocking() {
	idMap = map[string]MockedHTTPClient{}
}
