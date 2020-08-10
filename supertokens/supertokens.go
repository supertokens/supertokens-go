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

type contextKey int

const sessionContext contextKey = iota

// ConfigMap add key value params for session behaviour
type ConfigMap struct {
	Hosts           string
	AccessTokenPath string
	RefreshAPIPath  string
	CookieDomain    string
	CookieSecure    *bool
	CookieSameSite  string
	APIKey          string
}

// Config used to set locations of SuperTokens instances
func Config(config ConfigMap) {
	configCookieAndHeaders(config)
	core.Config(config.Hosts, config.APIKey)
}

// CreateNewSession function used to create a new SuperTokens session
func CreateNewSession(response http.ResponseWriter,
	userID string, payload ...map[string]interface{}) (Session, error) {

	var jwtPayload = map[string]interface{}{}
	var sessionData = map[string]interface{}{}
	if len(payload) == 1 && payload[0] != nil {
		jwtPayload = payload[0]
	} else if len(payload) == 2 {
		if payload[0] != nil {
			jwtPayload = payload[0]
		}
		if payload[1] != nil {
			sessionData = payload[1]
		}
	}

	session, err := core.CreateNewSession(userID, jwtPayload, sessionData)

	if err != nil {
		return Session{}, err
	}

	//attach token to cookies
	accessToken := session.AccessToken
	refreshToken := session.RefreshToken
	idRefreshToken := session.IDRefreshToken

	attachAccessTokenToCookie(
		response,
		accessToken.Token,
		accessToken.Expiry,
		accessToken.Domain,
		accessToken.CookieSecure,
		accessToken.CookiePath,
		accessToken.SameSite,
	)

	attachRefreshTokenToCookie(
		response,
		refreshToken.Token,
		refreshToken.Expiry,
		refreshToken.Domain,
		refreshToken.CookieSecure,
		refreshToken.CookiePath,
		refreshToken.SameSite,
	)

	setIDRefreshTokenInHeaderAndCookie(
		response,
		idRefreshToken.Token,
		idRefreshToken.Expiry,
		idRefreshToken.Domain,
		idRefreshToken.CookieSecure,
		idRefreshToken.CookiePath,
		idRefreshToken.SameSite,
	)

	if session.AntiCsrfToken != nil {
		setAntiCsrfTokenInHeaders(response, *session.AntiCsrfToken)
	}

	return Session{
		accessToken:   accessToken.Token,
		sessionHandle: session.Handle,
		userID:        session.UserID,
		userDataInJWT: session.UserDataInJWT,
		response:      response,
	}, nil

}

// GetSession function used to verify a session
func GetSession(response http.ResponseWriter, request *http.Request,
	doAntiCsrfCheck bool) (Session, error) {
	saveFrontendInfoFromRequest(request)

	idRefreshToken := getIDRefreshTokenFromCookie(request)
	if idRefreshToken == nil {
		handShakeInfo, handshakeInfoError := core.GetHandshakeInfoInstance()
		if handshakeInfoError != nil {
			return Session{}, handshakeInfoError
		}
		clearSessionFromCookie(response,
			handShakeInfo.CookieDomain,
			handShakeInfo.CookieSecure,
			handShakeInfo.AccessTokenPath,
			handShakeInfo.RefreshTokenPath,
			handShakeInfo.IDRefreshTokenPath,
			handShakeInfo.CookieSameSite,
		)
		return Session{}, errors.UnauthorizedError{
			Msg: "idRefreshToken missing",
		}
	}

	accessToken := getAccessTokenFromCookie(request)
	if accessToken == nil {
		// maybe the access token has expired.
		return Session{}, errors.TryRefreshTokenError{
			Msg: "access token missing in cookies",
		}
	}

	antiCsrfToken := getAntiCsrfTokenFromHeaders(request)

	session, getSessionError := core.GetSession(*accessToken, antiCsrfToken, doAntiCsrfCheck)

	if getSessionError != nil {
		if errors.IsUnauthorizedError(getSessionError) {
			handShakeInfo, handshakeInfoError := core.GetHandshakeInfoInstance()
			if handshakeInfoError != nil {
				return Session{}, handshakeInfoError
			}
			clearSessionFromCookie(response,
				handShakeInfo.CookieDomain,
				handShakeInfo.CookieSecure,
				handShakeInfo.AccessTokenPath,
				handShakeInfo.RefreshTokenPath,
				handShakeInfo.IDRefreshTokenPath,
				handShakeInfo.CookieSameSite,
			)
		}
		return Session{}, getSessionError
	}

	if session.AccessToken != nil {
		attachAccessTokenToCookie(
			response,
			session.AccessToken.Token,
			session.AccessToken.Expiry,
			session.AccessToken.Domain,
			session.AccessToken.CookieSecure,
			session.AccessToken.CookiePath,
			session.AccessToken.SameSite,
		)
		accessToken = &session.AccessToken.Token
	}

	return Session{
		accessToken:   *accessToken,
		response:      response,
		sessionHandle: session.Handle,
		userDataInJWT: session.UserDataInJWT,
		userID:        session.UserID,
	}, nil
}

