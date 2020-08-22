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
	"strconv"
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
	"github.com/supertokens/supertokens-go/supertokens/errors"
)

func TestTokenTheftDetected(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080",
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "id1")

	})
	mux.HandleFunc("/session/verify", func(response http.ResponseWriter, request *http.Request) {
		supertokens.GetSession(response, request, true)

	})

	mux.HandleFunc("/session/refresh", func(response http.ResponseWriter, request *http.Request) {
		_, err := supertokens.RefreshSession(response, request)
		if err != nil {
			if errors.IsTokenTheftDetectedError(err) {
				json.NewEncoder(response).Encode(map[string]interface{}{
					"success": true,
				})
				return
			}
		}
		json.NewEncoder(response).Encode(map[string]interface{}{
			"success": false,
		})

	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)
	response := extractInfoFromResponseHeader(res)

	var response2 map[string]string
	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/refresh", nil)
		req.Header.Add("Cookie", "sRefreshToken="+response["refreshToken"])
		res, _ = client.Do(req)
		response2 = extractInfoFromResponseHeader(res)
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/verify", nil)
		req.Header.Add("Cookie", "sAccessToken="+response2["accessToken"]+";sIdRefreshToken="+response2["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response2["antiCsrf"])
		res, _ = client.Do(req)
	}

	var response3 map[string]string
	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/refresh", nil)
		req.Header.Add("Cookie", "sRefreshToken="+response["refreshToken"])
		res, _ = client.Do(req)
		response3 = extractInfoFromResponseHeader(res)

		var jsonResponse map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&jsonResponse)
		if err != nil {
			t.Error("error when parsing body")
		}
		res.Body.Close()
		if !jsonResponse["success"].(bool) {
			t.Error("incorrect json response")
		}
	}

	{
		version, _ := core.GetQuerierInstance().GetAPIVersion()
		if core.MaxVersion(version, "2.1") == "2.1" {
			if (response3["accessTokenDomain"] != "localhost" &&
				response3["accessTokenDomain"] != "supertokens.io") ||
				(response3["refreshTokenDomain"] != "localhost" &&
					response3["refreshTokenDomain"] != "supertokens.io") ||
				(response3["idRefreshTokenDomain"] != "localhost" &&
					response3["idRefreshTokenDomain"] != "supertokens.io") {
				t.Error("incorrect cookie domain")
			}
		} else {
			if response3["accessTokenDomain"] != "" ||
				response3["refreshTokenDomain"] != "" ||
				response3["idRefreshTokenDomain"] != "" {
				t.Error("incorrect cookie domain")
			}
		}
		if response3["antiCsrf"] != "" {
			t.Error("antiCsrf is not empty")
		}
		if response3["accessToken"] != "" {
			t.Error("accessToken is not empty")
		}
		if response3["refreshToken"] != "" {
			t.Error("refreshToken is not empty")
		}
		if response3["idRefreshTokenFromHeader"] != "remove" {
			t.Error("incorrect value")
		}
		if response3["idRefreshTokenFromCookie"] != "" {
			t.Error("idRefreshTokenFromCookie is not empty")
		}
		if response3["accessTokenExpiry"] != "Thu, 01 Jan 1970 00:00:00 GMT" {
			t.Error("incorrect value")
		}
		if response3["idRefreshTokenExpiry"] != "Thu, 01 Jan 1970 00:00:00 GMT" {
			t.Error("incorrect value")
		}
		if response3["refreshTokenExpiry"] != "Thu, 01 Jan 1970 00:00:00 GMT" {
			t.Error("incorrect value")
		}
	}

}

