package supertokens

import (
	"net/http"

	"github.com/supertokens/supertokens-go/supertokens/core"
	"github.com/supertokens/supertokens-go/supertokens/errors"
)

type contextKey int

// SessionContext string to get the session struct from context
const SessionContext contextKey = iota

// Config used to set locations of SuperTokens instances
func Config(hosts string) error {
	return core.Config(hosts)
}

// CreateNewSession function used to create a new SuperTokens session
func CreateNewSession(response *http.ResponseWriter,
	userID string, payload ...map[string]interface{}) (Session, error) {

	var jwtPayload = map[string]interface{}{}
	var sessionData = map[string]interface{}{}
	if len(payload) == 1 {
		jwtPayload = payload[0]
	} else if len(payload) == 2 {
		jwtPayload = payload[0]
		sessionData = payload[1]
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
func GetSession(response *http.ResponseWriter, request *http.Request,
	doAntiCsrfCheck bool) (Session, error) {
	saveFrontendInfoFromRequest(request)
	accessToken := getAccessTokenFromCookie(request)
	if accessToken == nil {
		// maybe the access token has expired.
		return Session{}, errors.TryRefreshTokenError{
			Msg: "access token missing in cookies",
		}
	}

	antiCsrfToken := getAntiCsrfTokenFromHeaders(request)
	idRefreshToken := getIDRefreshTokenFromCookie(request)

	session, getSessionError := core.GetSession(*accessToken, antiCsrfToken, doAntiCsrfCheck, idRefreshToken)

	if getSessionError != nil {
		if errors.IsUnauthorisedError(getSessionError) {
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
func RefreshSession(response *http.ResponseWriter, request *http.Request) (Session, error) {
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
		return Session{}, errors.UnauthorisedError{
			Msg: "Missing auth tokens in cookies. Have you set the correct refresh API path in your frontend and SuperTokens config?",
		}
	}

	session, refreshError := core.RefreshSession(*inputRefreshToken)

	if refreshError != nil {

		if errors.IsUnauthorisedError(refreshError) || errors.IsTokenTheftDetectedError(refreshError) {
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
func SetRelevantHeadersForOptionsAPI(response *http.ResponseWriter) {
	setRelevantHeadersForOptionsAPI(response)
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
