package xlsx

import (
	"fmt"
)

type XlsxFile struct {
	Name          string
	sheetsIndexes map[string]int
	sheets        []*XlsxSheet
}

func (x *XlsxFile) AddSheet(idx int, name string) *XlsxSheet {
	sheet := &XlsxSheet{
		idx:            idx,
		headersNames:   map[int]string{},
		headersIndexes: map[string]int{},

		Name:      name,
		Data:      map[int]map[int]string{},
		MaxRow:    -1,
		MaxColumn: -1,
	}

	x.sheets = append(x.sheets, sheet)
	x.sheetsIndexes[name] = len(x.sheets) - 1
	return sheet
}

func (x *XlsxFile) Sheet(name string) *XlsxSheet {
	return x.sheets[x.sheetsIndexes[name]]
}

func (x *XlsxFile) Sheets() <-chan *XlsxSheet {
	chnl := make(chan *XlsxSheet)
	go func() {
		for _, sh := range x.sheets {
			chnl <- sh
		}
		close(chnl)
	}()
	return chnl
}

func (x *XlsxFile) String() string {
	return fmt.Sprintf("%s %v", x.Name, x.sheetsIndexes)
}
