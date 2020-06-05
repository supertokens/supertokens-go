package supertokens

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const accessTokenCookieKey = "sAccessToken"
const refreshTokenCookieKey = "sRefreshToken"

const idRefreshTokenCookieKey = "sIdRefreshToken"
const idRefreshTokenHeaderKey = "id-refresh-token"

const antiCsrfHeaderKey = "anti-csrf"
const frontendSDKNameHeaderKey = "supertokens-sdk-name"
const frontendSDKVersionHeaderKey = "supertokens-sdk-version"

func attachAccessTokenToCookie(response *http.ResponseWriter, token string,
	expiry uint64, domain string, path string, secure bool, sameSite string) {
	setCookie(response, accessTokenCookieKey, token, domain, secure, true, expiry, path, sameSite)
}

func attachRefreshTokenToCookie(response *http.ResponseWriter, token string,
	expiry uint64, domain string, path string, secure bool, sameSite string) {
	setCookie(response, refreshTokenCookieKey, token, domain, secure, true, expiry, path, sameSite)
}

func setIDRefreshTokenInHeaderAndCookie(response *http.ResponseWriter, token string,
	expiry uint64, domain string, path string, secure bool, sameSite string) {
	setHeader(response, idRefreshTokenHeaderKey, token+";"+fmt.Sprint(expiry))
	setHeader(response, "Access-Control-Expose-Headers", idRefreshTokenHeaderKey)

	setCookie(response, idRefreshTokenCookieKey, token, domain, secure, true, expiry, path, sameSite)
}

func setAntiCsrfTokenInHeaders(response *http.ResponseWriter, antiCsrfToken string) {
	setHeader(response, antiCsrfHeaderKey, antiCsrfToken)
	setHeader(response, "Access-Control-Expose-Headers", antiCsrfHeaderKey)
}

func saveFrontendInfoFromRequest(request *http.Request) {

	name := getHeader(request, frontendSDKNameHeaderKey)
	version := getHeader(request, frontendSDKVersionHeaderKey)
	if name != nil && version != nil {
		//TODO: add to device Driver info
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

func clearSessionFromCookie(response *http.ResponseWriter, domain string,
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

func setCookie(response *http.ResponseWriter, name string, value string,
	domain string, secure bool, httpOnly bool, expires uint64, path string, sameSite string) {

	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
		Expires:  time.Unix(int64(expires), 0),
		Path:     path,
	}
	http.SetCookie(*response, &cookie)
}

func setHeader(response *http.ResponseWriter, key string, value string) {
	existingValue := (*response).Header().Get(strings.ToLower(key))
	if existingValue == "" {
		(*response).Header().Set(key, value)
	} else {
		(*response).Header().Set(key, existingValue+", "+value)
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

	/* // parse JSON cookies
	    cookies = JSONCookies(cookies);

		return (cookies as any)[key];*/
	cookies := request.Cookies()
	for _, value := range cookies {
		if value.Name == key {
			return &value.Value
		}
	}
	return nil
}
