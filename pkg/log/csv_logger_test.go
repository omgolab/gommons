package log_utils_test

import (
	"testing"

	fu "github.com/omar391/go-commons/pkg/file"
	lu "github.com/omar391/go-commons/pkg/log"
)

func TestNewCsvLogger(t *testing.T) {
	type args struct {
		headers []string
		path    string
	}

	tests := []struct {
		name string
		args args
		fns  func(*testing.T, *lu.CsvLogger, args)
	}{
		{
			name: "writes headers to test.csv log file correctly",
			args: args{
				path:    "./tmp/test.csv",
				headers: []string{"My MSG", "Your MSG"},
			},
			fns: func(t *testing.T, l *lu.CsvLogger, a args) {
				c, err := fu.ContainsAllTexts(l.FileName(), 1, 1, a.headers...)
				if err != nil {
					t.Error(err)
				}
				if !c {
					t.Errorf("expected %v, got %v", true, c)
				}
				l.Close()
			},
		},
		{
			name: "writes msg to test.csv log file correctly without the csv indicator",
			args: args{
				path:    "./tmp/test.csv",
				headers: []string{"My MSG 2", "Your MSG 2"},
			},
			fns: func(t *testing.T, l *lu.CsvLogger, a args) {
				c, err := fu.ContainsAllTexts(l.FileName(), 1, 1, a.headers...)
				if err != nil {
					t.Error(err)
				}
				if !c {
					t.Errorf("expected %v, got %v", true, c && l.FileName() != a.path)
				}
				l.Close()
			},
		},
		{
			name: "writes doesn't write the header if the headers exist in the file",
			args: args{
				path:    "./tmp/test.csv",
				headers: []string{"My MSG 3", "Your MSG 3"},
			},
			fns: func(t *testing.T, l *lu.CsvLogger, a args) {
				l.Csv("hello", "world", lu.CsvIndicator, "ssd", "{json: true, msg: \"dd\"}")
				var err error
				c, err := fu.ContainsAllTexts(l.FileName(), 2, 0, "hello", "world")
				if err != nil {
					t.Error(err)
				}
				cc, err := fu.ContainsAllTexts(l.FileName(), 2, 0, lu.CsvIndicator)
				if err != nil {
					t.Error(err)
				}
				if !c {
					t.Errorf("expected %v, got %v", true, c && !cc)
				}
				l.Close()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lu.NewCsvLogger(tt.args.path, tt.args.headers)
			if err != nil {
				t.Error(err)
			}
			tt.fns(t, got, tt.args)
		})
	}
}
