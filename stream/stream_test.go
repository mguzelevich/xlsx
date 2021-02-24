package stream

import (
	"os"
	// "io"
	// "io/ioutil"
	"testing"
	// "github.com/tealeg/xlsx/v3"
)

func TestOpenStream(t *testing.T) {
	//flags := os.O_CREATE | os.O_WRONLY //| os.O_APPEND
	// file, err := os.OpenFile("../samples/test-00.xlsx", flags, 0666)
	file, err := os.Open("../samples/test-00.xlsx")
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	//t.Fatalf("init %v", db)
	// t1 := db.Table("t1")
	// for row := range t1.Select() {
	// 	t.Logf("%v", row)
	// }
	// if err != nil {
	// 	t.Fatal("Unable to format entry: ", err)
	// }
	stream, err := Open(file)
	t.Logf("stream %v %v", stream, err)
}

func TestSheets(t *testing.T) {
	//flags := os.O_CREATE | os.O_WRONLY //| os.O_APPEND
	// file, err := os.OpenFile("../samples/test-00.xlsx", flags, 0666)
	file, err := os.Open("../samples/test-00.xlsx")
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	stream, err := Open(file)
	t.Logf("stream %v %v", stream, err)

	for sh := range stream.Sheets() {
		t.Logf("sheet %v", sh)
	}

}

func TestSheet(t *testing.T) {
	//flags := os.O_CREATE | os.O_WRONLY //| os.O_APPEND
	// file, err := os.OpenFile("../samples/test-00.xlsx", flags, 0666)
	file, err := os.Open("../samples/test-00.xlsx")
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	stream, err := Open(file)
	t.Logf("stream %v %v", stream, err)

	sh, err := stream.Sheet("users")
	t.Logf("sheet %v %v", sh, err)

	sh, err = stream.Sheet("users111")
	t.Logf("sheet %v %v", sh, err)

}
func TestHeaders(t *testing.T) {
	//flags := os.O_CREATE | os.O_WRONLY //| os.O_APPEND
	// file, err := os.OpenFile("../samples/test-00.xlsx", flags, 0666)
	file, err := os.Open("../samples/test-00.xlsx")
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	stream, err := Open(file)
	t.Logf("stream %v %v", stream, err)

	for sh := range stream.Sheets() {
		t.Logf("sheet %v - %v", sh, sh.Header())
	}
}

func TestRows(t *testing.T) {
	//flags := os.O_CREATE | os.O_WRONLY //| os.O_APPEND
	// file, err := os.OpenFile("../samples/test-00.xlsx", flags, 0666)
	file, err := os.Open("../samples/test-00.xlsx")
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	stream, err := Open(file)
	t.Logf("stream %v %v", stream, err)

	for sh := range stream.Sheets() {
		t.Logf("sheet %v - %v", sh, sh.Header())
		for r := range sh.Rows() {
			t.Logf("row %v", r)
		}
	}
}

func TestTables(t *testing.T) {
	//flags := os.O_CREATE | os.O_WRONLY //| os.O_APPEND
	// file, err := os.OpenFile("../samples/test-00.xlsx", flags, 0666)
	file, err := os.Open("../samples/test-00.xlsx")
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	stream, err := Open(file)
	t.Logf("stream %v %v", stream, err)

	for sh := range stream.Sheets() {
		t.Logf("sheet %v - %v", sh, sh.Header())
		tbl, err := stream.Table(sh.name)
		t.Logf("table %v %v %v", tbl, tbl.header, err)
		for r := range tbl.Rows() {
			t.Logf("row %v", r)
		}
	}
}

// func init() {
// 	logger.Init(log.TraceLevel, "")
// 	//logger.Init(log.FatalLevel, "")
// }
