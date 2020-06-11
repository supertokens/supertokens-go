package testing

import (
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
	"github.com/supertokens/supertokens-go/supertokens/errors"
)

func TestTokenTheftDetection(t *testing.T) {
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")

	response, err := core.CreateNewSession("", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	response2, err := core.RefreshSession(response.RefreshToken.Token)
	if err != nil {
		t.Error(err)
	}
	_, err = core.GetSession(response2.AccessToken.Token, response2.AntiCsrfToken, true, &response2.IDRefreshToken.Token)
	if err != nil {
		t.Error(err)
	}
	_, err = core.RefreshSession(response.RefreshToken.Token)
	if err == nil {
		t.Error("should not have come here")
	} else if !errors.IsTokenTheftDetectedError(err) {
		t.Error("failed")
	}

}
