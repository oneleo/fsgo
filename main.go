package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/oneleo/fsgo/fileserver"
)

var (
	port   *int64  = new(int64)
	folder *string = new(string)
	browse *bool   = new(bool)
)

const (
	handlePrefix string = "/"
)

func init() {
	const (
		defaultPort = 80
		usagePort   = "Set listening port"

		defaultFolder = "./"
		usageFolder   = "Set shared folder"

		defaultBrowse = false
		usageBrowse   = "Open URL in browser"
	)
	// flag.Int64Var(&port, "port", defaultPort, usagePort)
	flag.Int64Var(port, "p", defaultPort, usagePort+" (shorthand)")
	// flag.StringVar(&folder, "folder", defaultFolder, usageFolder)
	flag.StringVar(folder, "f", defaultFolder, usageFolder+" (shorthand)")
	// flag.BoolVar(browse, "browse", defaultBrowse, usageBrowse)
	flag.BoolVar(browse, "b", defaultBrowse, usageBrowse+" (shorthand)")

	flag.Usage = func() {
		// template: go help doc
		fmt.Fprintf(os.Stderr, "Gofs is a small file server for Browsing specify local folder and files.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\n\tgofs [-p port] [-f folder] [-b]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  gofs\n\tStart this file server listen on default port \":80\" for browse \"./\" current folder.\n")
		fmt.Fprintf(os.Stderr, "  gofs -b\n\tStart this file server listen on default port \":80\" for browse \"./\" current folder, and \"auto-open\" in browser.\n")
		fmt.Fprintf(os.Stderr, "  gofs -p 8081 -f ./www -b\n\tStart this file server listen on specify port \":8081\" for browse \"./www\" specify folder, and \"auto-open\" in browser.")
	}
}

// openURI will open URI in browser.
// Reconmannd to use this new project: https://github.com/toqueteos/webbrowser
func openURI(uri string) {
	var c *string = new(string)
	var args *[]string = new([]string)

	// Obtain OS type.
	switch runtime.GOOS {
	case "windows":
		// *c = "cmd"
		// *args = []string{"/c", "start"}
		*c = "rundll32"
		*args = []string{"url.dll,FileProtocolHandler"}
	case "darwin":
		*c = "open"
	case "linux":
		*c = "xdg-open"
	default:
		log.Fatal("Unsupported operating system\n")
	}
	*args = append(*args, uri)
	// Open URI in browser.
	if err := exec.Command(*c, *args...).Start(); err != nil {
		log.Fatal(err, "\n")
	}
}

// GetUsedIP preferred outbound ip of this machine.
func GetUsedIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	if err != nil {
		return "127.0.0.1", err
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// StartGraceServer will start server gracefully.
func StartGraceServer(srv *http.Server) {
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// SIGTERM:	Termination signal defaultly
	// SIGINT:	Interrupt from keyboard		= $ kill -2 PID	= CTRL + C
	// SIGKILL:	Kill signal					= $ kill -9 PID
	// SIGQUIT:	Quit from keyboard			= CTRL + \
	// SIGTSTP:	Stop typed at terminal		= CTRL + Z
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
	<-stop

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("Server shutteddown gracefully")
}

func main() {
	flag.Parse()
	// Obtain absolute folder even in different OS.
	fullname, err := filepath.Abs(*folder)
	if err != nil {
		log.Fatal(err, "\n")
	}
	// Configuration setting for file server.
	// fs := http.StripPrefix(handlePrefix, http.FileServer(http.Dir(fullname)))
	fs := &fileserver.ReRender{
		HandlePattern: handlePrefix,
		FileServerDir: fullname}
	pt := fmt.Sprintf(":%d", *port)
	ip, _ := GetUsedIP()
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr: pt,
	}
	(*srv).Handler = mux

	// Handle function(s).
	//(*mux).Handle(handlePrefix, fs)

	(*mux).Handle(handlePrefix, fs)

	mux.HandleFunc("/0838117db8c2f8882e65bfddaccaab2e", (&fileserver.Module{Path: "."}).Compress)

	// If '-b' flag is set, open URI in browser.
	if *browse == true {
		go openURI("http://" + ip + pt)
	}

	// Start file server gracefully.
	fmt.Print("File Server is Started to listen on http://", ip, pt, " for browse ", fullname, "\nPlease press CTRL + C to finish the file server...\n")
	StartGraceServer(srv)
}
