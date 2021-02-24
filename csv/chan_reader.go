package csv

import (
	"encoding/csv"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// ChanReader - async csv readed
// function return headers
// records and errors sended to channel
func ChanReader(filename string, records chan<- map[string]string, errs chan<- error, shutdown <-chan bool) ([]string, error) {
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

	fileHeaders := make(chan []string)
	go func() {
		headers := []string{}
		idx := -1
		for {
			idx++

			record, err := r.Read()
			if err == io.EOF {
				logger.Debugf("eof")
				errs <- io.EOF
				break
			}
			if err != nil {
				logger.WithFields(log.Fields{
					"err": err,
				}).Debugf("error")
				errs <- err
				break
			}
			if idx == 0 {
				// process headers
				for _, h := range record {
					headers = append(headers, h)
				}
				fileHeaders <- headers
				continue
			}
			rec := map[string]string{}
			for i, value := range record {
				rec[headers[i]] = value
			}

			logger.WithFields(log.Fields{
				"idx": idx,
				"rec": rec,
			}).Debugf("readed")
			records <- rec
		}
	}()
	fh := <-fileHeaders
	close(fileHeaders)
	return fh, nil
}