func TestBasicUsage(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080",
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "id1")

	})
	mux.HandleFunc("/session/verify", func(response http.ResponseWriter, request *http.Request) {
		supertokens.GetSession(response, request, true)

	})

	mux.HandleFunc("/session/refresh", func(response http.ResponseWriter, request *http.Request) {
		supertokens.RefreshSession(response, request)

	})

	mux.HandleFunc("/session/revoke", func(response http.ResponseWriter, request *http.Request) {
		session, err := supertokens.GetSession(response, request, true)
		if err != nil {
			return
		}
		session.RevokeSession()
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)

	response := extractInfoFromResponseHeader(res)

	if response["antiCsrf"] == "" {
		t.Error("antiCsrf is empty")
	}
	if response["accessToken"] == "" {
		t.Error("accessToken is empty")
	}
	if response["refreshToken"] == "" {
		t.Error("refreshToken is empty")
	}
	if response["idRefreshTokenFromHeader"] == "" {
		t.Error("idRefreshTokenFromHeader is empty")
	}
	if response["idRefreshTokenFromCookie"] == "" {
		t.Error("idRefreshTokenFromCookie is empty")
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/verify", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
		if core.GetProcessStateInstance().GetLastEventByName(core.CallingServiceInVerify) != nil {
			t.Error("processState contains CallingServiceInVerify")
		}
	}

	var response2 map[string]string
	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/refresh", nil)
		req.Header.Add("Cookie", "sRefreshToken="+response["refreshToken"])
		res, _ = client.Do(req)
		response2 = extractInfoFromResponseHeader(res)
	}

	if response2["antiCsrf"] == "" {
		t.Error("antiCsrf is empty")
	}
	if response2["accessToken"] == "" {
		t.Error("accessToken is empty")
	}
	if response2["refreshToken"] == "" {
		t.Error("refreshToken is empty")
	}
	if response2["idRefreshTokenFromHeader"] == "" {
		t.Error("idRefreshTokenFromHeader is empty")
	}
	if response2["idRefreshTokenFromCookie"] == "" {
		t.Error("idRefreshTokenFromCookie is empty")
	}

	var response3 map[string]string
	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/verify", nil)
		req.Header.Add("Cookie", "sAccessToken="+response2["accessToken"]+";sIdRefreshToken="+response2["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response2["antiCsrf"])
		res, _ = client.Do(req)
		if core.GetProcessStateInstance().GetLastEventByName(core.CallingServiceInVerify) == nil {
			t.Error("processState does not contain CallingServiceInVerify")
		}
		response3 = extractInfoFromResponseHeader(res)
	}

	core.ResetProcessState()

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/verify", nil)
		req.Header.Add("Cookie", "sAccessToken="+response3["accessToken"]+";sIdRefreshToken="+response3["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response2["antiCsrf"])
		res, _ = client.Do(req)
		if core.GetProcessStateInstance().GetLastEventByName(core.CallingServiceInVerify) != nil {
			t.Error("processState contains CallingServiceInVerify")
		}
	}

	var response4 map[string]string
	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/revoke", nil)
		req.Header.Add("Cookie", "sAccessToken="+response3["accessToken"]+";sIdRefreshToken="+response3["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response2["antiCsrf"])
		res, _ = client.Do(req)

		response4 = extractInfoFromResponseHeader(res)
	}

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

func TestSessionVerifyWithAntiCsrf(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080",
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "id1")
	})
	mux.HandleFunc("/session/verify", func(response http.ResponseWriter, request *http.Request) {
		supertokens.GetSession(response, request, true)
	})
	mux.HandleFunc("/session/verifyAntiCsrfFalse", func(response http.ResponseWriter, request *http.Request) {
		supertokens.GetSession(response, request, false)
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
		res, _ = client.Do(req)

		if res.StatusCode != 200 {
			t.Error("response status code was not 200")
		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/verifyAntiCsrfFalse", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)

		if res.StatusCode != 200 {
			t.Error("response status code was not 200")
		}
	}
}

func TestSessionVerifyWithoutAntiCsrf(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080",
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "id1")
	})
	mux.HandleFunc("/session/verify", func(response http.ResponseWriter, request *http.Request) {
		supertokens.GetSession(response, request, false)
	})
	mux.HandleFunc("/session/verifyAntiCsrfFalse", func(response http.ResponseWriter, request *http.Request) {
		supertokens.GetSession(response, request, false)
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
		res, _ = client.Do(req)
		if res.StatusCode != 200 {
			t.Error("response status code was not 200")
		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/verifyAntiCsrfFalse", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		res, _ = client.Do(req)
		if res.StatusCode != 200 {
			t.Error("response status code was not 200")
		}
	}

}

func TestSessionRevoking(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080",
	})
	supertokens.RevokeAllSessionsForUser("someUniqueUserId")

	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "")
	})

	mux.HandleFunc("/usercreate", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "someUniqueUserId")
	})
	mux.HandleFunc("/session/revoke", func(response http.ResponseWriter, request *http.Request) {
		session, err := supertokens.GetSession(response, request, true)
		if err != nil {
			return
		}
		session.RevokeSession()
	})
	mux.HandleFunc("/session/revokeUserId", func(response http.ResponseWriter, request *http.Request) {
		session, err := supertokens.GetSession(response, request, true)
		if err != nil {
			return
		}
		supertokens.RevokeAllSessionsForUser(session.GetUserID())
	})
	mux.HandleFunc("/session/getSessionWithUserId1", func(response http.ResponseWriter, request *http.Request) {
		sessionHandles, err := supertokens.GetAllSessionHandlesForUser("someUniqueUserId")
		if err != nil {
			return
		}
		response.Write([]byte(strconv.Itoa(len(sessionHandles))))
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)
	response := extractInfoFromResponseHeader(res)

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/revoke", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)

		sessionRevokedResponse := extractInfoFromResponseHeader(res)

		if sessionRevokedResponse["antiCsrf"] != "" {
			t.Error("antiCsrf is not empty")
		}
		if sessionRevokedResponse["accessToken"] != "" {
			t.Error("accessToken is not empty")
		}
		if sessionRevokedResponse["refreshToken"] != "" {
			t.Error("refreshToken is not empty")
		}
		if sessionRevokedResponse["idRefreshTokenFromHeader"] != "remove" {
			t.Error("incorrect value")
		}
		if sessionRevokedResponse["idRefreshTokenFromCookie"] != "" {
			t.Error("idRefreshTokenFromCookie is not empty")
		}
		if sessionRevokedResponse["accessTokenExpiry"] != "Thu, 01 Jan 1970 00:00:00 GMT" {
			t.Error("incorrect value")
		}
		if sessionRevokedResponse["idRefreshTokenExpiry"] != "Thu, 01 Jan 1970 00:00:00 GMT" {
			t.Error("incorrect value")
		}
		if sessionRevokedResponse["refreshTokenExpiry"] != "Thu, 01 Jan 1970 00:00:00 GMT" {
			t.Error("incorrect value")
		}
	}

	req, _ = http.NewRequest("POST", ts.URL+"/create", nil)
	client.Do(req)

	req, _ = http.NewRequest("POST", ts.URL+"/usercreate", nil)
	res, _ = client.Do(req)
	userCreateResponse := extractInfoFromResponseHeader(res)

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/revokeUserId", nil)
		req.Header.Add("Cookie", "sAccessToken="+userCreateResponse["accessToken"]+";sIdRefreshToken="+userCreateResponse["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", userCreateResponse["antiCsrf"])
		res, _ = client.Do(req)
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/getSessionWithUserId1", nil)
		res, _ = client.Do(req)

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		res.Body.Close()
		if string(body) != "0" {
			t.Error("sessions were not revoked")
		}
	}
}

