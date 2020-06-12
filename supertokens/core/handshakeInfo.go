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
	"sync"
)

type handshakeInfo struct {
	JwtSigningPublicKey            string
	CookieDomain                   string
	CookieSecure                   bool
	AccessTokenPath                string
	RefreshTokenPath               string
	EnableAntiCsrf                 bool
	AccessTokenBlacklistingEnabled bool
	JwtSigningPublicKeyExpiryTime  uint64
	CookieSameSite                 string
	IDRefreshTokenPath             string
	SessionExpiredStatusCode       int
}

var handshakeInfoInstantiated *handshakeInfo

// GetHandshakeInfoInstance returns handshake info.
func GetHandshakeInfoInstance() (*handshakeInfo, error) {
	if handshakeInfoInstantiated == nil {
		handshakeInfoLock.Lock()
		defer handshakeInfoLock.Unlock()
		if handshakeInfoInstantiated == nil {
			response, err := GetQuerierInstance().SendPostRequest("handshake", "/handshake", map[string]interface{}{})
			if err != nil {
				return nil, err
			}
			handshakeInfoInstantiated = &handshakeInfo{
				JwtSigningPublicKey:            response["jwtSigningPublicKey"].(string),
				CookieDomain:                   response["cookieDomain"].(string),
				CookieSecure:                   response["cookieSecure"].(bool),
				AccessTokenPath:                response["accessTokenPath"].(string),
				RefreshTokenPath:               response["refreshTokenPath"].(string),
				EnableAntiCsrf:                 response["enableAntiCsrf"].(bool),
				AccessTokenBlacklistingEnabled: response["accessTokenBlacklistingEnabled"].(bool),
				JwtSigningPublicKeyExpiryTime:  uint64(response["jwtSigningPublicKeyExpiryTime"].(float64)),
				CookieSameSite:                 response["cookieSameSite"].(string),
				IDRefreshTokenPath:             response["idRefreshTokenPath"].(string),
				SessionExpiredStatusCode:       int(response["sessionExpiredStatusCode"].(float64)),
			}
		}
	}
	return handshakeInfoInstantiated, nil
}

var handshakeInfoLock sync.Mutex

func (info *handshakeInfo) UpdateJwtSigningPublicKeyInfo(newKey string, newExpiry uint64) {
	handshakeInfoLock.Lock()
	defer handshakeInfoLock.Unlock()
	info.JwtSigningPublicKey = newKey
	info.JwtSigningPublicKeyExpiryTime = newExpiry
}

// ResetHandshakeInfo to be used for testing only
func ResetHandshakeInfo() {
	handshakeInfoInstantiated = nil
}
