package stream

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/tealeg/xlsx/v3"

	log "github.com/sirupsen/logrus"
)

type xlsxStream struct {
	workBook *xlsx.File

	sheets map[string]*xlsxSheet
}

type xlsxSheet struct {
	idx   int
	sheet *xlsx.Sheet
	name  string
}

type xlsxRow struct {
	idx    int
	empty  bool
	values map[int]string
}
type xlsxRows chan *xlsxRow

type Row struct {
	idx    int
	Values map[string]string
}
type Rows chan *Row

type xlsxHeader struct {
	maxIdx int
	values map[int]string
}

func (x *xlsxHeader) Values() []string {
	values := []string{}
	for i := 0; i <= x.maxIdx; i++ {
		values = append(values, x.values[i])
	}
	return values
}

func Open(reader io.Reader) (*xlsxStream, error) {
	stream := &xlsxStream{
		sheets: map[string]*xlsxSheet{},
	}
	if body, err := ioutil.ReadAll(reader); err != nil {
		return stream, fmt.Errorf("ioutil.ReadAll: %v", err)
	} else {
		if wb, err := xlsx.OpenBinary(body); err != nil {
			return stream, fmt.Errorf("xlsx.OpenBinary: %v", err)
		} else {
			stream.workBook = wb
		}
	}

	for sh := range stream.Sheets() {
		stream.sheets[sh.name] = sh
	}

	return stream, nil
	// xlsxFile, _ := xlsx.Open(file.Name, wb)
}

func (x *xlsxStream) Sheets() <-chan *xlsxSheet {
	chnl := make(chan *xlsxSheet)
	go func() {
		for idx, sh := range x.workBook.Sheets {
			chnl <- &xlsxSheet{
				idx:   idx,
				sheet: sh,
				name:  sh.Name,
			}
		}
		close(chnl)
	}()
	return chnl
}

func (x *xlsxStream) Sheet(name string) (*xlsxSheet, error) {
	sh, ok := x.sheets[name]
	if !ok {
		return nil, fmt.Errorf("sheet [%s] not found", name)
	}
	return sh, nil
}

func (x *xlsxStream) Table(name string) (*Table, error) {
	table := &Table{}
	sh, err := x.Sheet(name)
	if err == nil {
		table.header = sh.Header()
		table.rows = sh.Rows()
	}
	return table, err
}

func (x *xlsxSheet) Header() *xlsxHeader {
	hdr := &xlsxHeader{
		maxIdx: -1,
		values: map[int]string{},
	}
	row, _ := x.sheet.Row(0)
	cellVisitor := func(c *xlsx.Cell) error {
		x, y := c.GetCoordinates()
		value, err := c.FormattedValue()
		if err != nil {
			log.Errorf("(%d, %d) %v", x, y, err.Error())
		} else {
			valueString := fmt.Sprint(value)
			hdr.values[x] = valueString
			if x > hdr.maxIdx {
				hdr.maxIdx = x
			}
		}
		return err
	}
	if err := row.ForEachCell(cellVisitor); err != nil {
		//
	}
	return hdr
}

func (x *xlsxSheet) Rows() xlsxRows {
	chnl := make(chan *xlsxRow)
	go func() {
		rowVisitor := func(r *xlsx.Row) error {
			rowIdx := r.GetCoordinate()
			if rowIdx == 0 {
				// skip header row
				return nil
			}

			row := &xlsxRow{
				idx:    rowIdx,
				empty:  true,
				values: map[int]string{},
			}
			//log.Printf("(%d)", rowIdx)
			cellVisitor := func(c *xlsx.Cell) error {
				x, y := c.GetCoordinates()
				value, err := c.FormattedValue()
				// log.Printf("(%d, %d) %v %v", x, y, value, err)
				if err != nil {
					log.Errorf("(%d, %d) %v", x, y, err)
				} else {
					valueString := fmt.Sprint(value)
					if valueString != "" {
						row.values[x] = valueString
						row.empty = false
					}
				}
				return err
			}
			err := r.ForEachCell(cellVisitor)
			if !row.empty {
				chnl <- row
			}
			return err
		}

		x.sheet.ForEachRow(rowVisitor)
		close(chnl)
	}()
	return chnl
}

/*
func (s *XlsxSheet) Rows() <-chan *map[string]string {
	chnl := make(chan *map[string]string)
	go func() {
		for rowIdx := 0; rowIdx < s.MaxRow; rowIdx++ {
			item := map[string]string{}

			r, _ := s.Data[rowIdx+1]
			for hk, hv := range s.headersIndexes {
				v, _ := r[hv]
				item[hk] = v
			}
			chnl <- &item
		}
		close(chnl)
	}()
	return chnl
}

func (s *XlsxSheet) Headers() []string {
	headers := []string{}
	for idx := 0; idx < len(s.headersNames); idx++ {
		headers = append(headers, s.headersNames[idx])
	}
	return headers
}

func Open(name string, input *xlsx.File) (*XlsxFile, error) {
	log.WithFields(log.Fields{
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
					log.Errorf("(%d, %d) %v", x, y, err.Error())
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
*/
