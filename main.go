package main

import (
	"flag"
	"fmt"
	"net/http"
	"log"
	"strconv"
)

const (
	ipString = "specifies the IP address the server will bind to"
	portString = "specifies the port the server will bind to"
	quietString = "specifies whether quiet mode is enabled"
)

var ip string
var port int
var quiet bool

func main() {
	flag.StringVar(&ip, "ip", "127.0.0.1", ipString)
	flag.IntVar(&port, "port", 5000, portString)
	flag.BoolVar(&quiet, "quiet", false, quietString)
	flag.Parse()
	binding := ip + ":" + strconv.Itoa(port)
	http.HandleFunc("/", indexHandler)
	if !quiet {
		fmt.Printf("SERVER: Starting on \033[0;31m")
		fmt.Printf("%v\033[0m:\033[0;34m%v\033[0m\n", ip, port)
	}
	log.Fatal(http.ListenAndServe(binding, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if !quiet {
		fmt.Println(r.Method, r.URL)
	}
	fmt.Fprintf(w, "Hello, World!")
}
