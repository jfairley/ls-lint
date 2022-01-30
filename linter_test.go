package main

import (
	"errors"
	"io/fs"
	"reflect"
	"sync"
	"testing"
	"testing/fstest"
)

func TestLinterRun(t *testing.T) {
	var tests = []*struct {
		filesystem        fs.FS
		config            *Config
		linter            *Linter
		expectedErr       error
		expectedStatistic *Statistic
		expectedErrors    []*Error
	}{
		{
			filesystem: fstest.MapFS{
				"snake_case.png": new(fstest.MapFile),
			},
			config: &Config{
				Ls: map[string]interface{}{
					".png": "snake_case",
				},
				Ignore: []string{
					"node_modules",
				},
				RWMutex: new(sync.RWMutex),
			},
			linter: &Linter{
				Statistic: &Statistic{
					Files:     0,
					FileSkips: 0,
					Dirs:      0,
					DirSkips:  0,
					RWMutex:   new(sync.RWMutex),
				},
				Errors:  []*Error{},
				RWMutex: new(sync.RWMutex),
			},
			expectedErr: nil,
			expectedStatistic: &Statistic{
				Files:     1,
				FileSkips: 0,
				Dirs:      1,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*Error{},
		},
	}

	var i = 0
	for _, test := range tests {
		err := test.linter.Run(test.filesystem, true, test.config)

		if !errors.Is(err, test.expectedErr) {
			t.Errorf("Test %d failed with unmatched error value - %v", i, err)
		}

		if !reflect.DeepEqual(test.linter.getStatistic(), test.expectedStatistic) {
			t.Errorf("Test %d failed with unmatched linter statistic values\nexpected: %+v\nactual: %+v", i, test.expectedStatistic, test.linter.getStatistic())
		}

		if !reflect.DeepEqual(test.linter.getErrors(), test.expectedErrors) {
			t.Errorf("Test %d failed with unmatched linter errors value\nexpected: %+v\nactual: %+v", i, test.expectedErrors, test.linter.getErrors())
		}

		i++
	}
}
