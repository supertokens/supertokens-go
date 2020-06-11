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

func TestQuerierCalledWithoutInit(t *testing.T) {
	beforeEach()
	core.GetQuerierInstance()
}

func TestCoreNotAvailable(t *testing.T) {
	beforeEach()
	supertokens.Config("localhost:8080;localhost:8081")
	q := core.GetQuerierInstance()
	_, err := q.SendGetRequest("/", map[string]string{})
	if err == nil && err.Error() != "Error while querying SuperTokens core" {
		t.Error("failed")
	}
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
	response, _ = q.SendDeleteRequest("/hello", map[string]interface{}{})
	if response == nil || response["result"].(string) != "Hello\n" {
		t.Error("failed")
		return
	}
	hostAlive := q.GetHostsAliveForTesting()

	if len(hostAlive) != 3 {
		t.Error("failed")
	}

	if !(containsHost(hostAlive, "localhost:8080") &&
		containsHost(hostAlive, "localhost:8081") && containsHost(hostAlive, "localhost:8082")) {
		t.Error("failed")
	}
}

func TestThreeCoresOneDeadRoundRobin(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	startST("localhost", "8082")
	supertokens.Config("localhost:8080;localhost:8081;localhost:8082")
	q := core.GetQuerierInstance()
	response, _ := q.SendGetRequest("/hello", map[string]string{})
	if response == nil || response["result"].(string) != "Hello\n" {
		t.Error("failed")
		return
	}
	response, _ = q.SendDeleteRequest("/hello", map[string]interface{}{})
	if response == nil || response["result"].(string) != "Hello\n" {
		t.Error("failed")
		return
	}
	hostAlive := q.GetHostsAliveForTesting()
	if len(hostAlive) != 2 {
		t.Error("failed")
	}

	response, _ = q.SendGetRequest("/hello", map[string]string{})
	if response == nil || response["result"].(string) != "Hello\n" {
		t.Error("failed")
		return
	}

	hostAlive = q.GetHostsAliveForTesting()
	if len(hostAlive) != 2 {
		t.Error("failed")
		return
	}
	if !(containsHost(hostAlive, "localhost:8080") &&
		containsHost(hostAlive, "localhost:8082")) {
		t.Error("failed")
	}
	if containsHost(hostAlive, "localhost:8081") {
		t.Error("failed")
	}
}
