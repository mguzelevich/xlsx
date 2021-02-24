package xlsx

import (
	"fmt"

	logger "github.com/mguzelevich/go-ext-log"
	log "github.com/sirupsen/logrus"

	"github.com/tealeg/xlsx/v3"
)

func Open(name string, input *xlsx.File) (*XlsxFile, error) {
	logger.Log.WithFields(log.Fields{
		"name": name,
	}).Infof("open xlsx file")

	output := &XlsxFile{
		sheetsIndexes: map[string]int{},

		Name:   name,
		sheets: []*XlsxSheet{},
	}

	for idx, sh := range input.Sheets {
		sheet := output.AddSheet(idx, sh.Name)

		rowVisitor := func(r *xlsx.Row) error {
			rowIdx := r.GetCoordinate()
			cellVisitor := func(c *xlsx.Cell) error {
				x, y := c.GetCoordinates()
				value, err := c.FormattedValue()
				if err != nil {
					logger.Log.Errorf("(%d, %d) %v", x, y, err.Error())
				} else {
					valueString := fmt.Sprint(value)
					if rowIdx == 0 {
						if valueString == "" {
							valueString = fmt.Sprintf("col_%02d", x)
						}
						sheet.headersNames[x] = valueString
						sheet.headersIndexes[valueString] = x
					} else {
						if valueString != "" {
							sheet.Set(x, y, valueString)
						}
					}
				}
				return err
			}
			err := r.ForEachCell(cellVisitor)
			return err
		}

		sh.ForEachRow(rowVisitor)
		// log.Infof("- %v -> [%s] (%d rows ~~ %d rows)", sh.Name, output, sh.MaxRow, lastNotEmptyRow)
	}

	return output, nil
}
