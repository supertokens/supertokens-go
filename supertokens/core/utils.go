package core

import "time"

func convertJSONResponseToSessionInfo(response map[string]interface{}) SessionInfo {
	sessionJSON := response["session"].(map[string]interface{})
	accessTokenJSON := response["accessToken"].(map[string]interface{})
	refreshTokenJSON := response["refreshToken"].(map[string]interface{})
	idRefreshTokenJSON := response["idRefreshToken"].(map[string]interface{})
	antiCSRFTokenJSON := response["antiCsrfToken"]

	var accessToken *TokenInfo = nil
	if accessTokenJSON != nil {
		accessToken = &TokenInfo{
			Token:        accessTokenJSON["token"].(string),
			Expiry:       accessTokenJSON["expiry"].(int64),
			CreatedTime:  accessTokenJSON["createdTime"].(int64),
			CookiePath:   accessTokenJSON["cookiePath"].(string),
			CookieSecure: accessTokenJSON["cookieSecure"].(bool),
			Domain:       accessTokenJSON["domain"].(string),
			SameSite:     accessTokenJSON["sameSite"].(string),
		}
	}
	var refreshToken *TokenInfo = nil
	if refreshTokenJSON != nil {
		refreshToken = &TokenInfo{
			Token:        refreshTokenJSON["token"].(string),
			Expiry:       refreshTokenJSON["expiry"].(int64),
			CreatedTime:  refreshTokenJSON["createdTime"].(int64),
			CookiePath:   refreshTokenJSON["cookiePath"].(string),
			CookieSecure: refreshTokenJSON["cookieSecure"].(bool),
			Domain:       refreshTokenJSON["domain"].(string),
			SameSite:     refreshTokenJSON["sameSite"].(string),
		}
	}
	var idRefreshToken *TokenInfo = nil
	if idRefreshTokenJSON != nil {
		idRefreshToken = &TokenInfo{
			Token:        idRefreshTokenJSON["token"].(string),
			Expiry:       idRefreshTokenJSON["expiry"].(int64),
			CreatedTime:  idRefreshTokenJSON["createdTime"].(int64),
			CookiePath:   idRefreshTokenJSON["cookiePath"].(string),
			CookieSecure: idRefreshTokenJSON["cookieSecure"].(bool),
			Domain:       idRefreshTokenJSON["domain"].(string),
			SameSite:     idRefreshTokenJSON["sameSite"].(string),
		}
	}

	var antiCSRFToken *string = nil
	if antiCSRFTokenJSON != nil {
		str := antiCSRFTokenJSON.(string)
		antiCSRFToken = &str
	}
	return SessionInfo{
		Handle:         sessionJSON["handle"].(string),
		UserID:         sessionJSON["userId"].(string),
		UserDataInJWT:  sessionJSON["userDataInJWT"].(map[string]interface{}),
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		IDRefreshToken: idRefreshToken,
		AntiCsrfToken:  antiCSRFToken,
	}
}

func getCurrTimeInMS() int64 {
	return time.Now().UnixNano() / 1000000
}
