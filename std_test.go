package log_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/ccmonky/log"
)

func TestLevelLogger(t *testing.T) {
	var cases = []struct {
		level   log.Level
		out     io.Writer
		lines   int
		keyword string
	}{
		{
			log.DebugLevel,
			new(bytes.Buffer),
			3,
			"logger debug",
		},
		{
			log.InfoLevel,
			new(bytes.Buffer),
			2,
			"logger info",
		},
		{
			log.ErrorLevel,
			new(bytes.Buffer),
			1,
			"logger error",
		},
	}
	for _, tc := range cases {
		std := log.New(tc.out, "", log.LstdFlags)
		logger := log.NewLevelLogger(
			log.WithLevel(tc.level),
			log.WithLogger(std),
		)
		logger.Debug("msg 1", "logger", tc.level)
		logger.Info("msg 2", "logger", tc.level)
		logger.Error("msg 3", "logger", tc.level)
		lines := strings.Split(strings.TrimSpace(tc.out.(*bytes.Buffer).String()), "\n")
		if len(lines) != tc.lines {
			for _, line := range lines {
				t.Log(line)
			}
			t.Fatalf("level %v, should == %d, got %d", tc.level, tc.lines, len(lines))
		}
		for _, line := range lines {
			if !strings.Contains(line, tc.keyword) {
				t.Fatalf("should contain %s, but got %s", tc.keyword, line)
			}
		}
	}
}

func TestStd(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	log.SetLogger("", log.NewLevelLogger(log.WithLogger(logger)))
	log.SetLogger("debug", log.NewLevelLogger(log.WithLogger(logger), log.WithLevel(log.D)))
	var cases = []struct {
		ps     []any
		expect string
	}{
		{
			[]any{1, 2, 3},
			"level: info, msg: , 1 2 3\n",
		},
		{
			[]any{log.Ctx("nop"), 1, 2, 3},
			"",
		},
		{
			[]any{log.Ctx("", log.E), 1, 2, 3},
			"level: error, msg: , 1 2 3\n",
		},
		{
			[]any{log.Ctx(log.D), 1, 2, 3},
			"",
		},
		{
			[]any{log.Ctx("debug", log.D), 1, 2, 3},
			"level: debug, msg: , 1 2 3\n",
		},
	}
	for i, tc := range cases {
		log.Println(tc.ps...)
		line := buf.String()
		result := line
		if len(line) > 20 {
			result = line[20:]
		}
		if result != tc.expect {
			t.Fatalf("case %d should == %s(%d), got %s(%d)", i, tc.expect, len(tc.expect), result, len(result))
		}
		buf.Reset()
	}
}