func TestManipulatingSessionData(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080",
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "")
	})
	mux.HandleFunc("/updateSessionData", func(response http.ResponseWriter, request *http.Request) {
		session, err := supertokens.GetSession(response, request, true)
		if err != nil {
			return
		}
		session.UpdateSessionData(map[string]interface{}{
			"key": "value",
		})

	})
	mux.HandleFunc("/getSessionData", func(response http.ResponseWriter, request *http.Request) {
		session, err := supertokens.GetSession(response, request, true)
		if err != nil {
			return
		}
		data, err := session.GetSessionData()
		if err != nil {
			return
		}
		response.Write([]byte(data["key"].(string)))
	})
	mux.HandleFunc("/updateSessionData2", func(response http.ResponseWriter, request *http.Request) {
		session, err := supertokens.GetSession(response, request, true)
		if err != nil {
			return
		}
		session.UpdateSessionData(map[string]interface{}{
			"key": "value2",
		})

	})
	mux.HandleFunc("/updateSessionDataInvalidSessionHandle", func(response http.ResponseWriter, request *http.Request) {

		err := supertokens.UpdateSessionData("InvalidHandle", map[string]interface{}{
			"key": "value2",
		})
		if err != nil {
			if errors.IsUnauthorizedError(err) {
				json.NewEncoder(response).Encode(map[string]interface{}{
					"success": true,
				})
				return
			}
		}
		json.NewEncoder(response).Encode(map[string]interface{}{
			"success": false,
		})
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)
	response := extractInfoFromResponseHeader(res)
	{
		req, _ = http.NewRequest("POST", ts.URL+"/updateSessionData", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/getSessionData", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		res.Body.Close()
		if string(body) != "value" {
			t.Error("incorrect response")
		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/updateSessionData2", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/getSessionData", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		res.Body.Close()
		if string(body) != "value2" {
			t.Error("incorrect response")
		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/updateSessionDataInvalidSessionHandle", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)

		var jsonResponse map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&jsonResponse)
		if err != nil {
			t.Error(err)
		}
		res.Body.Close()
		if !jsonResponse["success"].(bool) {
			t.Error("invalid response")
		}
	}
}

func TestManipulatingJWTData(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080",
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "")
	})
	mux.HandleFunc("/updateJWTPayload", func(response http.ResponseWriter, request *http.Request) {
		session, err := supertokens.GetSession(response, request, true)
		accessTokenBefore := session.GetAccessToken()
		if err != nil {
			return
		}

		session.UpdateJWTPayload(map[string]interface{}{
			"key": "value",
		})
		if session.GetAccessToken() == accessTokenBefore {
			response.WriteHeader(500)
			return
		}

	})
	mux.HandleFunc("/getJWTPayload", func(response http.ResponseWriter, request *http.Request) {
		session, err := supertokens.GetSession(response, request, true)
		if err != nil {
			return
		}
		data := session.GetJWTPayload()

		response.Write([]byte(data["key"].(string)))
	})
	mux.HandleFunc("/updateJWTPayload2", func(response http.ResponseWriter, request *http.Request) {
		session, err := supertokens.GetSession(response, request, true)
		if err != nil {
			return
		}
		session.UpdateJWTPayload(map[string]interface{}{
			"key": "value2",
		})
	})
	mux.HandleFunc("/updateJWTPayloadInvalidSessionHandle", func(response http.ResponseWriter, request *http.Request) {
		err := supertokens.UpdateJWTPayload("InvalidHandle", map[string]interface{}{
			"key": "value3",
		})
		if err != nil {
			if errors.IsUnauthorizedError(err) {
				json.NewEncoder(response).Encode(map[string]interface{}{
					"success": true,
				})
				return
			}
		}
		json.NewEncoder(response).Encode(map[string]interface{}{
			"success": false,
		})
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)
	response := extractInfoFromResponseHeader(res)

	var response2 map[string]string
	{
		req, _ = http.NewRequest("POST", ts.URL+"/updateJWTPayload", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
		response2 = extractInfoFromResponseHeader(res)
	}
	{
		req, _ = http.NewRequest("POST", ts.URL+"/getJWTPayload", nil)
		req.Header.Add("Cookie", "sAccessToken="+response2["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		res.Body.Close()
		if string(body) != "value" {
			t.Error("incorrect response")
		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/updateJWTPayload2", nil)
		req.Header.Add("Cookie", "sAccessToken="+response2["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
		response2 = extractInfoFromResponseHeader(res)
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/getJWTPayload", nil)
		req.Header.Add("Cookie", "sAccessToken="+response2["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		res.Body.Close()
		if string(body) != "value2" {
			t.Error("incorrect response")
		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/updateJWTPayloadInvalidSessionHandle", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)

		var jsonResponse map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&jsonResponse)
		if err != nil {
			t.Error(err)
		}
		res.Body.Close()
		if !jsonResponse["success"].(bool) {
			t.Error("invalid response")
		}
	}
}

func TestAppendingToExistingHeader(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080",
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Add("testHeader", "testValue")
		response.Header().Add("Access-Control-Expose-Headers", "customValue")
		supertokens.CreateNewSession(response, "")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)

	if res.Header.Get("testHeader") != "testValue" {
		t.Error("header value not set")
	}
	if res.Header.Get("Access-Control-Expose-Headers") != "customValue, id-refresh-token, anti-csrf" {
		t.Error("header value not set")
	}

	response := extractInfoFromResponseHeader(res)

	if response["antiCsrf"] == "" {
		t.Error("antiCsrf is empty")
	}
	if response["accessToken"] == "" {
		t.Error("accessToken is empty")
	}
	if response["refreshToken"] == "" {
		t.Error("refreshToken is empty")
	}
	if response["idRefreshTokenFromHeader"] == "" {
		t.Error("idRefreshTokenFromHeader is empty")
	}
	if response["idRefreshTokenFromCookie"] == "" {
		t.Error("idRefreshTokenFromCookie is empty")
	}
}

func TestAntiCsrfDisabled(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080",
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "")
	})

	mux.HandleFunc("/session/verify", func(response http.ResponseWriter, request *http.Request) {
		supertokens.GetSession(response, request, true)
	})

	mux.HandleFunc("/session/verifyAntiCsrfFalse", func(response http.ResponseWriter, request *http.Request) {
		supertokens.GetSession(response, request, false)
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
		res, _ = client.Do(req)
		if res.StatusCode != 200 {
			t.Error("non 200 status code")
		}
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/session/verifyAntiCsrfFalse", nil)
		req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
		res, _ = client.Do(req)
		if res.StatusCode != 200 {
			t.Error("non 200 status code")
		}
	}
}
