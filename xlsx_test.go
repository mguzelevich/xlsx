package xlsx

import (
	// "fmt"
	"testing"

	logger "github.com/mguzelevich/go-ext-log"
	log "github.com/sirupsen/logrus"
)

func initData() *XlsxFile {
	return nil
}

func TestS(t *testing.T) {
	//t.Fatalf("init %v", db)
	// t1 := db.Table("t1")
	// for row := range t1.Select() {
	// 	t.Logf("%v", row)
	// }
	// if err != nil {
	// 	t.Fatal("Unable to format entry: ", err)
	// }
}

func init() {
	logger.Init(log.TraceLevel, "")
	//logger.Init(log.FatalLevel, "")
}
