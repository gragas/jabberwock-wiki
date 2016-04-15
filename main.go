package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gragas/fsnotify"
	"github.com/gragas/jabberwock-lib/ingredient"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	ipString         = "specifies the IP address the server will bind to"
	portString       = "specifies the port the server will bind to"
	quietString      = "specifies whether quiet mode is enabled"
	assetsPathString = "specifies where assets are located"
)

var ip string
var port int
var quiet bool
var assetsPath string

var templateNames = [...]string{"templates/index.html", "templates/list.html", "templates/fungi.html"}
var templates = map[string]*template.Template{}
var pages = map[string][]byte{}

func main() {
	flag.StringVar(&ip, "ip", "127.0.0.1", ipString)
	flag.IntVar(&port, "port", 5000, portString)
	flag.BoolVar(&quiet, "quiet", false, quietString)
	flag.StringVar(&assetsPath, "assetsPath", "./", assetsPathString)
	flag.Parse()
	binding := ip + ":" + strconv.Itoa(port)

	parseTemplates()
	generateAllPages()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/list/", listHandler)
	http.HandleFunc("/fungi/", fungiHandler)

	cssHandler := http.FileServer(http.Dir("./css/"))
	imageHandler := http.FileServer(http.Dir("./images/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/images/", http.StripPrefix("/images/", imageHandler))

	if !quiet {
		fmt.Printf("SERVER: Starting on \033[0;31m")
		fmt.Printf("%v\033[0m:\033[0;34m%v\033[0m\n", ip, port)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case _ = <-watcher.Event:
				generateAllPages()
			case err := <-watcher.Error:
				fmt.Printf("SERVER: WATCHER: error: %v\n", err)
			}
		}
	}()
	err = watcher.Watch("./")
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	log.Fatal(http.ListenAndServe(binding, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if !quiet {
		fmt.Println(r.Method, r.URL)
	}
	if r.Method == "GET" && r.URL.String() == "/" {
		err := templates["templates/index.html"].Execute(w, nil)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Fprintf(w, "Bad request.")
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if !quiet {
		fmt.Println(r.Method, r.URL)
	}
	if r.Method == "GET" {
		err := templates["templates/list.html"].Execute(w, nil)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Fprintf(w, "Bad request.")
	}
}

func fungiHandler(w http.ResponseWriter, r *http.Request) {
	if !quiet {
		fmt.Println(r.Method, r.URL)
	}
	badRequest := func() {
		fmt.Fprintf(w, "Bad request.\nMETHOD: %v\nURL: %v", r.Method, r.URL)
	}
	if r.Method == "GET" {
		urlString := r.URL.String()
		urlLen := len(urlString)
		startIndex := len("/fungi/")
		if urlLen <= startIndex {
			badRequest()
			return
		}
		lastIndex := strings.LastIndex(urlString, "/")
		if lastIndex == -1 || lastIndex <= startIndex {
			lastIndex = urlLen
		}
		pagename := "./assets/fungi/" + urlString[startIndex:lastIndex]
		if pages[pagename] != nil {
			_, err := w.Write(pages[pagename])
			if err != nil {
				panic(err)
			}
			return
		}
	}
	badRequest()
}

func parseTemplates() {
	for _, templateName := range templateNames {
		tmpl, err := template.ParseFiles(templateName)
		if err != nil {
			panic(err)
		}
		if !quiet {
			fmt.Printf("Added '%v' to templates.\n", templateName)
		}
		templates[templateName] = tmpl
	}
}

func generateAllPages() {
	generateAllFungiPages()
	generateAllPlantaePages()
}

func generateAllFungiPages() {
	basePath := assetsPath + "assets/fungi"
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		dirPath := basePath + "/" + file.Name()
		dirInfo, err := os.Stat(dirPath)
		if err != nil {
			panic(err)
		}
		if dirInfo.IsDir() {
			generateFungiPage(dirPath)
		}
	}
}

func generateFungiPage(dirname string) {
	fungus := ingredient.FromFile(dirname + "/about")
	prettyFungus := fungus.PrettyIngredient()
	url := strings.Replace(dirname, " ", "_", -1)
	var buf bytes.Buffer
	err := templates["templates/fungi.html"].Execute(&buf, prettyFungus)
	if err != nil {
		panic(err)
	}
	if !quiet {
		fmt.Printf("Added '%v' to pages.\n", url)
	}
	pages[url] = buf.Bytes()
}

func generateAllPlantaePages() {
	basePath := assetsPath + "assets/plantae"
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		dirPath := basePath + "/" + file.Name()
		dirInfo, err := os.Stat(dirPath)
		if err != nil {
			panic(err)
		}
		if dirInfo.IsDir() {
			generateFungiPage(dirPath)
		}
	}
}
