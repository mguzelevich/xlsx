package csv

import (
	"encoding/csv"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func Reader(filename string) ([]map[string]string, error) {
	logger := log.WithFields(log.Fields{
		"func":     "csv reader",
		"filename": filename,
	})

	csvFile, err := os.Open(filename)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Debugf("CSV file opened")
	defer csvFile.Close()

	//r := csv.NewReader(strings.NewReader(in))
	r := csv.NewReader(csvFile)

	data := []map[string]string{}

	fileHeaders := []string{}
	idx := -1
	for {
		idx++

		record, err := r.Read()
		if err == io.EOF {
			logger.Debugf("eof")
			break
		}
		if err != nil {
			logger.WithFields(log.Fields{
				"err": err,
			}).Debugf("error")
			break
		}
		if idx == 0 {
			// process headers
			for _, h := range record {
				fileHeaders = append(fileHeaders, h)
			}
			continue
		}
		rec := map[string]string{}
		for i, value := range record {
			rec[fileHeaders[i]] = value
		}

		logger.WithFields(log.Fields{
			"idx": idx,
			"rec": rec,
		}).Debugf("readed")
		data = append(data, rec)
	}
	return data, nil
}
