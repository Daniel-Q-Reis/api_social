package integration_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Skip integration tests if running in a CI environment without database
	const skipIntegrationTests = "true"
	if os.Getenv("SKIP_INTEGRATION_TESTS") == skipIntegrationTests {
		return
	}

	code := m.Run()
	os.Exit(code)
}
