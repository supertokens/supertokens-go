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

package supertokens

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/supertokens/supertokens-go/supertokens/core"
)

const accessTokenCookieKey = "sAccessToken"
const refreshTokenCookieKey = "sRefreshToken"

const idRefreshTokenCookieKey = "sIdRefreshToken"
const idRefreshTokenHeaderKey = "id-refresh-token"

const antiCsrfHeaderKey = "anti-csrf"
const frontendSDKNameHeaderKey = "supertokens-sdk-name"
const frontendSDKVersionHeaderKey = "supertokens-sdk-version"

var configMap *ConfigMap = nil

func configCookieAndHeaders(config ConfigMap) {
	configMap = &config
}

func attachAccessTokenToCookie(response http.ResponseWriter, token string,
	expiry uint64, domain string, secure bool, path string, sameSite string) {
	setCookie(response, accessTokenCookieKey, token, domain, secure, true, expiry, path, sameSite)
}

func attachRefreshTokenToCookie(response http.ResponseWriter, token string,
	expiry uint64, domain string, secure bool, path string, sameSite string) {
	setCookie(response, refreshTokenCookieKey, token, domain, secure, true, expiry, path, sameSite)
}

func setIDRefreshTokenInHeaderAndCookie(response http.ResponseWriter, token string,
	expiry uint64, domain string, secure bool, path string, sameSite string) {
	setHeader(response, idRefreshTokenHeaderKey, token+";"+fmt.Sprint(expiry))
	setHeader(response, "Access-Control-Expose-Headers", idRefreshTokenHeaderKey)

	setCookie(response, idRefreshTokenCookieKey, token, domain, secure, true, expiry, path, sameSite)
}

func setAntiCsrfTokenInHeaders(response http.ResponseWriter, antiCsrfToken string) {
	setHeader(response, antiCsrfHeaderKey, antiCsrfToken)
	setHeader(response, "Access-Control-Expose-Headers", antiCsrfHeaderKey)
}

func saveFrontendInfoFromRequest(request *http.Request) {
	name := getHeader(request, frontendSDKNameHeaderKey)
	version := getHeader(request, frontendSDKVersionHeaderKey)
	if name != nil && version != nil {
		core.GetDeviceInfoInstance().AddToFrontendSDKs(*name, *version)
	}
}

func getAccessTokenFromCookie(request *http.Request) *string {
	return getCookieValue(request, accessTokenCookieKey)
}

func getAntiCsrfTokenFromHeaders(request *http.Request) *string {
	return getHeader(request, antiCsrfHeaderKey)
}

func getIDRefreshTokenFromCookie(request *http.Request) *string {
	return getCookieValue(request, idRefreshTokenCookieKey)
}

func clearSessionFromCookie(response http.ResponseWriter, domain string,
	secure bool, accessTokenPath string, refreshTokenPath string, idRefreshTokenPath string, sameSite string) {
	setCookie(response, accessTokenCookieKey, "", domain, secure, true, 0, accessTokenPath, sameSite)
	setCookie(response, refreshTokenCookieKey, "", domain, secure, true, 0, refreshTokenPath, sameSite)
	setCookie(response, idRefreshTokenCookieKey, "", domain, secure, true, 0, idRefreshTokenPath, sameSite)
	setHeader(response, idRefreshTokenHeaderKey, "remove")
	setHeader(response, "Access-Control-Expose-Headers", idRefreshTokenHeaderKey)
}

func getRefreshTokenFromCookie(request *http.Request) *string {
	return getCookieValue(request, refreshTokenCookieKey)
}

func setCookie(response http.ResponseWriter, name string, value string,
	domain string, secure bool, httpOnly bool, expires uint64, path string, sameSite string) {

	if configMap != nil {
		if configMap.CookieDomain != "" {
			domain = configMap.CookieDomain
		}
		if configMap.CookieSecure != nil {
			secure = *configMap.CookieSecure
		}
		if configMap.CookieSameSite == "none" || configMap.CookieSameSite == "lax" ||
			configMap.CookieSameSite == "strict" {
			sameSite = configMap.CookieSameSite
		}
		if name == accessTokenCookieKey && configMap.AccessTokenPath != "" {
			path = configMap.AccessTokenPath
		}
		if name == idRefreshTokenCookieKey && configMap.AccessTokenPath != "" {
			path = configMap.AccessTokenPath
		}
		if name == refreshTokenCookieKey && configMap.RefreshAPIPath != "" {
			path = configMap.RefreshAPIPath
		}
	}

	var sameSiteField = http.SameSiteNoneMode
	if sameSite == "lax" {
		sameSiteField = http.SameSiteLaxMode
	} else if sameSite == "strict" {
		sameSiteField = http.SameSiteStrictMode
	}

	cookie := http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
		Expires:  time.Unix(int64(expires/1000), 0),
		Path:     path,
		SameSite: sameSiteField,
	}
	http.SetCookie(response, &cookie)
}

func setHeader(response http.ResponseWriter, key string, value string) {
	existingValue := response.Header().Get(strings.ToLower(key))
	if existingValue == "" {
		response.Header().Set(key, value)
	} else {
		response.Header().Set(key, existingValue+", "+value)
	}
}

func getHeader(request *http.Request, key string) *string {
	value := request.Header.Get(key)
	if value == "" {
		return nil
	}
	return &value
}

func getCookieValue(request *http.Request, key string) *string {
	cookies := request.Cookies()
	for _, value := range cookies {
		if value.Name == key {
			val, err := url.QueryUnescape(value.Value)
			if err != nil {
				return nil
			}
			return &val
		}
	}
	return nil
}

func setRelevantHeadersForOptionsAPI(response http.ResponseWriter) {
	setHeader(response, "Access-Control-Allow-Headers", antiCsrfHeaderKey)
	setHeader(response, "Access-Control-Allow-Headers", frontendSDKNameHeaderKey)
	setHeader(response, "Access-Control-Allow-Headers", frontendSDKVersionHeaderKey)
	setHeader(response, "Access-Control-Allow-Credentials", "true")
}

func getCORSAllowedHeaders() []string {
	return []string{
		antiCsrfHeaderKey, frontendSDKNameHeaderKey, frontendSDKVersionHeaderKey,
	}
}
