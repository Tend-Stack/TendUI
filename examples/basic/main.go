// Run with: go run ./examples/basic
package main

import (
	"flag"
	"fmt"
	"github.com/Tend-Stack/GothUI/v2026"
	"github.com/Tend-Stack/GothUI/v2026/examples/site"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

var version = gothui.Version + "-dev"

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "local listen address")
	open := flag.Bool("open", false, "open the component playground in the default browser")
	showVersion := flag.Bool("version", false, "print the embedded build version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Println(version)
		return
	}
	server := &http.Server{
		Addr: *addr, Handler: site.Handler(),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second,
	}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}
	address := "http://" + listener.Addr().String()
	log.Printf("GothUI %s at %s/playground", version, address)
	if *open {
		go func() {
			program, args := browserOpener(runtime.GOOS, address+"/playground")
			if err := exec.Command(program, args...).Run(); err != nil {
				log.Printf("Open %s/playground in your browser (%v)", address, err)
			}
		}()
	}
	log.Fatal(server.Serve(listener))
}

// Pass the generated local URL as an argument, never through a command shell.
func browserOpener(goos, url string) (string, []string) {
	switch goos {
	case "windows":
		return "rundll32.exe", []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		return "open", []string{url}
	default:
		return "xdg-open", []string{url}
	}
}
