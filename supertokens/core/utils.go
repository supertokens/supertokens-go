package core

import (
	"strconv"
	"strings"
	"time"
)

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

func getLargestVersionFromIntersection(v1 []string, v2 []string) *string {
	var intersection = []string{}
	for i := 0; i < len(v1); i++ {
		for y := 0; y < len(v2); y++ {
			if v1[i] == v2[y] {
				intersection = append(intersection, v1[i])
			}
		}
	}
	if len(intersection) == 0 {
		return nil
	}
	maxVersionSoFar := intersection[0]
	for i := 1; i < len(intersection); i++ {
		maxVersionSoFar = maxVersion(intersection[i], maxVersionSoFar)
	}
	return &maxVersionSoFar
}

func maxVersion(version1 string, version2 string) string {
	var splittedv1 = strings.Split(version1, ".")
	var splittedv2 = strings.Split(version2, ".")
	var minLength = len(splittedv1)
	if minLength > len(splittedv2) {
		minLength = len(splittedv2)
	}
	for i := 0; i < minLength; i++ {
		var v1, _ = strconv.Atoi(splittedv1[i])
		var v2, _ = strconv.Atoi(splittedv2[i])
		if v1 > v2 {
			return version1
		} else if v2 > v1 {
			return version2
		}
	}
	if len(splittedv1) >= len(splittedv2) {
		return version1
	}
	return version2
}
