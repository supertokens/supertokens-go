package supertokens

import (
	"github.com/supertokens/supertokens-go/supertokens"
)

// Session object returned for managing a session
type Session struct {
	actualSession *supertokens.Session
}

// RevokeSession function used to revoke a session for this session
func (session *Session) RevokeSession() error {
	return session.actualSession.RevokeSession()
}

// GetSessionData function used to get session data for this session
func (session *Session) GetSessionData() (map[string]interface{}, error) {
	return session.actualSession.GetSessionData()
}

// UpdateSessionData function used to update session data for this session
func (session *Session) UpdateSessionData(newSessionData map[string]interface{}) error {
	return session.actualSession.UpdateSessionData(newSessionData)
}

// GetUserID function gets the user for this session
func (session *Session) GetUserID() string {
	return session.actualSession.GetUserID()
}

// GetJWTPayload function gets the jwt payload for this session
func (session *Session) GetJWTPayload() map[string]interface{} {
	return session.actualSession.GetJWTPayload()
}

// GetHandle function gets the session handle for this session
func (session *Session) GetHandle() string {
	return session.actualSession.GetHandle()
}

// GetAccessToken function gets the access token for this session
func (session *Session) GetAccessToken() string {
	return session.actualSession.GetAccessToken()
}

// UpdateJWTPayload function used to update jwt payload for this session
func (session *Session) UpdateJWTPayload(newJWTPayload map[string]interface{}) error {
	return session.actualSession.UpdateJWTPayload(newJWTPayload)
}
