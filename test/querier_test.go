package testing

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	killAllST()
	cleanST()
	os.Exit(code)
}

func TestExample(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
}
