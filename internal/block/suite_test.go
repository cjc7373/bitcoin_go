package block

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	bolt "go.etcd.io/bbolt"

	"github.com/cjc7373/bitcoin_go/internal/db"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

var testDBPath = "blockchain_test.db"
var testDB *bolt.DB
var testConf utils.Config
var testWalletName = "default"

func TestBlockchain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "blockchain Suite")
}

var _ = BeforeEach(func() {
	testConf = utils.Config{
		DBPath:  testDBPath,
		Wallets: map[string]string{},
	}
	testDB = db.OpenDB(&testConf)
})

var _ = AfterEach(func() {
	testDB.Close()
	os.Remove(testDBPath)
})
