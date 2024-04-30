package glog_test

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"testing"

	glog "github.com/omgolab/go-commons/pkg/log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func Test_logCfg_Event(t *testing.T) {
	l, _ := glog.New(glog.WithFileLogger("test.log"))
	err := fmt.Errorf("from error: %w", errors.New("error message"))
	l.Error("error msg", err)
}

func TestOriginalLogStack(t *testing.T) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	out := &bytes.Buffer{}
	log := zerolog.New(out)

	err := fmt.Errorf("from error: %w", errors.New("error message"))
	log.Log().Stack().Err(err).Msg("")

	got := out.String()
	want := `\{"stack":\[\{"func":"TestLogStack","line":"27","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}
