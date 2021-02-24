package xlsx

import (
	"fmt"
)

type XlsxRow struct {
	idx int
	// (col) = value
	Data map[int]string
}

type XlsxSheet struct {
	idx            int
	Name           string
	headersNames   map[int]string
	headersIndexes map[string]int

	// (row, col) = value
	Data      map[int]map[int]string
	MaxRow    int
	MaxColumn int
}

func (s *XlsxSheet) Set(x int, y int, value string) {
	if s.MaxColumn < x {
		s.MaxColumn = x
	}
	if s.MaxRow < y {
		s.MaxRow = y
	}
	if _, ok := s.Data[y]; !ok {
		s.Data[y] = map[int]string{}
	}
	s.Data[y][x] = value
}

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

func (s *XlsxSheet) String() string {
	return fmt.Sprintf("%s - %v rows %v", s.Name, s.MaxRow, s.headersNames)
}
