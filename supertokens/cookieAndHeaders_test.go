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
	cookieValue := joinCookieValues(w)
	assert.Equal(t,"abc=someVal", cookieValue)
}

func Test_setCookie_Twice(t *testing.T) {
	w := httptest.NewRecorder()
	setCookie(w, "abc", "someVal", nil, false,
		false, 0, "/cookie", "lax")
	joinedCookie := joinCookieValues(w)
	assert.Equal(t,"abc=someVal", joinedCookie)

	setCookie(w, "abc", "someOtherVal", nil, false,
		false, 0, "/cookie", "lax")
	updatedCookie := joinCookieValues(w)
	assert.Equal(t,"abc=someOtherVal", updatedCookie)
}

func joinCookieValues(w http.ResponseWriter) string {
	cookie := w.Header().Values("Set-Cookie")
	var values []string
	for _, v := range cookie {
		cv := strings.Split(v, ";")
		values = append(values, cv[0])
	}
	return strings.Join(values, ";")
}