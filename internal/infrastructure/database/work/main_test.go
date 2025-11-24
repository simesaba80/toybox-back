//go:build integration

package work_test

import (
	"os"
	"testing"

	"github.com/simesaba80/toybox-back/internal/infrastructure/database/testutil"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.Teardown()
	os.Exit(code)
}
