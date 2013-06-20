package main

import (
	// "fmt"
	"log"
	"net/http"
	// "io"
	"flag"
	"os"
	"path"
)

func createSaveHandler(prefix string, root string) http.HandlerFunc {
	lenPrefix := len(prefix)
	return func(w http.ResponseWriter, r *http.Request) {
		path := path.Clean(r.URL.Path[lenPrefix-1:])

		file, err := os.Create(root + path)

		if err == nil {
			_, err = file.WriteString(r.FormValue("data"))
		}

		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func main() {

	var root, workdir, port string

	flag.StringVar(&root, "root", "web", "http root")
	flag.StringVar(&workdir, "workdir", "project", "directory to serve")
	flag.StringVar(&port, "port", ":8080", "port for http server")
	flag.Parse()

	savePrefix := "/save/"
	http.HandleFunc(savePrefix, createSaveHandler(savePrefix, workdir))
	loadPrefix := "/load/"
	http.Handle(loadPrefix, http.StripPrefix(loadPrefix, http.FileServer(http.Dir(workdir))))
	http.Handle("/", http.FileServer(http.Dir(root)))
	log.Fatal(http.ListenAndServe(port, nil))
}
