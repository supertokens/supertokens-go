package supertokens

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-go/supertokens"
)

// SessionContext string to get session struct from context if using Gin
const sessionContext string = "supertokens_session_key"

// Config used to set locations of SuperTokens instances
func Config(hosts string) error {
	return supertokens.Config(hosts)
}

// CreateNewSession function used to create a new SuperTokens session
func CreateNewSession(c *gin.Context, userID string,
	payload ...map[string]interface{}) (supertokens.Session, error) {
	return supertokens.CreateNewSession(c.Writer, userID, payload...)
}

// GetSession function used to verify a session
func GetSession(c *gin.Context, doAntiCsrfCheck bool) (supertokens.Session, error) {
	return supertokens.GetSession(c.Writer, c.Request, doAntiCsrfCheck)
}

// RefreshSession function used to refresh a session
func RefreshSession(c *gin.Context) (supertokens.Session, error) {
	return supertokens.RefreshSession(c.Writer, c.Request)
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

// GetSessionFromRequest returns the verified session object if present, otherwise it panics
func GetSessionFromRequest(c *gin.Context) supertokens.Session {
	return c.MustGet(sessionContext).(supertokens.Session)
}
