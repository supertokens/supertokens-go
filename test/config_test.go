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
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
)

func TestTrailingSlashInHostPath(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts: "http://localhost:8080/",
	})

	response, err := core.CreateNewSession("", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	if response.AccessToken.Token == "" {
		t.Error("accessToken is empty")
	}
	if response.RefreshToken.Token == "" {
		t.Error("refreshToken is empty")
	}
	if response.IDRefreshToken.Token == "" {
		t.Error("idrefreshToken is empty")
	}
	if response.Handle == "" {
		t.Error("handle is empty")
	}
	if response.AntiCsrfToken == nil {
		t.Error("antiCsrfToken is nil")
	}

	_, err = core.GetSession(response.AccessToken.Token, response.AntiCsrfToken, true)
	if err != nil {
		t.Error(err)
	}

	p := core.GetProcessStateInstance()

	if p.GetLastEventByName(core.CallingServiceInVerify) != nil {
		t.Error("processState contains CallingServiceInVerify")
	}

	response2, err := core.RefreshSession(response.RefreshToken.Token, response.AntiCsrfToken)
	if err != nil {
		t.Error(err)
	}

	if response2.AccessToken.Token == "" {
		t.Error("accessToken is empty")
	}
	if response2.RefreshToken.Token == "" {
		t.Error("refreshToken is empty")
	}
	if response2.IDRefreshToken.Token == "" {
		t.Error("idrefreshToken is empty")
	}
	if response2.Handle == "" {
		t.Error("handle is empty")
	}
	if response2.AntiCsrfToken == nil {
		t.Error("antiCsrfToken is nil")
	}

	response3, err := core.GetSession(response2.AccessToken.Token, response2.AntiCsrfToken, true)
	if err != nil {
		t.Error(err)
	}

	if p.GetLastEventByName(core.CallingServiceInVerify) == nil {
		t.Error("processState does not contain CallingServiceInVerify")
	}
	if response3.Handle == "" {
		t.Error("handle is empty")
	}
	if response3.AccessToken.Token == "" {
		t.Error("accessToken is empty")
	}
	if response3.AntiCsrfToken != nil {
		t.Error("antiCsrfToken is not nil")
	}
	if response3.RefreshToken != nil {
		t.Error("refreshToken is not empty")
	}
	if response3.IDRefreshToken != nil {
		t.Error("idrefreshToken is not empty")
	}
	core.ResetProcessState()

	response4, err := core.GetSession(response3.AccessToken.Token, response2.AntiCsrfToken, true)
	if err != nil {
		t.Error(err)
	}
	if core.GetProcessStateInstance().GetLastEventByName(core.CallingServiceInVerify) != nil {
		t.Error("processState contains CallingServiceInVerify")
	}
	if response4.Handle == "" {
		t.Error("handle is empty")
	}
	if response4.AccessToken != nil {
		t.Error("accessToken is not empty")
	}
	if response4.AntiCsrfToken != nil {
		t.Error("antiCsrfToken is not nil")
	}
	if response4.RefreshToken != nil {
		t.Error("refreshToken is not empty")
	}
	if response4.IDRefreshToken != nil {
		t.Error("idrefreshToken is not empty")
	}

	revokeResponse, err := core.RevokeSession(response4.Handle)
	if err != nil {
		t.Error(err)
	}
	if !revokeResponse {
		t.Error("session was not revoked")
	}
}

func TestConfigPathsAreSet(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	cookieSecure := true
	supertokens.Config(supertokens.ConfigMap{
		Hosts:           "http://localhost:8080",
		AccessTokenPath: "/customAccessTokenPath",
		RefreshAPIPath:  "/customRefreshPath",
		CookieDomain:    "customCookieDomain",
		CookieSecure:    &cookieSecure,
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "id1")
	})
	mux.HandleFunc("/customRefreshPath", supertokens.Middleware(
		func(response http.ResponseWriter, request *http.Request) {
			json.NewEncoder(response).Encode(map[string]interface{}{
				"success": true,
			})
		}))

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)
	response := extractInfoFromResponseHeader(res)
	if response["accessTokenPath"] != "/customAccessTokenPath" {
		t.Error("accessToken path not set")
	}

	if response["refreshTokenPath"] != "/customRefreshPath" {
		t.Error("refresh Token path not set")
	}

	if response["accessTokenDomain"] != "customCookieDomain" {
		t.Error("custom domain not set")
	}
	if response["accessTokenSecure"] != strconv.FormatBool(cookieSecure) {
		t.Error("CookieSecure value not set")
	}
	if response["refreshTokenPath"] != "/customRefreshPath" {
		t.Error("refreshTokenPath not set")
	}

	{
		req, _ = http.NewRequest("POST", ts.URL+"/customRefreshPath", nil)
		req.Header.Add("Cookie", "sRefreshToken="+response["refreshToken"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
		res, _ = client.Do(req)
		response2 := extractInfoFromResponseHeader(res)
		if response2["accessToken"] == response["accesToken"] {
			t.Error("refresh did not take place")
		}

	}
}

func TestTrySupertokensHostPath(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config(supertokens.ConfigMap{
		Hosts:          "http://localhost:8080",
		RefreshAPIPath: "/refresh",
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, request *http.Request) {
		supertokens.CreateNewSession(response, "id1")
	})

	mux.HandleFunc("/session/verify", supertokens.Middleware(
		func(response http.ResponseWriter, request *http.Request) {
			supertokens.GetSession(response, request, true)
		}))

	mux.HandleFunc("/refresh", supertokens.Middleware(
		func(response http.ResponseWriter, request *http.Request) {
		}))

	mux.HandleFunc("/session/revoke", supertokens.Middleware(func(response http.ResponseWriter, request *http.Request) {
		session, err := supertokens.GetSession(response, request, true)
		if err != nil {
			return
		}
		session.RevokeSession()
	}))

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
		_, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
	}
	var response2 map[string]string
	{
		req, _ = http.NewRequest("POST", ts.URL+"/refresh", nil)
		req.Header.Add("Cookie", "sRefreshToken="+response["refreshToken"])
		req.Header.Add("anti-csrf", response["antiCsrf"])
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
		_, err := client.Do(req)
		if err != nil {
			t.Error(err)
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
