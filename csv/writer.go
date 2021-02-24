package csv

import (
	"encoding/csv"
	"os"

	log "github.com/sirupsen/logrus"
)

func Writer(output string, headers []string, data []map[string]string) error {
	writer := os.Stdout

	//flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if output != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if file, err := os.OpenFile(output, flags, 0666); err != nil {
			log.Info("Failed to log to file, using default stdout")
		} else {
			writer = file
		}
	}

	w := csv.NewWriter(writer)

	if err := w.Write(headers); err != nil {
		log.Fatalln("error writing headers to csv:", err)
	}
	for _, d := range data {
		record := []string{}
		for _, field := range headers {
			record = append(record, d[field])
		}
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return nil
}
