package helpers

import (
	"encoding/csv"
	"fmt"
	"net/http"
)

func WriteCSV(w http.ResponseWriter, filename string, data [][]string) error {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().
		Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.csv", filename))

	csvWriter := csv.NewWriter(w)

	return csvWriter.WriteAll(data)
}
