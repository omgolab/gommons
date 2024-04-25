package gcprocess_test

import (
	"bytes"
	"reflect"
	"runtime"
	"testing"

	gcprocess "github.com/omgolab/go-commons/pkg/runtime/process"
	ts "github.com/omgolab/go-commons/pkg/test"
	gomock "go.uber.org/mock/gomock"
)

type testMap = map[ts.GWTSteps]func(t *testing.T, ctrl *gomock.Controller) test

type args struct {
	fields gcprocess.ProcessFilter
}

type test struct {
	args    args
	want    []map[string]string
	wantErr bool
}

func sut(args args) gcprocess.ProcessFilter {
	return args.fields
}

func TestProcessFilter_ListProcesses(t *testing.T) {
	tests := testMap{
		// ts.ForScenario("Mock a process execution PASS for unix").
		// 	Given("a Unix-like system").
		// 	When("listing processes").
		// 	Then("it should return expected process details").
		// 	AndReturnsNoError(): func(t *testing.T, ctrl *gomock.Controller) test {
		// 	var stdout, stderr bytes.Buffer
		// 	mockExecutor := NewMockCommandExecutor(ctrl)
		// 	mockExecutor.EXPECT().ExecuteCommand(&stdout, &stderr, "sh", "-c", "ps aux").DoAndReturn(func(stdout, stderr *bytes.Buffer, command string, args ...string) error {
		// 		stdout.WriteString("USER PID COMMAND\nroot 1 /usr/lib/systemd/systemd")
		// 		return nil
		// 	})

		// 	return test{
		// 		args: args{
		// 			fields: gcprocess.ProcessFilter{
		// 				Attrs:           map[string]string{},
		// 				CommandExecutor: mockExecutor,
		// 			},
		// 		},
		// 		want:    []map[string]string{{"USER": "root", "PID": "1", "COMMAND": "/usr/lib/systemd/systemd"}},
		// 		wantErr: false,
		// 	}
		// },
		ts.ForScenario("Mock a process execution PASS for windows").
			Given("a Windows-like system").
			When("listing processes").
			Then("it should return expected process details").
			AndReturnsNoError(): func(t *testing.T, ctrl *gomock.Controller) test {
			var stdout, stderr bytes.Buffer
			mockExecutor := NewMockCommandExecutor(ctrl)
			if runtime.GOOS == "windows" {
				mockExecutor.EXPECT().ExecuteCommand(&stdout, &stderr, "cmd", "/C", "tasklist").DoAndReturn(func(stdout, stderr *bytes.Buffer, command string, args ...string) error {
					stdout.WriteString("ImageName PID\nSystem 1")
					return nil
				})
			}

			return test{
				args: args{
					fields: gcprocess.ProcessFilter{
						Attrs:           map[string]string{},
						CommandExecutor: mockExecutor,
					},
				},
				want:    []map[string]string{{"ImageName": "System", "PID": "1"}},
				wantErr: false,
			}
		},
		// ts.ForScenario("Mock a process execution FAIL for unix").
		// 	Given("a Unix-like system").
		// 	When("listing processes").
		// 	Then("it should return expected process details").
		// 	AndReturnsError(): func(t *testing.T, ctrl *gomock.Controller) test {
		// 	var stdout, stderr bytes.Buffer
		// 	mockExecutor := NewMockCommandExecutor(ctrl)
		// 	mockExecutor.EXPECT().ExecuteCommand(&stdout, &stderr, "sh", "-c", "ps aux").Return(errors.New("command failed"))

		// 	return test{
		// 		args: args{
		// 			fields: gcprocess.ProcessFilter{
		// 				Attrs:           map[string]string{},
		// 				CommandExecutor: mockExecutor,
		// 			},
		// 		},
		// 		want:    nil,
		// 		wantErr: true,
		// 	}
		// },
		// "Given an invalid output format, When listing processes, Then it should return an error": func(t *testing.T, ctrl *gomock.Controller) test {
		// 	var stdout, stderr bytes.Buffer
		// 	mockExecutor := NewMockCommandExecutor(ctrl)
		// 	mockExecutor.EXPECT().ExecuteCommand(&stdout, &stderr, "sh", "-c", "ps aux").DoAndReturn(func(stdout, stderr *bytes.Buffer, command string, args ...string) error {
		// 		stdout.WriteString("invalid output format")
		// 		return nil
		// 	})

		// 	return test{
		// 		args: args{
		// 			fields: gcprocess.ProcessFilter{
		// 				Attrs:           map[string]string{},
		// 				CommandExecutor: mockExecutor,
		// 			},
		// 		},
		// 		want:    nil,
		// 		wantErr: true,
		// 	}
		// },
		// "Given a real system with no filter, When listing processes, Then it should return at least one process": func(t *testing.T, ctrl *gomock.Controller) test {
		// 	return test{
		// 		args: args{
		// 			fields: gcprocess.ProcessFilter{
		// 				Attrs:           map[string]string{},
		// 				CommandExecutor: gcprocess.RealCommandExecutor{},
		// 			},
		// 		},
		// 		want:    []map[string]string{{ /* expect at least one process */ }},
		// 		wantErr: false,
		// 	}
		// },
		// "Given a real system with a filter, When listing processes, Then it should return filtered processes": func(t *testing.T, ctrl *gomock.Controller) test {
		// 	return test{
		// 		args: args{
		// 			fields: gcprocess.ProcessFilter{
		// 				Attrs:           map[string]string{"USER": "root"},
		// 				CommandExecutor: gcprocess.RealCommandExecutor{},
		// 			},
		// 		},
		// 		want:    []map[string]string{{ /* expect processes filtered by USER=root */ }},
		// 		wantErr: false,
		// 	}
		// },
	}

	for name, testFn := range tests {
		t.Run(name.ToString(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			// Given
			sut := sut(tt.args)

			// When
			got, err := sut.ListProcesses()

			// Then
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessFilter.ListProcesses() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) && (tt.want != nil && len(tt.want) > 0) {
				t.Errorf("ProcessFilter.ListProcesses() = %v, want %v", got, tt.want)
			}
		})
	}
}
