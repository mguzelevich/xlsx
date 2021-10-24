package main

import (
	"bufio"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	logger "github.com/mguzelevich/go-ext-log"
	log "github.com/sirupsen/logrus"
	xlsxLib "github.com/tealeg/xlsx/v3"

	"github.com/mguzelevich/xlsx"
	//	"github.com/mguzelevich/xlsx/csv"
)

var (
	Log = log.New()

	appStartedAt = time.Now()

	logLevel = flag.String("log-level", "fatal", "log level: []")
	logFile  = flag.String("log-file", "", "log file")

	outputPrefix      string
	appendMetaColumns bool
	mode              string
	mappingFile       string
)

func initFieldMapping() map[string]string {
	fieldsMapping := make(map[string]string)
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
				}
				fieldsMapping[item] = field
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return fieldsMapping
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
	Log.Infof("level: %v", level)

	os.Args = flag.Args()

	Log.Infof(
		"out-prefix=%v log=%v mode=%v mapping=%v append-meta=%v",
		outputPrefix, logFile, mode, mappingFile, appendMetaColumns,
	)
}

func main() {
	args := os.Args
	if len(args) <= 0 {
		log.Fatalf("empty input files list")
	}

	Log.Infof("xlsxcli started [%v], %v", appStartedAt, args)

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
			Log.WithFields(log.Fields{
				"err": err,
				// "path": file.Path,
				// "name": fileName,
				// "id": file.ID,
			}).Fatalf("xlsxLib.OpenFile")
		}
		xlsxFile, _ := xlsx.Open(tablespaceName, wb)
		// ioutil.WriteFile("/tmp/result.xlsx", body, 0600)
		// msgs := xlsxFile.Sheet("messages")
		// users := xlsxFile.Sheet("users")
		// for msg := range msgs.Rows() {
		// 	// _startedAt := time.Now()
		// 	p.processCampaign(file, *msg, users)

		// }

		for s := range xlsxFile.Sheets() {
			Log.WithFields(log.Fields{
				"sheet": s,
			}).Infof("sheet")
		}
		Log.WithFields(log.Fields{
			// "path": file.Path,
			// "name": fileName,
			// "id": file.ID,
		}).Infof("file processed")
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
	Log.Infof("xlsxcli finished %d ms", appDuration)
}
