package stream

type Table struct {
	header *xlsxHeader
	rows   xlsxRows
}

func (t *Table) Rows() Rows {
	chnl := make(Rows)
	go func() {
		for r := range t.rows {
			//log.Printf("Rows t.r %v", r)
			row := &Row{
				idx:    r.idx,
				Values: map[string]string{},
			}
			for k, v := range r.values {
				row.Values[t.header.values[k]] = v
			}
			//log.Printf("Rows r %v", r)
			chnl <- row
		}
		close(chnl)
	}()
	return chnl
}
