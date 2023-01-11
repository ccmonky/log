package log_test

import (
	stdlog "log"
	"testing"

	"github.com/ccmonky/log"
)

func TestPrintln(t *testing.T) {
	log.Println(1, 3, 4)
	log.Println = stdlog.Println
	t.Fatal(1)
}
