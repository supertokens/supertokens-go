package testing

import (
	"os"
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
)

func TestMain(m *testing.M) {
	code := m.Run()
	killAllST()
	cleanST()
	os.Exit(code)
}

func TestThreeCoresAndRoundRobin(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	startST("localhost", "8081")
	startST("localhost", "8082")
	supertokens.Config("localhost:8080;localhost:8081;localhost:8082")
	q := core.GetQuerierInstance()
	response, _ := q.SendGetRequest("/hello", map[string]string{})
	if response == nil || response["result"].(string) != "Hello\n" {
		t.Error("failed")
		return
	}
	// TODO:
}
