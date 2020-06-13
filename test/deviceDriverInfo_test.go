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
	"net/http/httptest"
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
	stErrors "github.com/supertokens/supertokens-go/supertokens/errors"
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
		newErr := err.(stErrors.GeneralError)
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
	beforeEach()
	startST("localhost", "8080")

	supertokens.Config("localhost:8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "")
	})
	mux.HandleFunc("/session/verify", func(response http.ResponseWriter, request *http.Request) {
		supertokens.GetSession(response, request, true)
	})
	mux.HandleFunc("/session/refresh", func(response http.ResponseWriter, request *http.Request) {
		supertokens.RefreshSession(response, request)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)

	response := extractInfoFromResponseHeader(res)
	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/verify", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		req.Header.Add("supertokens-sdk-name", "ios")
		req.Header.Add("supertokens-sdk-version", "0.0.0")
		res, _ = client.Do(req)

	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/verify", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/verify", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		req.Header.Add("supertokens-sdk-name", "android")
		req.Header.Add("supertokens-sdk-version", "0.0.1")
		res, _ = client.Do(req)
	}

	var createSessionOutput map[string]interface{}
	createHandler := func(req *http.Request) (*http.Response, error) {
		_ = json.NewDecoder(req.Body).Decode(&createSessionOutput)
		r := ioutil.NopCloser(bytes.NewReader([]byte("")))
		return &http.Response{
			Status: "500",
			Body:   r,
		}, errors.New("custom error")

	}
	mockCreateSession := mockedStruct{
		doFunc: createHandler,
	}

	var verifySessionOutput map[string]interface{}
	verifyHandler := func(req *http.Request) (*http.Response, error) {
		_ = json.NewDecoder(req.Body).Decode(&verifySessionOutput)
		r := ioutil.NopCloser(bytes.NewReader([]byte("")))
		return &http.Response{
			Status: "500",
			Body:   r,
		}, errors.New("custom error")
	}
	mockVerifySession := mockedStruct{
		doFunc: verifyHandler,
	}

	var refreshSessionOutput map[string]interface{}
	refreshHandler := func(req *http.Request) (*http.Response, error) {
		_ = json.NewDecoder(req.Body).Decode(&refreshSessionOutput)
		r := ioutil.NopCloser(bytes.NewReader([]byte("")))
		return &http.Response{
			Status: "500",
			Body:   r,
		}, errors.New("custom error")
	}
	mockRefreshSession := mockedStruct{
		doFunc: refreshHandler,
	}

	var handshakeOutput map[string]interface{}
	handshakeHandler := func(req *http.Request) (*http.Response, error) {
		_ = json.NewDecoder(req.Body).Decode(&handshakeOutput)
		r := ioutil.NopCloser(bytes.NewReader([]byte("")))
		return &http.Response{
			Status: "500",
			Body:   r,
		}, errors.New("custom error")
	}
	mockHandshake := mockedStruct{
		doFunc: handshakeHandler,
	}

	core.AddMockedHTTPHandler("newsession", &mockCreateSession)
	core.AddMockedHTTPHandler("verify", &mockVerifySession)
	core.AddMockedHTTPHandler("refresh", &mockRefreshSession)
	core.AddMockedHTTPHandler("handshake", &mockHandshake)

	{
		core.CreateNewSession("", map[string]interface{}{}, map[string]interface{}{})
		frontendSDK := createSessionOutput["frontendSDK"].([]interface{})

		if len(frontendSDK) != 2 {
			t.Error("incorrect number of frontendSDK")
		}
		{
			data := frontendSDK[0].(map[string]interface{})
			if data["name"] != "ios" {
				t.Error("invalid sdk name")
			}
			if data["version"] != "0.0.0" {
				t.Error("invalid sdk version")
			}

			data = frontendSDK[1].(map[string]interface{})
			if data["name"] != "android" {
				t.Error("invalid sdk name")
			}
			if data["version"] != "0.0.1" {
				t.Error("invalid sdk version")
			}
		}

		driver := createSessionOutput["drive"].(map[string]interface{})
		if driver["name"] != "go" {
			t.Error("invalid driver name")
		}
		if driver["version"] != core.VERSION {
			t.Error("invalid driver version")
		}
	}

	// {
	// 	core.GetSession("", nil, false)
	// 	frontendSDK := verifySessionOutput["frontendSDK"].([]interface{})
	// 	if len(frontendSDK) != 2 {
	// 		t.Error("incorrect number of frontendSDK")
	// 	}

	// 	{
	// 		data := frontendSDK[0].(map[string]interface{})
	// 		if data["name"] != "ios" {
	// 			t.Error("invalid sdk name")
	// 		}
	// 		if data["version"] != "0.0.0" {
	// 			t.Error("invalid sdk version")
	// 		}

	// 		data = frontendSDK[1].(map[string]interface{})
	// 		if data["name"] != "android" {
	// 			t.Error("invalid sdk name")
	// 		}
	// 		if data["version"] != "0.0.1" {
	// 			t.Error("invalid sdk version")
	// 		}
	// 	}

	// 	driver := verifySessionOutput["drive"].(map[string]interface{})
	// 	if driver["name"] != "go" {
	// 		t.Error("invalid driver name")
	// 	}
	// 	if driver["version"] != core.VERSION {
	// 		t.Error("invalid driver version")
	// 	}
	// }

	// {
	// 	core.RefreshSession("")
	// 	frontendSDK := refreshSessionOutput["frontendSDK"].([]interface{})
	// 	if len(frontendSDK) != 2 {
	// 		t.Error("incorrect number of frontendSDK")
	// 	}

	// 	{
	// 		data := frontendSDK[0].(map[string]interface{})
	// 		if data["name"] != "ios" {
	// 			t.Error("invalid sdk name")
	// 		}
	// 		if data["version"] != "0.0.0" {
	// 			t.Error("invalid sdk version")
	// 		}

	// 		data = frontendSDK[1].(map[string]interface{})
	// 		if data["name"] != "android" {
	// 			t.Error("invalid sdk name")
	// 		}
	// 		if data["version"] != "0.0.1" {
	// 			t.Error("invalid sdk version")
	// 		}
	// 	}

	// 	driver := refreshSessionOutput["drive"].(map[string]interface{})
	// 	if driver["name"] != "go" {
	// 		t.Error("invalid driver name")
	// 	}
	// 	if driver["version"] != core.VERSION {
	// 		t.Error("invalid driver version")
	// 	}
	// }

	// {
	// 	core.ResetHandshakeInfo()
	// 	core.GetHandshakeInfoInstance()
	// 	frontendSDK := handshakeOutput["frontendSDK"].([]interface{})
	// 	if len(frontendSDK) != 2 {
	// 		t.Error("incorrect number of frontendSDK")
	// 	}

	// 	{
	// 		data := frontendSDK[0].(map[string]interface{})
	// 		if data["name"] != "ios" {
	// 			t.Error("invalid sdk name")
	// 		}
	// 		if data["version"] != "0.0.0" {
	// 			t.Error("invalid sdk version")
	// 		}

	// 		data = frontendSDK[1].(map[string]interface{})
	// 		if data["name"] != "android" {
	// 			t.Error("invalid sdk name")
	// 		}
	// 		if data["version"] != "0.0.1" {
	// 			t.Error("invalid sdk version")
	// 		}
	// 	}

	// 	driver := handshakeOutput["drive"].(map[string]interface{})
	// 	if driver["name"] != "go" {
	// 		t.Error("invalid driver name")
	// 	}
	// 	if driver["version"] != core.VERSION {
	// 		t.Error("invalid driver version")
	// 	}
	// }
}
