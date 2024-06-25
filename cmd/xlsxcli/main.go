package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	logger "github.com/mguzelevich/go-ext-log"
	log "github.com/sirupsen/logrus"
	xlsxLib "github.com/tealeg/xlsx/v3"

	ts "github.com/mguzelevich/storages/table_storage"
	"github.com/mguzelevich/xlsx"
	//	"github.com/mguzelevich/xlsx/csv"
)

var (
	Log = log.New()

	appStartedAt = time.Now()

	logLevel = flag.String("log-level", "info", "log level: []")
	logFile  = flag.String("log-file", "", "log file")

	outputPrefix      string
	appendMetaColumns bool
	mode              string
	mappingFile       string

	fields  []string
	mapping map[string]string
)

type sheetDesc struct {
	space   string
	name    string
	headers []string
}

func initFieldMapping() ([]string, map[string]string) {
	fields := []string{}
	fieldsMapping := map[string]string{}
	if mappingFile != "" {
		file, err := os.Open(mappingFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			field := ""
			for i, item := range strings.Split(scanner.Text(), ",") {
				if i == 0 {
					field = item
					fields = append(fields, field)
				}
				fieldsMapping[item] = field
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return fields, fieldsMapping
}

func init() {
	flag.StringVar(&outputPrefix, "out-prefix", "./", "output file prefix")
	flag.StringVar(&mode, "mode", "one2many", "processing mode: [one2one|one2many]")
	flag.StringVar(&mappingFile, "mapping", "", "headers mapping file")
	flag.BoolVar(&appendMetaColumns, "append-meta", false, "add columns with metainfo (source sheet, source row number, etc.)")

	flag.Parse()

	level, _ := log.ParseLevel(*logLevel)
	Log, _ = logger.Init(level, *logFile)
	logger.Apply(Log)
	log.Debugf("level: %v", level)

	os.Args = flag.Args()

	log.Infof(
		"out-prefix=%v log=%v mode=%v mapping=%v append-meta=%v",
		outputPrefix, logFile, mode, mappingFile, appendMetaColumns,
	)
}

func one2many(db *ts.Database, sheetChan chan *sheetDesc, rowChan chan *ts.Row, shutdownChan chan bool) {
	var sheet *ts.Table
	for {
		select {
		case sh := <-sheetChan:
			log.WithFields(log.Fields{
				"sheet": sh,
			}).Debugf("sheet msg received")
			sheet = db.NewTable(fmt.Sprintf("%s / %s", sh.space, sh.name))
			sheet.Init(ts.FieldsFromStrings(sh.headers)...)
		case r := <-rowChan:
			log.WithFields(log.Fields{
				"row": r,
			}).Debugf("row msg received")
			sheet.Insert(r)
			// case <-shutdownChan:
			// 	log.WithFields(log.Fields{
			// 		// "row": r,
			// 	}).Infof("shutdown msg received")
			// 	break
		}
	}
}

func one2one(db *ts.Database, sheetChan chan *sheetDesc, rowChan chan *ts.Row, shutdownChan chan bool) {
	sheet := db.NewTable("table")
	sheet.Init(ts.FieldsFromStrings([]string{"_space", "_table", "_idx"})...)
	var tblRow *ts.Row
	for {
		select {
		case sh := <-sheetChan:
			sheet.AppendFields(ts.FieldsFromStrings(sh.headers)...)
			tblRow = &ts.Row{"", ts.Record{"_space": sh.space, "_table": sh.name}}
		case r := <-rowChan:
			log.WithFields(log.Fields{
				"row": r,
			}).Debugf("row msg received")

			r.Join(tblRow)
			sheet.Insert(r)
			// case <-shutdownChan:
			// 	log.WithFields(log.Fields{
			// 		// "row": r,
			// 	}).Infof("shutdown msg received")
			// 	break
		}
	}
}

func _tst() error {
	/*
	p.s.Init(params.GdriveCredentials, params.GdriveToken)

	if db, err := sql.Open("sqlite", "database.s3db"); err != nil {
		log.Fatal(err)
	} else {
		p.DB = db
	}

	sh, _ := p.s.Sheets.ReadSpreadsheet(file.ID, true)

	if _, err := p.DB.Exec("drop VIEW IF EXISTS source_data;"); err != nil {
		log.Error(err)
		return err
	}

	for _, table := range sh {
		stmt := table.CreateTableStmt()
		if _, err := p.DB.Exec(stmt); err != nil {
			log.Debug(stmt)
			log.Error(err)
			return err
		}

		istmt := table.InsertAllStmt()
		if _, err := p.DB.Exec(istmt); err != nil {
			log.Debug(istmt)
			log.Error(err)
			return err
		}
	}

	stmt := p.mld.CreateSourceViewStmt()
	if _, err := p.DB.Exec(stmt); err != nil {
		log.Debug(stmt)
		log.Error(err)
		return err
	}
	*/
	return nil
}

func main() {
	args := os.Args
	if len(args) <= 0 {
		log.Fatalf("empty input files list")
	}

	log.Infof("xlsxcli started [%v], %v", appStartedAt, args)

	fields, mapping = initFieldMapping()

	db, _ := ts.New("db")

	idx := 0
	for _, input := range args {
		_, inputFilename := filepath.Split(input)
		tablespaceName := inputFilename
		if strings.HasSuffix(inputFilename, ".xlsx") {
			tablespaceName = inputFilename[:len(inputFilename)-5]
		}

		// body, _ := ioutil.ReadAll(f)
		// wb, _ := xlsxLib.OpenBinary(body)
		wb, err := xlsxLib.OpenFile(input)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
				// "path": file.Path,
				// "name": fileName,
				// "id": file.ID,
			}).Fatalf("xlsxLib.OpenFile")
		}
		xlsxFile, _ := xlsx.Open(tablespaceName, wb)
		// ioutil.WriteFile("/tmp/result.xlsx", body, 0600)

		sheetChan := make(chan *sheetDesc)
		rowChan := make(chan *ts.Row)
		shutdownChan := make(chan bool)

		if mode == "one2many" {
			go one2many(db, sheetChan, rowChan, shutdownChan)
		} else if mode == "one2one" {
			go one2one(db, sheetChan, rowChan, shutdownChan)
		}

		for s := range xlsxFile.Sheets() {
			sd := &sheetDesc{tablespaceName, s.Name, s.Headers()}
			if len(fields) > 0 {
				sd.headers = fields
			}
			sheetChan <- sd

			log.WithFields(log.Fields{"sheet": s}).Debugf("sheet")

			firstRow := true
			rIdx := 0
			for row := range s.Rows() {
				rIdx++
				if firstRow {
					firstRow = false
					continue
				}
				rec := ts.Record{"_idx": fmt.Sprintf("%d", rIdx)}
				for k, v := range *row {
					field := ts.Field(k)
					if m, ok := mapping[k]; ok {
						field = ts.Field(m)
					}
					rec[field] = v
				}
				rowChan <- &ts.Row{fmt.Sprintf("%d", idx), rec}
				idx++
			}
			log.WithFields(log.Fields{
				// "path": file.Path,
				// "headers": sheet.Headers(),
				// "dump":    sheet.DumpAsJson(),
			}).Debugf("sheet processed")

		}
		close(shutdownChan)

		// db.DumpJson()
		log.WithFields(log.Fields{
			// "path": file.Path,
			// "name": fileName,
			// "id": file.ID,
		}).Debugf("file processed")
	}

	if mode == "one2many" {
		for sheet := range db.Tables() {
			sheet.DumpAsCsv()
		}
	} else if mode == "one2one" {
		for sheet := range db.Tables() {
			sheet.DumpAsCsv()
		}
	}
	/*
		if mode == "one2many" {
			// // store.SaveChildren("output")
			tablesPathes, _ := store.Tables()
			for _, path := range tablesPathes {
				headers, data, _ := store.Table(path)
				outPath := ""
				outPath = fmt.Sprintf("%s%s.csv", outputPrefix, strings.Join(path, "."))
				log.Infof("one2many out [%s]", outPath)
				csv.Writer(outPath, headers, data)
			}
			// headers, data := store.Tables("output.csv")
		} else if mode == "one2one" {
			headers := []string{}
			hMap := map[string]int{}
			data := []map[string]string{}
			tablesPathes, _ := store.Tables()
			for _, path := range tablesPathes {
				h, d, _ := store.Table(path)
				for _, hi := range h {
					if _, ok := hMap[hi]; !ok {
						hMap[hi] = len(hMap)
						headers = append(headers, hi)
					}
				}
				data = append(data, d...)
			}
			outPath := ""
			//outPath = fmt.Sprintf("%s%s.csv", outputPrefix, strings.Join(path, "."))
			log.Infof("one2one out [%s]", outPath)
			headers, data, _ = store.ReMapData(headers, data)
			csv.Writer(outPath, headers, data)
		} else {
			store.DumpJson("output.csv")
		}
	*/
	appFinishedAt := time.Now()
	appDuration := int64(appFinishedAt.Sub(appStartedAt) / time.Millisecond)
	log.Infof("xlsxcli finished %d ms", appDuration)
}
