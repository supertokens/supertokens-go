package supertokens

import "net/http"

// Session object returned for managing a session
type Session struct {
	sessionHandle string
	userID        string
	userDataInJWT map[string]interface{}
	accessToken   string
	response      *http.ResponseWriter
}

// RevokeSession function used to revoke a session for this session
func (session *Session) RevokeSession() {
	// TODO:
}

// GetSessionData function used to get session data for this session
func (session *Session) GetSessionData() map[string]interface{} {
	// TODO:
	return map[string]interface{}{}
}

// UpdateSessionData function used to update session data for this session
func (session *Session) UpdateSessionData(newSessionData map[string]interface{}) {
	// TODO:
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
func (session *Session) UpdateJWTPayload(newJWTPayload map[string]interface{}) {
	// TODO:
}
