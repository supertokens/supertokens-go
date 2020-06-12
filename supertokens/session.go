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
	"net/http"

	"github.com/supertokens/supertokens-go/supertokens/core"
	"github.com/supertokens/supertokens-go/supertokens/errors"
)

// Session object returned for managing a session
type Session struct {
	sessionHandle string
	userID        string
	userDataInJWT map[string]interface{}
	accessToken   string
	response      http.ResponseWriter
}

// RevokeSession function used to revoke a session for this session
func (session *Session) RevokeSession() error {
	success, err := RevokeSession(session.sessionHandle)
	if err != nil {
		return err
	}
	if success {
		handShakeInfo, handShakeInfoErr := core.GetHandshakeInfoInstance()
		if handShakeInfoErr != nil {
			return handShakeInfoErr
		}
		clearSessionFromCookie(session.response,
			handShakeInfo.CookieDomain,
			handShakeInfo.CookieSecure,
			handShakeInfo.AccessTokenPath,
			handShakeInfo.RefreshTokenPath,
			handShakeInfo.IDRefreshTokenPath,
			handShakeInfo.CookieSameSite,
		)
	}
	return nil
}

// GetSessionData function used to get session data for this session
func (session *Session) GetSessionData() (map[string]interface{}, error) {
	data, err := GetSessionData(session.sessionHandle)
	if err != nil {
		if errors.IsUnauthorizedError(err) {
			handShakeInfo, handShakeInfoErr := core.GetHandshakeInfoInstance()
			if handShakeInfoErr != nil {
				return nil, handShakeInfoErr
			}
			clearSessionFromCookie(session.response,
				handShakeInfo.CookieDomain,
				handShakeInfo.CookieSecure,
				handShakeInfo.AccessTokenPath,
				handShakeInfo.RefreshTokenPath,
				handShakeInfo.IDRefreshTokenPath,
				handShakeInfo.CookieSameSite,
			)
		}
		return nil, err
	}
	return data, nil
}

// UpdateSessionData function used to update session data for this session
func (session *Session) UpdateSessionData(newSessionData map[string]interface{}) error {
	err := UpdateSessionData(session.sessionHandle, newSessionData)
	if err != nil {
		if errors.IsUnauthorizedError(err) {
			handShakeInfo, handShakeInfoErr := core.GetHandshakeInfoInstance()
			if handShakeInfoErr != nil {
				return handShakeInfoErr
			}
			clearSessionFromCookie(session.response,
				handShakeInfo.CookieDomain,
				handShakeInfo.CookieSecure,
				handShakeInfo.AccessTokenPath,
				handShakeInfo.RefreshTokenPath,
				handShakeInfo.IDRefreshTokenPath,
				handShakeInfo.CookieSameSite,
			)
		}
		return err
	}
	return nil
}

// GetUserID function gets the user for this session
func (session *Session) GetUserID() string {
	return session.userID
}

// GetJWTPayload function gets the jwt payload for this session
func (session *Session) GetJWTPayload() map[string]interface{} {
	return session.userDataInJWT
}

// GetHandle function gets the session handle for this session
func (session *Session) GetHandle() string {
	return session.sessionHandle
}

// GetAccessToken function gets the access token for this session
func (session *Session) GetAccessToken() string {
	return session.accessToken
}

// UpdateJWTPayload function used to update jwt payload for this session
func (session *Session) UpdateJWTPayload(newJWTPayload map[string]interface{}) error {
	sessionInfo, err := core.RegenerateSession(session.accessToken, newJWTPayload)
	if err != nil {
		if errors.IsUnauthorizedError(err) {
			handShakeInfo, handShakeInfoErr := core.GetHandshakeInfoInstance()
			if handShakeInfoErr != nil {
				return handShakeInfoErr
			}
			clearSessionFromCookie(session.response,
				handShakeInfo.CookieDomain,
				handShakeInfo.CookieSecure,
				handShakeInfo.AccessTokenPath,
				handShakeInfo.RefreshTokenPath,
				handShakeInfo.IDRefreshTokenPath,
				handShakeInfo.CookieSameSite,
			)
		}
		return err
	}
	session.userDataInJWT = sessionInfo.UserDataInJWT
	if sessionInfo.AccessToken != nil {
		session.accessToken = (*sessionInfo.AccessToken).Token
		attachAccessTokenToCookie(
			session.response,
			(*sessionInfo.AccessToken).Token,
			(*sessionInfo.AccessToken).Expiry,
			(*sessionInfo.AccessToken).Domain,
			(*sessionInfo.AccessToken).CookieSecure,
			(*sessionInfo.AccessToken).CookiePath,
			(*sessionInfo.AccessToken).SameSite,
		)
	}
	return nil
}
