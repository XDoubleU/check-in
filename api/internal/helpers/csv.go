package helpers

import (
	"fmt"
	"net/http"

	"github.com/gocarina/gocsv"
)

func WriteCSV(w http.ResponseWriter, filename string, data any) error {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().
		Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.csv", filename))

	return gocsv.Marshal(data, w)
}
