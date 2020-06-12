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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
)

func TestSessionVerifyWithAntiCsrf(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, requeset *http.Request) {
		supertokens.CreateNewSession(response, "id1")

	})
	mux.HandleFunc("/session/verify", func(response http.ResponseWriter, requeset *http.Request) {
		supertokens.GetSession(response, requeset, true)

	})
	mux.HandleFunc("/session/verifyAntiCsrfFalse", func(response http.ResponseWriter, requeset *http.Request) {
		supertokens.GetSession(response, requeset, false)

	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)
	response := extractInfoFromResponseHeader(res)

	req, _ = http.NewRequest("POST", ts.URL+"/session/verify", nil)
	req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
	req.Header.Add("anti-csrf", response["antiCsrf"])
	res, _ = client.Do(req)

	if res.StatusCode != 200 {
		t.Error("response status code was not 200")
	}

	req, _ = http.NewRequest("POST", ts.URL+"/session/verifyAntiCsrfFalse", nil)
	req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
	req.Header.Add("anti-csrf", response["antiCsrf"])
	res, _ = client.Do(req)

	if res.StatusCode != 200 {
		t.Error("response status code was not 200")
	}
}
