package gccsvlog_test

import (
	"testing"

	fu "github.com/omar391/go-commons/pkg/file"
	gclog "github.com/omar391/go-commons/pkg/log"
	lu "github.com/omar391/go-commons/pkg/log/custom/csv"
)

func TestNewCsvLogger(t *testing.T) {
	type args struct {
		headers []string
		path    string
	}

	type ts struct {
		name string
		args args
		fns  func(*testing.T, lu.CsvLogger, args)
	}

	csvPath := "./tmp/test.csv"
	errFmt := "expected %v, got %v"

	testCases := []ts{
		{
			name: "writes headers to test.csv log file correctly",
			args: args{
				path:    csvPath,
				headers: []string{"My MSG", "Your MSG"},
			},
			fns: func(t *testing.T, l lu.CsvLogger, a args) {
				c, err := fu.ContainsAllTexts(l.FileName(), 1, 1, a.headers...)
				if err != nil {
					t.Error(err)
				}
				if !c {
					t.Errorf(errFmt, true, c)
				}
			},
		},
		{
			name: "writes msg to test.csv log file correctly without the csv indicator",
			args: args{
				path:    csvPath,
				headers: []string{"My MSG 2", "Your MSG 2"},
			},
			fns: func(t *testing.T, l lu.CsvLogger, a args) {
				c, err := fu.ContainsAllTexts(l.FileName(), 1, 1, a.headers...)
				if err != nil {
					t.Error(err)
				}
				if !c {
					t.Errorf(errFmt, true, c && csvPath != a.path)
				}
			},
		},
		{
			name: "writes doesn't write the header if the headers exist in the file",
			args: args{
				path:    csvPath,
				headers: []string{"My MSG 3", "Your MSG 3"},
			},
			fns: func(t *testing.T, l lu.CsvLogger, a args) {
				l.Csv("hello", "world", "ssd", "{json: true, msg: \"dd\"}")
				var err error
				c, err := fu.ContainsAllTexts(l.FileName(), 2, 0, "hello", "world")
				if err != nil {
					t.Error(err)
				}
				cc, err := fu.ContainsAllTexts(l.FileName(), 2, 0, "My MSG 3", "Your MSG 3")
				if err != nil {
					t.Error(err)
				}
				if !c {
					t.Errorf(errFmt, true, c && !cc)
				}
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lu.New(tt.args.path, tt.args.headers, nil, nil)
			got.UpdateBaseLogger(got.SetMinCallerAttachLevel(gclog.DebugLevel))
			if err != nil {
				t.Error(err)
			}
			tt.fns(t, got, tt.args)
		})
	}
}
