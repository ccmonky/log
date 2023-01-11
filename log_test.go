package log_test

import (
	stdlog "log"
	"testing"

	"github.com/ccmonky/log"
)

func TestPrintln(t *testing.T) {
	stdlog.Println(log.ErrorLevel, 1, 3, 4)
	log.Println(log.ErrorLevel, 1, 3, 4)
	t.Fatal(1)
}
