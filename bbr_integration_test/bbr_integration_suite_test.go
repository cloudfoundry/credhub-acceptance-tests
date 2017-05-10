package bbr_integration

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"
)

func TestBbrIntegrationTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Backup and Restore integration suite")
}

var uniqueTestID string
var jumpBoxSession *JumpBoxSession

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(15 * time.Minute)
	uniqueTestID = strconv.FormatInt(time.Now().UnixNano(), 16)
	jumpBoxSession = NewJumpBoxSession(uniqueTestID)
})

var _ = AfterSuite(func() {
	jumpBoxSession.Cleanup()
})
