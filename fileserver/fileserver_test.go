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

func TestDirAddName(t *testing.T) {
	type args struct {
		dir  string
		name string
	}

	type test struct {
		name           string
		args           args
		wantRelatePath string
		wantErr        bool
	}

	wintests := []test{
		// TODO: Add test cases.
		{
			name: "Windows test 01",
			args: args{
				dir:  ".\\testdir\\",
				name: "test.file",
			},
			wantRelatePath: "testdir\\test.file",
			wantErr:        false,
		},
		{
			name: "Windows test 02",
			args: args{
				dir:  "",
				name: "testdir",
			},
			wantRelatePath: "testdir",
			wantErr:        false,
		},
	}

	unixtests := []test{
		// TODO: Add test cases.
		{
			name: "Linux test 01",
			args: args{
				dir:  "./testdir/",
				name: "test.file",
			},
			wantRelatePath: "testdir/test.file",
			wantErr:        false,
		},
		{
			name: "Linux test 02",
			args: args{
				dir:  "",
				name: "testdir",
			},
			wantRelatePath: "testdir",
			wantErr:        false,
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
			gotRelatePath, err := DirAddName(tt.args.dir, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("DirAddName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRelatePath != tt.wantRelatePath {
				t.Errorf("DirAddName() = %v, want %v", gotRelatePath, tt.wantRelatePath)
			}
		})
	}
}
