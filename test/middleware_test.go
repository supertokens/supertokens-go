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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
)

func TestMiddleware(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080/",
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "testing-userID")
	})

	mux.HandleFunc("/user/id", supertokens.Middleware(func(response http.ResponseWriter, request *http.Request) {
		session := supertokens.GetSessionFromRequest(request)

		if session != nil {
			json.NewEncoder(response).Encode(map[string]interface{}{
				"message": session.GetUserID(),
			})
		}
	}))

	mux.HandleFunc("/user/handle", supertokens.Middleware(func(response http.ResponseWriter, request *http.Request) {
		session := supertokens.GetSessionFromRequest(request)
		if session != nil {
			json.NewEncoder(response).Encode(map[string]interface{}{
				"message": session.GetHandle(),
			})
		}
	}))

	mux.HandleFunc("/refresh", supertokens.Middleware(func(response http.ResponseWriter, request *http.Request) {

		json.NewEncoder(response).Encode(map[string]interface{}{
			"message": true,
		})

	}))

	mux.HandleFunc("/logout", supertokens.Middleware(func(response http.ResponseWriter, request *http.Request) {
		session := supertokens.GetSessionFromRequest(request)
		if session != nil {
			err := session.RevokeSession()
			if err != nil {
				return
			}
			json.NewEncoder(response).Encode(map[string]interface{}{
				"message": true,
			})
		}
	}))
	supertokens.OnTryRefreshToken(func(err error, response http.ResponseWriter) {
		response.WriteHeader(401)
		json.NewEncoder(response).Encode(map[string]interface{}{
			"message": "try refresh token",
		})
	})
	supertokens.OnTokenTheftDetected(func(val1 string, val2 string, response http.ResponseWriter) {
		response.WriteHeader(403)
		json.NewEncoder(response).Encode(map[string]interface{}{
			"message": "token theft detected",
		})
	})
	supertokens.OnGeneralError(func(err error, response http.ResponseWriter) {
		response.WriteHeader(400)
		json.NewEncoder(response).Encode(map[string]interface{}{
			"message": "general error",
		})
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)

	response := extractInfoFromResponseHeader(res)
	{
		req, _ = http.NewRequest("POST", ts.URL+"/user/id", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)

		var jsonResponse map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&jsonResponse)
		if err != nil {
			t.Error("error when parsing json body")
		}
		res.Body.Close()
		if jsonResponse["message"] != "testing-userID" {
			t.Error("incorrect response body")
		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/user/handle", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)

		if res.StatusCode != 200 {
			t.Error("response has non 200 status code ")
		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/user/handle", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		res, _ = client.Do(req)

		if res.StatusCode != 401 {
			t.Error("response does not have 401 status code ")
		}
		var jsonResponse2 map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&jsonResponse2)
		res.Body.Close()
		if err != nil {
			t.Error("error when parsing json body")
		}
		if jsonResponse2["message"] != "try refresh token" {
			t.Error("incorrect response body")
		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/user/handle", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)

		if res.StatusCode != 440 {
			t.Error("incorrect status code")
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Error("error when parsing body")
		}
		if !strings.Contains(string(body), "Unauthorized") {
			t.Error("incorrect response")
		}
	}

	req, _ = http.NewRequest("POST", ts.URL+"/refresh", nil)
	req.Header.Add("Cookie", "sRefreshToken="+response["refreshToken"])
	res, _ = client.Do(req)

	response2 := extractInfoFromResponseHeader(res)
	var response3 map[string]string
	{
		req, _ = http.NewRequest("POST", ts.URL+"/user/id", nil)
		req.Header.Add("Cookie", "sAccessToken="+response2["accessToken"]+";sIdRefreshToken="+response2["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response2["antiCsrf"])
		res, _ = client.Do(req)

		var jsonResponse3 map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&jsonResponse3)
		if err != nil {
			t.Error("error when parsing body")
		}
		res.Body.Close()

		if jsonResponse3["message"] != "testing-userID" {
			t.Error("incorrect json response")
		}

		response3 = extractInfoFromResponseHeader(res)
	}
	{
		req, _ = http.NewRequest("POST", ts.URL+"/user/handle", nil)
		req.Header.Add("Cookie", "sAccessToken="+response3["accessToken"]+";sIdRefreshToken="+response2["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response2["antiCsrf"])
		res, _ = client.Do(req)

		if res.StatusCode != 200 {
			t.Error("non 200 status code")

		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/refresh", nil)
		req.Header.Add("Cookie", "sRefreshToken="+response["refreshToken"])
		res, _ = client.Do(req)

		if res.StatusCode != 403 {
			t.Error("incorrect status code")
		}

		var jsonResponse4 map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&jsonResponse4)
		if err != nil {
			t.Error("error when parsing body")
		}
		res.Body.Close()

		if jsonResponse4["message"] != "token theft detected" {
			t.Error("incorrect response body")
		}
	}
	var response4 map[string]string
	{
		req, _ = http.NewRequest("POST", ts.URL+"/logout", nil)
		req.Header.Add("Cookie", "sAccessToken="+response3["accessToken"]+";sIdRefreshToken="+response2["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response2["antiCsrf"])
		res, _ = client.Do(req)

		if res.StatusCode != 200 {
			t.Error("non 200 response code")
		}

		response4 = extractInfoFromResponseHeader(res)
		if response4["antiCsrf"] != "" {
			t.Error("antiCsrf is not empty")
		}
		if response4["accessToken"] != "" {
			t.Error("accessToken is not empty")
		}
		if response4["refreshToken"] != "" {
			t.Error("refreshToken is not empty")
		}
		if response4["idRefreshTokenFromHeader"] != "remove" {
			t.Error("incorrect value")
		}
		if response4["idRefreshTokenFromCookie"] != "" {
			t.Error("idRefreshTokenFromCookie is not empty")
		}
		if response4["accessTokenExpiry"] != "Thu, 01 Jan 1970 00:00:00 GMT" {
			t.Error("incorrect value")
		}
		if response4["idRefreshTokenExpiry"] != "Thu, 01 Jan 1970 00:00:00 GMT" {
			t.Error("incorrect value")
		}
		if response4["refreshTokenExpiry"] != "Thu, 01 Jan 1970 00:00:00 GMT" {
			t.Error("incorrect value")
		}
	}

	req, _ = http.NewRequest("POST", ts.URL+"/user/handle", nil)
	req.Header.Add("Cookie", "sAccessToken="+response4["accessToken"]+";sIdRefreshToken="+response4["idRefreshTokenFromCookie"])
	res, _ = client.Do(req)
	if res.StatusCode != 401 {
		t.Error("incorrect status code")
	}
}