// RefreshSession function used to refresh a session
func RefreshSession(response http.ResponseWriter, request *http.Request) (Session, error) {
	saveFrontendInfoFromRequest(request)
	inputRefreshToken := getRefreshTokenFromCookie(request)
	if inputRefreshToken == nil {
		handShakeInfo, handshakeInfoError := core.GetHandshakeInfoInstance()
		if handshakeInfoError != nil {
			return Session{}, handshakeInfoError
		}
		clearSessionFromCookie(
			response,
			handShakeInfo.CookieDomain,
			handShakeInfo.CookieSecure,
			handShakeInfo.AccessTokenPath,
			handShakeInfo.RefreshTokenPath,
			handShakeInfo.IDRefreshTokenPath,
			handShakeInfo.CookieSameSite)
		return Session{}, errors.UnauthorizedError{
			Msg: "Missing auth tokens in cookies. Have you set the correct refresh API path in your frontend and SuperTokens config?",
		}
	}

	session, refreshError := core.RefreshSession(*inputRefreshToken)

	if refreshError != nil {

		if errors.IsUnauthorizedError(refreshError) || errors.IsTokenTheftDetectedError(refreshError) {
			handShakeInfo, handshakeInfoError := core.GetHandshakeInfoInstance()
			if handshakeInfoError != nil {
				return Session{}, handshakeInfoError
			}
			clearSessionFromCookie(
				response,
				handShakeInfo.CookieDomain,
				handShakeInfo.CookieSecure,
				handShakeInfo.AccessTokenPath,
				handShakeInfo.RefreshTokenPath,
				handShakeInfo.IDRefreshTokenPath,
				handShakeInfo.CookieSameSite)
		}
		return Session{}, refreshError
	}

	//attach cookies
	accessToken := session.AccessToken
	refreshToken := session.RefreshToken
	idRefreshToken := session.IDRefreshToken

	attachAccessTokenToCookie(
		response,
		accessToken.Token,
		accessToken.Expiry,
		accessToken.Domain,
		accessToken.CookieSecure,
		accessToken.CookiePath,
		accessToken.SameSite,
	)

	attachRefreshTokenToCookie(
		response,
		refreshToken.Token,
		refreshToken.Expiry,
		refreshToken.Domain,
		refreshToken.CookieSecure,
		refreshToken.CookiePath,
		refreshToken.SameSite,
	)

	setIDRefreshTokenInHeaderAndCookie(
		response,
		idRefreshToken.Token,
		idRefreshToken.Expiry,
		idRefreshToken.Domain,
		idRefreshToken.CookieSecure,
		idRefreshToken.CookiePath,
		idRefreshToken.SameSite,
	)

	if session.AntiCsrfToken != nil {
		setAntiCsrfTokenInHeaders(response, *session.AntiCsrfToken)
	}

	return Session{
		accessToken:   accessToken.Token,
		sessionHandle: session.Handle,
		userID:        session.UserID,
		userDataInJWT: session.UserDataInJWT,
		response:      response,
	}, nil
}

// RevokeAllSessionsForUser function used to revoke all sessions for a user
func RevokeAllSessionsForUser(userID string) ([]string, error) {
	return core.RevokeAllSessionsForUser(userID)
}

// GetAllSessionHandlesForUser function used to get all sessions for a user
func GetAllSessionHandlesForUser(userID string) ([]string, error) {
	return core.GetAllSessionHandlesForUser(userID)
}

// RevokeSession function used to revoke a specific session
func RevokeSession(sessionHandle string) (bool, error) {
	return core.RevokeSession(sessionHandle)
}

// RevokeMultipleSessions function used to revoke a list of sessions
func RevokeMultipleSessions(sessionHandles []string) ([]string, error) {
	return core.RevokeMultipleSessions(sessionHandles)
}

// GetSessionData function used to get session data for the given handle
func GetSessionData(sessionHandle string) (map[string]interface{}, error) {
	return core.GetSessionData(sessionHandle)
}

// UpdateSessionData function used to update session data for the given handle
func UpdateSessionData(sessionHandle string, newSessionData map[string]interface{}) error {
	return core.UpdateSessionData(sessionHandle, newSessionData)
}

// SetRelevantHeadersForOptionsAPI function is used to set headers specific to SuperTokens for OPTIONS API
func SetRelevantHeadersForOptionsAPI(response http.ResponseWriter) {
	setRelevantHeadersForOptionsAPI(response)
}

// GetCORSAllowedHeaders function is used to get header keys that are used by SuperTokens
func GetCORSAllowedHeaders() []string {
	return getCORSAllowedHeaders()
}

// GetJWTPayload function used to get jwt payload for the given handle
func GetJWTPayload(sessionHandle string) (map[string]interface{}, error) {
	return core.GetJWTPayload(sessionHandle)
}

// UpdateJWTPayload function used to update jwt payload for the given handle
func UpdateJWTPayload(sessionHandle string, newJWTPayload map[string]interface{}) error {
	return core.UpdateJWTPayload(sessionHandle, newJWTPayload)
}

// OnTokenTheftDetected function to override default behaviour of handling token thefts
func OnTokenTheftDetected(handler func(string, string, http.ResponseWriter)) {
	core.GetErrorHandlersInstance().OnTokenTheftDetectedErrorHandler = handler
}

// OnUnauthorized function to override default behaviour of handling Unauthorized error
func OnUnauthorized(handler func(error, http.ResponseWriter)) {
	core.GetErrorHandlersInstance().OnUnauthorizedErrorHandler = handler
}

// OnTryRefreshToken function to override default behaviour of handling try refresh token errors
func OnTryRefreshToken(handler func(error, http.ResponseWriter)) {
	core.GetErrorHandlersInstance().OnTryRefreshTokenErrorHandler = handler
}

// OnGeneralError function to override default behaviour of handling general errors
func OnGeneralError(handler func(error, http.ResponseWriter)) {
	core.GetErrorHandlersInstance().OnGeneralErrorHandler = handler
}

// GetSessionFromRequest returns the verified session object if present, otherwise returns nil
func GetSessionFromRequest(r *http.Request) *Session {
	value := r.Context().Value(sessionContext)
	if value == nil {
		return nil
	}
	temp := value.(Session)
	return &temp
}
