package fileserver

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

var debug bool = false

var headerTpl string = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" />
		<!--Patterns name fixed, rooted paths, like "/favicon.ico".-->
		<link rel="icon" href="data:;base64,iVBORw0KGgo=">
		<title>File Server In Go</title>
	</head>
	<body>
		<b>
			{{ range $i, $rd := $.RootDir }}{{ $rdn := index $.RootDirName $i }}／<a href="{{ $.HandlePrefix }}{{ $rd }}">{{ $rdn }}</a>
			{{ end }}
		</b>
		<hr>
		<form action="/0838117db8c2f8882e65bfddaccaab2e" method="POST">
			<button type="submit" value="{{ $.RelateName }}" name="execute">Archive this Dir</button>
			<input type="hidden" name="op" value="archive">
		</form>
`
var footerTpl string = `	</body>
</html>`

// ReRender will re-render Body with http.FileServer
type ReRender struct {
	HandlePattern string
	FileServerDir string
}

func (fsrd *ReRender) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// At the first time visiting website, browser will check "favicon.ico" file, that is r.URL.Path = "/favicon.ico".
	// This will cause os.Stat(f) to cannot find the path specified error.
	// reference: https://stackoverflow.com/questions/1321878/how-to-prevent-favicon-ico-requests
	// if r.URL.Path == "/favicon.ico" {
	// 	return
	// }

	// Get relative path f that in the server's filesystem.
	name := RequestName(fsrd.HandlePattern, r)
	f, err := DirAddName(fsrd.FileServerDir, name)
	if err != nil {
		return
	}

	// Check if path d does exist.
	d, err := os.Stat(f)
	if err != nil {
		return
	}

	if debug {
		fmt.Println("fsrd.HandlePattern: ", fsrd.HandlePattern)
		fmt.Println("fsrd.FileServerDir: ", fsrd.FileServerDir)
		fmt.Println("name: ", name)
		fmt.Println("fullname: ", f)
		fmt.Println("fullname stat: ", d)
		fmt.Println("err: ", err)
	}

	// Render Body if target path is dir.
	if d.IsDir() {
		t := template.New("body")
		body, err := t.Parse(headerTpl)
		if err != nil {
			return
		}

		rp, rpn := RootPath(name)
		for i, path := range rp {
			rp[i] = filepath.ToSlash(path + "/")
			if rpn[i] == "." {
				rpn[i] = "【HOME】"
			}
		}
		if debug {
			fmt.Println("rp: ", rp, "\trpn: ", rpn)
		}

		body.Execute(w, struct {
			HandlePrefix string
			RelateName   string
			RootDir      []string
			RootDirName  []string
		}{
			HandlePrefix: fsrd.HandlePattern,
			RelateName:   f,
			RootDir:      rp,
			RootDirName:  rpn,
		})
		if debug {
			fmt.Println("***Is dir (render site body)***")
		}
	}

	// Render Body via http.FileServer
	fs := http.StripPrefix(fsrd.HandlePattern, http.FileServer(http.Dir(fsrd.FileServerDir)))
	fs.ServeHTTP(w, r)

	// Render Footer if target path is dir.
	if d.IsDir() {
		io.WriteString(w, footerTpl)
	}
}

// RequestName returns a string name that serves HTTP requests by removing the given prefix from the request URL's Path and invoking the request r. The name is '/'-separated.
func RequestName(prefix string, r *http.Request) (name string) {
	upath := r.URL.Path
	if debug {
		fmt.Println("----------")
		fmt.Println("r.URL.Path: ", upath)
	}

	// Reference: net/http/server.go: 2070: func StripPrefix(prefix string, h Handler) Handler {
	if prefix != "" {
		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			upath = p
			if debug {
				fmt.Println("upath: ", upath)
			}
		}
	}

	// Reference: net/http/fs.go:724: func (f *fileHandler) ServeHTTP(w ResponseWriter, r *Request) {
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		if debug {
			fmt.Println("hasprefix: ", upath)
		}
	}
	return path.Clean(upath)
}

// DirAddName implements FileSystem using os.Stat, checking files for reading rooted and relative to the directory dir.
func DirAddName(dir, name string) (relatePath string, err error) {
	// Reference: net/http/fs.go:70: func (d Dir) Open(name string) (File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return "", errors.New("invalid character in file path")
	}
	d := string(dir)
	if d == "" {
		d = "."
	}
	relatePath = filepath.Join(d, filepath.FromSlash(path.Clean("/"+name)))
	if debug {
		fmt.Println("dir: ", dir)
		fmt.Println("d: ", d)
		fmt.Println("name: ", name)
		fmt.Println("relatePath: ", relatePath)
	}
	if _, err := os.Stat(relatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("%s does not exist", relatePath)
	}
	return relatePath, nil
}

// RootPath traverse all root path for giving dir.
func RootPath(dir string) (rootPaths, rootPathName []string) {
	clean := filepath.Clean(dir)
	dp := strings.TrimLeft(clean, `/\`)
	count := strings.Count(dp, `/`) + strings.Count(dp, `\`) + 1 // Has at least one root dir itselfe.
	// Remove count of drive letter's trailing backslash
	if v := filepath.VolumeName(dp); len(v) == 2 && v[1] == ':' {
		count--
	}
	rootPaths = make([]string, count)
	rootPathName = make([]string, count)
	// Dir returns all but the last element of path, typically the path's directory.
	// Given "C:\foo\bar" it returns "C:\foo" on Windows.
	// Given "host\share\foo" it returns "host\share" on Unix.
	rootPaths[count-1] = filepath.Dir(dp)
	_, rootPathName[count-1] = filepath.Split(rootPaths[count-1])

	for i := count - 2; i >= 0; i-- {
		rootPaths[i] = filepath.Dir(rootPaths[i+1])
		_, rootPathName[i] = filepath.Split(rootPaths[i])
	}
	return rootPaths, rootPathName
}
