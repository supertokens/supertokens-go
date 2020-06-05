package core

import (
	"net/http"

	"github.com/supertokens/supertokens-go/supertokens/errors"
)

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
	IDRefreshToken *TokenInfo
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
	sessionData map[string]interface{}) (SessionInfo, error) {
	response, err := getQuerierInstance().SendPostRequest("/session",
		map[string]interface{}{
			"userId":             userID,
			"userDataInJWT":      jwtPayload,
			"userDataInDatabase": sessionData,
		})
	if err != nil {
		return SessionInfo{}, err
	}
	return convertJSONResponseToSessionInfo(response), nil
}

// GetSession function used to verify a session
func GetSession(accessToken string, antiCsrfToken *string, doAntiCsrfCheck bool, idRefreshToken *string) (SessionInfo, error) {
	// TODO: Try verifying it here

	body := map[string]interface{}{
		"accessToken":     accessToken,
		"doAntiCsrfCheck": doAntiCsrfCheck,
	}
	if antiCsrfToken != nil {
		body["antiCsrfToken"] = *antiCsrfToken
	}
	response, err := getQuerierInstance().SendPostRequest("/session/verify", body)
	if err != nil {
		return SessionInfo{}, err
	}
	if response["status"] == "OK" {
		handShakeInfo, handShakeError := GetHandshakeInfoInstance()
		if handShakeError != nil {
			if err != nil {
				return SessionInfo{}, handShakeError
			}
		}
		handShakeInfo.UpdateJwtSigningPublicKeyInfo(
			response["jwtSigningPublicKey"].(string), response["jwtSigningPublicKeyExpiryTime"].(uint64))
		return convertJSONResponseToSessionInfo(response), nil
	} else if response["status"] == "UNAUTHORISED" {
		return SessionInfo{}, errors.UnauthorisedError{
			Msg: response["message"].(string),
		}
	} else {
		return SessionInfo{}, errors.TryRefreshTokenError{
			Msg: response["message"].(string),
		}
	}
}

// RefreshSession function used to refresh a session
func RefreshSession(refreshToken string) (SessionInfo, error) {
	response, err := getQuerierInstance().SendPostRequest("/session/refresh",
		map[string]interface{}{
			"refreshToken": refreshToken,
		})
	if err != nil {
		return SessionInfo{}, err
	}
	if response["status"] == "OK" {
		return convertJSONResponseToSessionInfo(response), nil
	} else if response["status"] == "UNAUTHORISED" {
		return SessionInfo{}, errors.UnauthorisedError{
			Msg: response["message"].(string),
		}
	} else {
		return SessionInfo{}, errors.TokenTheftDetectedError{
			Msg:           "Token theft detected",
			SessionHandle: (response["session"].(map[string]interface{}))["handle"].(string),
			UserID:        (response["session"].(map[string]interface{}))["userId"].(string),
		}
	}
}

// RevokeAllSessionsForUser function used to revoke all sessions for a user
func RevokeAllSessionsForUser(userID string) ([]string, error) {
	// TODO:
	return make([]string, 0), nil
}

// GetAllSessionHandlesForUser function used to get all sessions for a user
func GetAllSessionHandlesForUser(userID string) ([]string, error) {
	// TODO:
	return make([]string, 0), nil
}

// RevokeSession function used to revoke a specific session
func RevokeSession(sessionHandle string) (bool, error) {
	// TODO:
	return false, nil
}

// RevokeMultipleSessions function used to revoke a list of sessions
func RevokeMultipleSessions(sessionHandles []string) ([]string, error) {
	// TODO:
	return make([]string, 0), nil
}

// GetSessionData function used to get session data for the given handle
func GetSessionData(sessionHandle string) (map[string]interface{}, error) {
	// TODO:
	return map[string]interface{}{}, nil
}

// UpdateSessionData function used to update session data for the given handle
func UpdateSessionData(sessionHandle string, newSessionData map[string]interface{}) error {
	// TODO:
	return nil
}

// SetRelevantHeadersForOptionsAPI function is used to set headers specific to SuperTokens for OPTIONS API
func SetRelevantHeadersForOptionsAPI(response *http.ResponseWriter) error {
	// TODO:
	return nil
}

// GetJWTPayload function used to get jwt payload for the given handle
func GetJWTPayload(sessionHandle string) (map[string]interface{}, error) {
	// TODO:
	return map[string]interface{}{}, nil
}

// UpdateJWTPayload function used to update jwt payload for the given handle
func UpdateJWTPayload(sessionHandle string, newJWTPayload map[string]interface{}) error {
	// TODO:
	return nil
}

// RegenerateSession function used to regenerate a session
func RegenerateSession(accessToken string, newJWTPayload map[string]interface{}) (SessionInfo, error) {
	// TODO:
	return SessionInfo{}, nil
}
