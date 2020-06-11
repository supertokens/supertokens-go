package testing

import (
	"fmt"
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
	"github.com/supertokens/supertokens-go/supertokens/errors"
)

func TestTokenTheftDetection(t *testing.T) {
	beforeEach()
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
	_, err = core.GetSession(response2.AccessToken.Token, response2.AntiCsrfToken, true)
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

func TestProcessState(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")

	p := core.GetProcessStateInstance()
	fmt.Println(p.GetLastEventByName(core.CallingServiceInVerify))
	p.AddState(core.CallingServiceInVerify)
	fmt.Println(*p.GetLastEventByName(core.CallingServiceInVerify))
}
func TestBasicSessionUse(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")

	response, err := core.CreateNewSession("", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	if response.AccessToken.Token == "" {
		t.Error("accessToken is empty")
	}
	if response.RefreshToken.Token == "" {
		t.Error("refreshToken is empty")
	}
	if response.IDRefreshToken.Token == "" {
		t.Error("idrefreshToken is empty")
	}
	if response.Handle == "" {
		t.Error("handle is empty")
	}
	if response.AntiCsrfToken == nil {
		t.Error("antiCsrfToken is nil")
	}

	_, err = core.GetSession(response.AccessToken.Token, response.AntiCsrfToken, true)
	if err != nil {
		t.Error(err)
	}

	p := core.GetProcessStateInstance()

	if p.GetLastEventByName(core.CallingServiceInVerify) != nil {
		t.Error("processState contains CallingServiceInVerify")
	}

	response2, err := core.RefreshSession(response.RefreshToken.Token)
	if err != nil {
		t.Error(err)
	}

	if response2.AccessToken.Token == "" {
		t.Error("accessToken is empty")
	}
	if response2.RefreshToken.Token == "" {
		t.Error("refreshToken is empty")
	}
	if response2.IDRefreshToken.Token == "" {
		t.Error("idrefreshToken is empty")
	}
	if response2.Handle == "" {
		t.Error("handle is empty")
	}
	if response2.AntiCsrfToken == nil {
		t.Error("antiCsrfToken is nil")
	}

	// response3, err := core.GetSession(response2.AccessToken.Token, response2.AntiCsrfToken, true)
	// if err != nil {
	// 	t.Error(err)
	// }

	// if p.GetLastEventByName(core.CallingServiceInVerify) == nil {
	// 	t.Error("processState does not contain CallingServiceInVerify")
	// }
	// if response3.Handle == "" {
	// 	t.Error("handle is empty")
	// }
	// if response3.AccessToken.Token == "" {
	// 	t.Error("accessToken is empty")
	// }
	// if response3.AntiCsrfToken != nil {
	// 	t.Error("antiCsrfToken is nil")
	// }
	// if response3.RefreshToken.Token == "" {
	// 	t.Error("refreshToken is empty")
	// }
	// if response3.IDRefreshToken.Token == "" {
	// 	t.Error("idrefreshToken is empty")
	// }
	// core.ResetProcessState()

	// response4, err := core.GetSession(response3.AccessToken.Token, response2.AntiCsrfToken, true)
	// if err != nil {
	// 	t.Error(err)
	// }
	// if core.GetProcessStateInstance().GetLastEventByName(core.CallingServiceInVerify) != nil {
	// 	t.Error("processState contains CallingServiceInVerify")
	// }
	// if response4.Handle == "" {
	// 	t.Error("handle is empty")
	// }
	// if response4.AccessToken.Token == "" {
	// 	t.Error("accessToken is empty")
	// }
	// if response4.AntiCsrfToken != nil {
	// 	t.Error("antiCsrfToken is nil")
	// }
	// if response4.RefreshToken.Token == "" {
	// 	t.Error("refreshToken is empty")
	// }
	// if response4.IDRefreshToken.Token == "" {
	// 	t.Error("idrefreshToken is empty")
	// }

	// revokeResponse, err := core.RevokeSession(response4.Handle)
	// if err != nil {
	// 	t.Error(err)
	// }
	// if !revokeResponse {
	// 	t.Error("session was not revoked")
	// }

}
