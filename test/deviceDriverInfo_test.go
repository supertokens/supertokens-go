package testing

import (
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
)

func TestDeviceDriveInfoWithoutFrontendSDK(t *testing.T) {
	beforeEach()

	startST("localhost", "8080")
	supertokens.Config("localhost:8080")
}
