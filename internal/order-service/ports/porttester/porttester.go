package porttester

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
)

type PortTester[T any] struct {
	Subject   T
	PostCheck func(t testing.TB)
}

func Noop(testing.TB) {}

func VerifyNoInUseConnections(db *sqlx.DB) func(t testing.TB) {
	//preStats := db.Stats()
	return func(t testing.TB) {
		t.Helper()
		gomega.NewWithT(t).Eventually(db.Stats).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{"InUse": gomega.Equal(0)}))
	}
}
