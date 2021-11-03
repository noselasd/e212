package routes

import (
	"e212/store"
	"encoding/csv"
)

func writeRow(e *store.E212Entry, wr *csv.Writer) error {
	row := []string{e.E212Country.Name, e.Operator, e.Code.Mcc, e.Code.Mnc}
	return wr.Write(row)
}

func csvExport(ctx *AppContext) {
	entries, err := store.E212GetAll()
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	separator := ','
	reqSep := ctx.Query("separator")
	if len(reqSep) == 1 {
		separator = rune(reqSep[0])

	}

	csvWr := csv.NewWriter(ctx.Resp)
	csvWr.Comma = separator

	ctx.Resp.Header().Set("content-disposition", "attachment; filename=\"e212.csv\"")
	ctx.Resp.Header().Set("content-type", "text/csv")
	defer csvWr.Flush()

	csvHeaders := []string{"Country", "Operator", "MCC", "MNC"}
	csvWr.Write(csvHeaders)
	for i := range entries {
		err := writeRow(&entries[i], csvWr)
		if err != nil {
			ctx.Logger.Println("csvexport error:", err.Error())
			break
		}
	}
}
