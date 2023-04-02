package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

type UILogger struct {
	out    *bytes.Buffer
	logger *log.Logger
}

func (l *UILogger) alloc() {
	l.out = bytes.NewBuffer(nil)
}

func NewLogger() *UILogger {
	l := &UILogger{}
	l.alloc()

	l.logger = log.New(l.out, "", log.LstdFlags)

	return l
}

func (l *UILogger) Log(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *UILogger) Dump() error {
	fmt.Fprintf(os.Stdout, "<< Begin loggerUI output >>\n")
	fmt.Fprintf(os.Stdout, l.out.String())
	fmt.Fprintf(os.Stdout, "<< End loggerUI output >>\n")
	return nil
}
