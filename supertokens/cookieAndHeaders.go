package supertokens

import "net/http"

func attachAccessTokenToCookie(response *http.ResponseWriter, token string,
	expiry uint64, domain string, path string, secure bool, sameSite string) {

}

func attachRefreshTokenToCookie(response *http.ResponseWriter, token string,
	expiry uint64, domain string, path string, secure bool, sameSite string) {

}

func setIDRefreshTokenInHeaderAndCookie(response *http.ResponseWriter, token string,
	expiry uint64, domain string, path string, secure bool, sameSite string) {

}

func setAntiCsrfTokenInHeaders(response *http.ResponseWriter, antiCsrfToken string) {

}

func saveFrontendInfoFromRequest(request *http.Request) {

}

func getAccessTokenFromCookie(request *http.Request) *string {
	return nil
}

func getAntiCsrfTokenFromHeaders(request *http.Request) *string {
	return nil
}

func getIDRefreshTokenFromCookie(request *http.Request) *string {
	return nil
}
