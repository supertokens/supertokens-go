package supertokens

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_setCookie_Once(t *testing.T) {
	w := httptest.NewRecorder()
	setCookie(w, "abc", "someVal", nil, false,
		false, 0, "/cookie", "lax")
	cookieMap := getCookieNameValuesMap(w)
	assert.Equal(t,"someVal", cookieMap["abc"])
}

func Test_setCookie_Replace(t *testing.T) {
	w := httptest.NewRecorder()
	setCookie(w, "abc", "someVal", nil, false,
		false, 0, "/cookie", "lax")
	cookieMap := getCookieNameValuesMap(w)
	assert.Equal(t,"someVal", cookieMap["abc"])

	setCookie(w, "abc", "someOtherVal", nil, false,
		false, 0, "/cookie", "lax")
	cookieMap = getCookieNameValuesMap(w)
	assert.Equal(t,"someOtherVal", cookieMap["abc"])
}

func Test_setCookie_Multiple(t *testing.T) {
	w := httptest.NewRecorder()
	setCookie(w, "abc", "valOne", nil, false,
		false, 0, "/cookie", "lax")
	cookieMap := getCookieNameValuesMap(w)
	assert.Equal(t,"valOne", cookieMap["abc"])

	setCookie(w, "xyz", "valOne", nil, false,
		false, 0, "/cookie", "lax")
	cookieMap = getCookieNameValuesMap(w)
	assert.Equal(t,"valOne", cookieMap["abc"])
	assert.Equal(t,"valOne", cookieMap["xyz"])

	setCookie(w, "abc", "valTwo", nil, false,
		false, 0, "/cookie", "lax")
	cookieMap = getCookieNameValuesMap(w)
	assert.Equal(t,"valTwo", cookieMap["abc"])
	assert.Equal(t,"valOne", cookieMap["xyz"])
}

func getCookieNameValuesMap(w http.ResponseWriter) map[string]string {
	cookieHeader := w.Header().Values("Set-Cookie")
	cookies := make(map[string]string, len(cookieHeader))
	for _, ch := range cookieHeader {
		parts := strings.Split(ch, ";")
		nameValue := strings.Split(parts[0], "=")
		cookies[getCookieName(ch)] = nameValue[1]
	}
	return cookies
}