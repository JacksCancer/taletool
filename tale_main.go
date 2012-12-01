package main

import (
	"fmt"
	"log"
	"net/http"
	//"html"
	"io"
	//"os"
)

func main() {

	http.HandleFunc("/update-rules", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, r.Body)
		fmt.Println(r.Body)
		//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	//fmt.Println(os.Getwd())

	http.Handle("/edit/", http.StripPrefix("/edit/", http.FileServer(http.Dir("./server/editor/"))))
	//http.Handle("/doc/", http.FileServer(http.Dir("/usr/share/doc")))

	//log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/usr/share/doc"))))

	http.Handle("/animview/", http.StripPrefix("/animview", http.FileServer(http.Dir("server/animview/web/"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
