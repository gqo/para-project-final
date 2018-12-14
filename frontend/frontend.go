// Graeme Ferguson | ggf221 | 10/24/18
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var backendAddr string

/* Http request handler for "/" page. Displays items and links for
creation, editing, and deletion. */
func handler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("main.html"))
	// Read data from backend
	rData := bReadAll()
	t.Execute(w, rData)
}

/* Http request handler for "/save/" pages. Receives form value from either
/edit/ or /create/ and saves item info accordingly (creating a new one
if necessary) */
func saveHandler(w http.ResponseWriter, r *http.Request) {
	/* dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump)) */
	ID, _ := strconv.Atoi(r.URL.Path[len("/save/"):])
	ID32 := int32(ID)
	Album := r.PostFormValue("albumName")
	Artist := r.PostFormValue("artistName")
	Rating, _ := strconv.Atoi(r.PostFormValue("rating"))
	Rating32 := int32(Rating)
	Body := r.PostFormValue("bodyText")
	// Could do this logic on backend
	if ID == 0 {
		bCreate(Album, Artist, Rating32, Body)
	} else {
		if _, exist := bRead(ID32); exist {
			bUpdate(ID32, Album, Artist, Rating32, Body)
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

/* Http request handler for "/create/" pages. Allows users to enter form
values and then redirects form submissions to "/save/" */
func createHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./create.html")
}

/* Http request handler for "/edit/" pages. Similar to "/create/" with
the caveat that the forms have values filled in from the item to be
edited. Upon form submission, redirects to "/save/" */
func editHandler(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.URL.Path[len("/edit/"):])
	ID32 := int32(ID)
	review, exist := bRead(ID32)
	if !exist {
		http.Redirect(w, r, "/create/", http.StatusFound)
		return
	}
	t := template.Must(template.ParseFiles("edit.html"))
	t.Execute(w, review)
}

/* Http request handler for "/delete/" pages. Receives id of item and
sends delete request to backend with id data */
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.URL.Path[len("/delete/"):])
	ID32 := int32(ID)
	// Could do this logic on backend
	_, exist := bRead(ID32)
	if exist {
		bDelete(ID32)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	httpPort := flag.Int("listen", 8080, "Specify port"+
		" for web server to listen on.")
	backendPortsString := flag.String("backend", "localhost:8090", `Specify hostname:port of backends in comma
		`+`separated list.`)
	flag.Parse()
	backendPorts := strings.Split(*backendPortsString, ",")

	backendAddr = backendPorts[0]
	portStr := ":" + strconv.Itoa(*httpPort)

	http.HandleFunc("/", handler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/delete/", deleteHandler)

	fmt.Println("Listening to HTTP requests on port:", portStr[1:])
	fmt.Println("Connect made to backend at:", backendAddr)
	go bPing()

	bStart(backendPorts)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
