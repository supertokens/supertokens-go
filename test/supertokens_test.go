package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
)

func TestSessionVerifyWithAntiCsrf(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(response http.ResponseWriter, requeset *http.Request) {
		supertokens.CreateNewSession(response, "id1")

	})
	mux.HandleFunc("/session/verify", func(response http.ResponseWriter, requeset *http.Request) {
		supertokens.GetSession(response, requeset, true)

	})
	mux.HandleFunc("/session/verifyAntiCsrfFalse", func(response http.ResponseWriter, requeset *http.Request) {
		supertokens.GetSession(response, requeset, false)

	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	res, _ := client.Do(req)
	response := extractInfoFromResponse(res)

	req, _ = http.NewRequest("POST", ts.URL+"/session/verify", nil)
	req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
	req.Header.Add("anti-csrf", response["antiCsrf"])
	res, _ = client.Do(req)

	if res.StatusCode != 200 {
		t.Error("response status code was not 200")
	}

	req, _ = http.NewRequest("POST", ts.URL+"/session/verifyAntiCsrfFalse", nil)
	req.Header.Add("Cookie", "sAccessToken="+response["accessToken"]+";sIdRefreshToken="+response["idRefreshTokenFromCookie"])
	req.Header.Add("anti-csrf", response["antiCsrf"])
	res, _ = client.Do(req)

	if res.StatusCode != 200 {
		t.Error("response status code was not 200")
	}
}
