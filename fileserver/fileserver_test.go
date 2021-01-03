package fileserver

import (
	"reflect"
	"runtime"
	"testing"
)

func TestRootPath(t *testing.T) {
	type args struct {
		dir string
	}

	type test struct {
		name             string
		args             args
		wantRootPaths    []string
		wantRootPathName []string
	}

	wintests := []test{
		// TODO: Add test cases.
		{
			name: "Windows test 01",
			args: args{
				dir: "\\foo\\bar\\baz\\",
			},
			wantRootPaths:    []string{".", "foo", "foo\\bar"},
			wantRootPathName: []string{".", "foo", "bar"},
		},
		{
			name: "Windows test 02",
			args: args{
				dir: "\\foo\\bar\\baz",
			},
			wantRootPaths:    []string{".", "foo", "foo\\bar"},
			wantRootPathName: []string{".", "foo", "bar"},
		},
		{
			name: "Windows test 03",
			args: args{
				dir: "foo\\bar\\baz\\",
			},
			wantRootPaths:    []string{".", "foo", "foo\\bar"},
			wantRootPathName: []string{".", "foo", "bar"},
		},
		{
			name: "Windows test 04",
			args: args{
				dir: "foo\\bar\\baz",
			},
			wantRootPaths:    []string{".", "foo", "foo\\bar"},
			wantRootPathName: []string{".", "foo", "bar"},
		},
		{
			name: "Windows test 05",
			args: args{
				dir: "C:\\foo\\bar\\baz\\",
			},
			wantRootPaths:    []string{"C:\\", "C:\\foo", "C:\\foo\\bar"},
			wantRootPathName: []string{"", "foo", "bar"},
		},
		{
			name: "Windows test 06",
			args: args{
				dir: "C:\\foo\\bar\\baz",
			},
			wantRootPaths:    []string{"C:\\", "C:\\foo", "C:\\foo\\bar"},
			wantRootPathName: []string{"", "foo", "bar"},
		},
	}

	unixtests := []test{
		// TODO: Add test cases.
		{
			name: "Linux test 01",
			args: args{
				dir: "/foo/bar/baz/",
			},
			wantRootPaths:    []string{".", "foo", "foo/bar"},
			wantRootPathName: []string{".", "foo", "bar"},
		},
		{
			name: "Linux test 02",
			args: args{
				dir: "/foo/bar/baz",
			},
			wantRootPaths:    []string{".", "foo", "foo/bar"},
			wantRootPathName: []string{".", "foo", "bar"},
		},
		{
			name: "Linux test 03",
			args: args{
				dir: "foo/bar/baz/",
			},
			wantRootPaths:    []string{".", "foo", "foo/bar"},
			wantRootPathName: []string{".", "foo", "bar"},
		},
		{
			name: "Linux test 04",
			args: args{
				dir: "foo/bar/baz",
			},
			wantRootPaths:    []string{".", "foo", "foo/bar"},
			wantRootPathName: []string{".", "foo", "bar"},
		},
	}

	var tests []test
	if runtime.GOOS == "windows" {
		tests = append(tests, wintests...)
	} else {
		tests = append(tests, unixtests...)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRootPaths, gotRootPathName := RootPath(tt.args.dir)
			if !reflect.DeepEqual(gotRootPaths, tt.wantRootPaths) {
				t.Errorf("RootPath() gotRootPaths = %v, want %v", gotRootPaths, tt.wantRootPaths)
			}
			if !reflect.DeepEqual(gotRootPathName, tt.wantRootPathName) {
				t.Errorf("RootPath() gotRootPathName = %v, want %v", gotRootPathName, tt.wantRootPathName)
			}
		})
	}
}
