package log

import (
	"fmt"
	syslog "log"
	"os"
)

// Log is the type that implements log functionality
type Log struct {
	errors   chan error
	filePath string
}

// New method that constructs a new instance of log
func New(filePath string) *Log {
	log := Log{make(chan error, 100), filePath}
	go log.discover()

	return &log
}

func (log *Log) discover() {
	for e := range log.errors {
		log.writeToFile(e)
		syslog.Println(e.Error())
	}
}

func (log *Log) writeToFile(e error) {
	f, err := os.OpenFile(log.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		syslog.Fatal(err)
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%s\n", e.Error()))
	f.Sync()
}

// Error method that writes an error level message to log
func (log *Log) Error(err error) {
	log.errors <- err
}
