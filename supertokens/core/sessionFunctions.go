package core

import "net/http"

// Config used to set locations of SuperTokens instances
func Config(hosts string) error {
	return InitQuerier(hosts)
}

// SessionInfo carrier of session token information
type SessionInfo struct {
	// some of the fields are points cause they can be nil too
	Handle         string
	UserID         string
	UserDataInJWT  map[string]interface{}
	AccessToken    *TokenInfo
	RefreshToken   *TokenInfo
	IdRefreshToken *TokenInfo
	AntiCsrfToken  *string
}

// TokenInfo carrier of cookie related info for a token
type TokenInfo struct {
	Token        string
	Expiry       uint64
	CreatedTime  uint64
	CookiePath   string
	CookieSecure bool
	Domain       string
	SameSite     string
}

// CreateNewSession function used to create a new SuperTokens session
func CreateNewSession(userID string, jwtPayload map[string]interface{},
	sessionData map[string]interface{}) SessionInfo {
	// TODO:
	return SessionInfo{}
}

// GetSession function used to verify a session
func GetSession(response *http.ResponseWriter, request *http.Request,
	doAntiCsrfCheck bool) SessionInfo {
	// TODO:
	return SessionInfo{}
}

// RefreshSession function used to refresh a session
func RefreshSession(response *http.ResponseWriter, request *http.Request) SessionInfo {
	// TODO:
	return SessionInfo{}
}

// RevokeAllSessionsForUser function used to revoke all sessions for a user
func RevokeAllSessionsForUser(userID string) []string {
	// TODO:
	return make([]string, 0)
}

// GetAllSessionHandlesForUser function used to get all sessions for a user
func GetAllSessionHandlesForUser(userID string) []string {
	// TODO:
	return make([]string, 0)
}

// RevokeSession function used to revoke a specific session
func RevokeSession(sessionHandle string) bool {
	// TODO:
	return false
}

// RevokeMultipleSessions function used to revoke a list of sessions
func RevokeMultipleSessions(sessionHandles []string) []string {
	// TODO:
	return make([]string, 0)
}

// GetSessionData function used to get session data for the given handle
func GetSessionData(sessionHandle string) map[string]interface{} {
	// TODO:
	return map[string]interface{}{}
}

// UpdateSessionData function used to update session data for the given handle
func UpdateSessionData(sessionHandle string, newSessionData map[string]interface{}) {
	// TODO:
}

// SetRelevantHeadersForOptionsAPI function is used to set headers specific to SuperTokens for OPTIONS API
func SetRelevantHeadersForOptionsAPI(response *http.ResponseWriter) {
	// TODO:
}

// GetJWTPayload function used to get jwt payload for the given handle
func GetJWTPayload(sessionHandle string) map[string]interface{} {
	// TODO:
	return map[string]interface{}{}
}

// UpdateJWTPayload function used to update jwt payload for the given handle
func UpdateJWTPayload(sessionHandle string, newJWTPayload map[string]interface{}) {
	// TODO:
}

// RegenerateSession function used to regenerate a session
func RegenerateSession(accessToken string, newJWTPayload map[string]interface{}) SessionInfo {
	// TODO:
	return SessionInfo{}
}
