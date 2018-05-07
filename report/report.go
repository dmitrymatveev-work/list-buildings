package report

import (
	"fmt"
	"list-buildings/log"
	"list-buildings/model"
	"os"
)

// Report is the type that implements reporting functionality
type Report struct {
	filePath string
	log      *log.Log
}

// New method that constructs a new instance of report
func New(filePath string, log *log.Log) *Report {
	return &Report{filePath, log}
}

// Write method writes an entry to output
func (r *Report) Write(b *model.Building) {
	f, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		r.log.Error(err)
		panic(err)
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%s,%s,%s\n", b.Street, b.Building, b.URL))
	f.Sync()
}
