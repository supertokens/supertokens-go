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

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-go/supertokens"
)

// SessionContext string to get session struct from context if using Gin
const sessionContext string = "supertokens_session_key"

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
	supertokens.Config(supertokens.ConfigMap{
		Hosts:           config.Hosts,
		AccessTokenPath: config.AccessTokenPath,
		RefreshAPIPath:  config.RefreshAPIPath,
		CookieDomain:    config.CookieDomain,
		CookieSecure:    config.CookieSecure,
		CookieSameSite:  config.CookieSameSite,
		APIKey:          config.APIKey,
	})
}

// CreateNewSession function used to create a new SuperTokens session
func CreateNewSession(c *gin.Context, userID string,
	payload ...map[string]interface{}) (Session, error) {
	actualSession, err := supertokens.CreateNewSession(c.Writer, userID, payload...)
	if err != nil {
		return Session{}, err
	}
	return Session{
		actualSession: &actualSession,
	}, nil
}

// GetSession function used to verify a session
func GetSession(c *gin.Context, doAntiCsrfCheck bool) (Session, error) {
	actualSession, err := supertokens.GetSession(c.Writer, c.Request, doAntiCsrfCheck)
	if err != nil {
		return Session{}, err
	}
	return Session{
		actualSession: &actualSession,
	}, nil
}

// RefreshSession function used to refresh a session
func RefreshSession(c *gin.Context) (Session, error) {
	actualSession, err := supertokens.RefreshSession(c.Writer, c.Request)
	if err != nil {
		return Session{}, err
	}
	return Session{
		actualSession: &actualSession,
	}, nil
}

// RevokeAllSessionsForUser function used to revoke all sessions for a user
func RevokeAllSessionsForUser(userID string) ([]string, error) {
	return supertokens.RevokeAllSessionsForUser(userID)
}

// GetAllSessionHandlesForUser function used to get all sessions for a user
func GetAllSessionHandlesForUser(userID string) ([]string, error) {
	return supertokens.GetAllSessionHandlesForUser(userID)
}

// RevokeSession function used to revoke a specific session
func RevokeSession(sessionHandle string) (bool, error) {
	return supertokens.RevokeSession(sessionHandle)
}

// RevokeMultipleSessions function used to revoke a list of sessions
func RevokeMultipleSessions(sessionHandles []string) ([]string, error) {
	return supertokens.RevokeMultipleSessions(sessionHandles)
}

// GetSessionData function used to get session data for the given handle
func GetSessionData(sessionHandle string) (map[string]interface{}, error) {
	return supertokens.GetSessionData(sessionHandle)
}

// UpdateSessionData function used to update session data for the given handle
func UpdateSessionData(sessionHandle string, newSessionData map[string]interface{}) error {
	return supertokens.UpdateSessionData(sessionHandle, newSessionData)
}

// SetRelevantHeadersForOptionsAPI function is used to set headers specific to SuperTokens for OPTIONS API
func SetRelevantHeadersForOptionsAPI(c *gin.Context) {
	supertokens.SetRelevantHeadersForOptionsAPI(c.Writer)
}

// GetCORSAllowedHeaders function is used to get header keys that are used by SuperTokens
func GetCORSAllowedHeaders() []string {
	return supertokens.GetCORSAllowedHeaders()
}

// GetJWTPayload function used to get jwt payload for the given handle
func GetJWTPayload(sessionHandle string) (map[string]interface{}, error) {
	return supertokens.GetJWTPayload(sessionHandle)
}

// UpdateJWTPayload function used to update jwt payload for the given handle
func UpdateJWTPayload(sessionHandle string, newJWTPayload map[string]interface{}) error {
	return supertokens.UpdateJWTPayload(sessionHandle, newJWTPayload)
}

// OnTokenTheftDetected function to override default behaviour of handling token thefts
func OnTokenTheftDetected(handler func(string, string, http.ResponseWriter)) {
	supertokens.OnTokenTheftDetected(handler)
}

// OnUnauthorized function to override default behaviour of handling Unauthorized error
func OnUnauthorized(handler func(error, http.ResponseWriter)) {
	supertokens.OnUnauthorized(handler)
}

// OnTryRefreshToken function to override default behaviour of handling try refresh token errors
func OnTryRefreshToken(handler func(error, http.ResponseWriter)) {
	supertokens.OnTryRefreshToken(handler)
}

// OnGeneralError function to override default behaviour of handling general errors
func OnGeneralError(handler func(error, http.ResponseWriter)) {
	supertokens.OnGeneralError(handler)
}

// GetSessionFromRequest returns the verified session object if present, otherwise returns nil
func GetSessionFromRequest(c *gin.Context) *Session {
	value, exists := c.Get(sessionContext)
	if exists {
		temp := value.(*Session)
		return temp
	}
	return nil
}
