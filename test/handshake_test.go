package testing

import (
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
)

func TestDriverInfoCheckWithoutFrontendSDK(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")
	info, err := core.GetHandshakeInfoInstance()
	if err != nil {
		t.Error(err)
	}

	if info.AccessTokenPath != "/" {
		t.Error("AccessToken path is not /")
	}
	if info.CookieDomain != "supertokens.io" {
		if info.CookieDomain != "localhost" {
			t.Error("incorrect cookie domain")
		}
	}
	if info.CookieSecure {
		t.Error("cookie secure set as true")
	}
	if info.RefreshTokenPath != "/refresh" {
		t.Error("incorrect refresh token path")
	}
	if !info.EnableAntiCsrf {
		t.Error("enable anticsrf set to false")
	}
	if info.AccessTokenBlacklistingEnabled {
		t.Error("accessTokenBlacklisting enabled")
	}
	info.UpdateJwtSigningPublicKeyInfo("hello", 100)

	info2, err2 := core.GetHandshakeInfoInstance()
	if err2 != nil {
		t.Error(err2)
	}

	if info2.JwtSigningPublicKey != "hello" {
		t.Error("JwtSigningPublicKey value does not match")
	}
	if info2.JwtSigningPublicKeyExpiryTime != 100 {
		t.Error("JwtSigningPublicKeyExpiryTime does not match")
	}
}
