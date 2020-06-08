package supertokens

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/supertokens/supertokens-go/supertokens/core"
)

const accessTokenCookieKey = "sAccessToken"
const refreshTokenCookieKey = "sRefreshToken"

const idRefreshTokenCookieKey = "sIdRefreshToken"
const idRefreshTokenHeaderKey = "id-refresh-token"

const antiCsrfHeaderKey = "anti-csrf"
const frontendSDKNameHeaderKey = "supertokens-sdk-name"
const frontendSDKVersionHeaderKey = "supertokens-sdk-version"

func attachAccessTokenToCookie(response http.ResponseWriter, token string,
	expiry int64, domain string, secure bool, path string, sameSite string) {
	setCookie(response, accessTokenCookieKey, token, domain, secure, true, expiry, path, sameSite)
}

func attachRefreshTokenToCookie(response http.ResponseWriter, token string,
	expiry int64, domain string, secure bool, path string, sameSite string) {
	setCookie(response, refreshTokenCookieKey, token, domain, secure, true, expiry, path, sameSite)
}

func setIDRefreshTokenInHeaderAndCookie(response http.ResponseWriter, token string,
	expiry int64, domain string, secure bool, path string, sameSite string) {
	setHeader(response, idRefreshTokenHeaderKey, token+";"+fmt.Sprint(expiry))
	setHeader(response, "Access-Control-Expose-Headers", idRefreshTokenHeaderKey)

	setCookie(response, idRefreshTokenCookieKey, token, domain, secure, true, expiry, path, sameSite)
}

func setAntiCsrfTokenInHeaders(response http.ResponseWriter, antiCsrfToken string) {
	setHeader(response, antiCsrfHeaderKey, antiCsrfToken)
	setHeader(response, "Access-Control-Expose-Headers", antiCsrfHeaderKey)
}

func saveFrontendInfoFromRequest(request *http.Request) {
	name := getHeader(request, frontendSDKNameHeaderKey)
	version := getHeader(request, frontendSDKVersionHeaderKey)
	if name != nil && version != nil {
		core.GetDeviceInfoInstance().AddToFrontendSDKs(*name, *version)
	}
}

func getAccessTokenFromCookie(request *http.Request) *string {
	return getCookieValue(request, accessTokenCookieKey)
}

func getAntiCsrfTokenFromHeaders(request *http.Request) *string {
	return getHeader(request, antiCsrfHeaderKey)
}

func getIDRefreshTokenFromCookie(request *http.Request) *string {
	return getCookieValue(request, idRefreshTokenCookieKey)
}

func clearSessionFromCookie(response http.ResponseWriter, domain string,
	secure bool, accessTokenPath string, refreshTokenPath string, idRefreshTokenPath string, sameSite string) {
	setCookie(response, accessTokenCookieKey, "", domain, secure, true, 0, accessTokenPath, sameSite)
	setCookie(response, refreshTokenCookieKey, "", domain, secure, true, 0, refreshTokenPath, sameSite)
	setCookie(response, idRefreshTokenCookieKey, "", domain, secure, true, 0, idRefreshTokenPath, sameSite)
	setHeader(response, idRefreshTokenHeaderKey, "remove")
	setHeader(response, "Access-Control-Expose-Headers", idRefreshTokenHeaderKey)
}

func getRefreshTokenFromCookie(request *http.Request) *string {
	return getCookieValue(request, refreshTokenCookieKey)
}

func setCookie(response http.ResponseWriter, name string, value string,
	domain string, secure bool, httpOnly bool, expires int64, path string, sameSite string) {

	var sameSiteField = http.SameSiteNoneMode
	if sameSite == "lax" {
		sameSiteField = http.SameSiteLaxMode
	} else if sameSite == "strict" {
		sameSiteField = http.SameSiteStrictMode
	}

	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
		Expires:  time.Unix(int64(expires), 0),
		Path:     path,
		SameSite: sameSiteField,
	}
	http.SetCookie(response, &cookie)
}

func setHeader(response http.ResponseWriter, key string, value string) {
	existingValue := response.Header().Get(strings.ToLower(key))
	if existingValue == "" {
		response.Header().Set(key, value)
	} else {
		response.Header().Set(key, existingValue+", "+value)
	}
}

func getHeader(request *http.Request, key string) *string {
	value := request.Header.Get(key)
	if value == "" {
		return nil
	}
	return &value
}

func getCookieValue(request *http.Request, key string) *string {
	cookies := request.Cookies()
	for _, value := range cookies {
		if value.Name == key {
			return &value.Value
		}
	}
	return nil
}

func setRelevantHeadersForOptionsAPI(response http.ResponseWriter) {
	setHeader(response, "Access-Control-Allow-Headers", antiCsrfHeaderKey)
	setHeader(response, "Access-Control-Allow-Headers", frontendSDKNameHeaderKey)
	setHeader(response, "Access-Control-Allow-Headers", frontendSDKVersionHeaderKey)
	setHeader(response, "Access-Control-Allow-Credentials", "true")
}
