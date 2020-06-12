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

package testing

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
	spErrors "github.com/supertokens/supertokens-go/supertokens/errors"
)

type mockedStruct struct {
	doFunc DoFunc
}
type DoFunc func(req *http.Request) (*http.Response, error)

func (m *mockedStruct) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestDeviceDriveInfoWithoutFrontendSDK(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")

	var output map[string]interface{}
	handler := func(req *http.Request) (*http.Response, error) {
		_ = json.NewDecoder(req.Body).Decode(&output)
		r := ioutil.NopCloser(bytes.NewReader([]byte("")))
		return &http.Response{
			Status: "500",
			Body:   r,
		}, errors.New("custom error")
	}
	mock := mockedStruct{
		doFunc: handler,
	}
	core.AddMockedHTTPHandler("newsession", &mock)
	supertokens.Config("localhost:8080")
	_, err := core.CreateNewSession("", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		newErr := err.(spErrors.GeneralError)
		if newErr.ActualError.Error() != "custom error" {
			t.Error(newErr)
		}
	}

	drive := output["drive"].(map[string]interface{})
	frontendSDK := output["frontendSDK"].([]interface{})
	if len(frontendSDK) != 0 {
		t.Error("contains frontendSDK values")
	}
	if drive["name"] != "go" && drive["version"] != core.VERSION {
		t.Error("incorrect values set for driver")
	}
}

func TestFrontendSDK(t *testing.T) {

}
