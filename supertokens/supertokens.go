package supertokens

import (
	"net/http"

	"github.com/supertokens/supertokens-go/supertokens/core"
)

// Config used to set locations of SuperTokens instances
func Config(hosts string) error {
	return core.Config(hosts)
}

// CreateNewSession function used to create a new SuperTokens session
func CreateNewSession(response *http.ResponseWriter,
	userID string, payload ...map[string]interface{}) (Session, error) {
	// TODO:

	//check for size of payload <= 2
	if len(payload) > 2 {
		//handle error
	}

	session, err := core.CreateNewSession(userID, payload[0], payload[1])
	if err != nil {
		//handle error
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
		accessToken.CookiePath,
		accessToken.CookieSecure,
		accessToken.SameSite,
	)

	attachRefreshTokenToCookie(
		response,
		refreshToken.Token,
		refreshToken.Expiry,
		refreshToken.Domain,
		refreshToken.CookiePath,
		refreshToken.CookieSecure,
		refreshToken.SameSite,
	)

	setIDRefreshTokenInHeaderAndCookie(
		response,
		idRefreshToken.Token,
		idRefreshToken.Expiry,
		idRefreshToken.Domain,
		idRefreshToken.CookiePath,
		idRefreshToken.CookieSecure,
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
func GetSession(response *http.ResponseWriter, request *http.Request,
	doAntiCsrfCheck bool) (Session, error) {
	// TODO:

	saveFrontendInfoFromRequest(request)

	accessToken := getAccessTokenFromCookie(request)

	if accessToken == nil {
		// maybe the access token has expired.
		// return try refresh token error msg:"access token missing in cookies"
	}

	// antiCsrfToken := getAntiCsrfTokenFromHeaders(request)
	// idRefreshToken := getIDRefreshTokenFromCookie(request)

	// session := core.GetSession(response, request, doAntiCsrfCheck)

	return Session{}, nil
}

// RefreshSession function used to refresh a session
func RefreshSession(response *http.ResponseWriter, request *http.Request) (Session, error) {
	// TODO:
	return Session{}, nil
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
func SetRelevantHeadersForOptionsAPI(response *http.ResponseWriter) error {
	return core.SetRelevantHeadersForOptionsAPI(response)
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

// OnUnauthorised function to override default behaviour of handling unauthorised error
func OnUnauthorised(handler func(error, http.ResponseWriter)) {
	core.GetErrorHandlersInstance().OnUnauthorisedErrorHandler = handler
}

// OnTryRefreshToken function to override default behaviour of handling try refresh token errors
func OnTryRefreshToken(handler func(error, http.ResponseWriter)) {
	core.GetErrorHandlersInstance().OnTryRefreshTokenErrorHandler = handler
}

// OnGeneralError function to override default behaviour of handling general errors
func OnGeneralError(handler func(error, http.ResponseWriter)) {
	core.GetErrorHandlersInstance().OnGeneralErrorHandler = handler
}
