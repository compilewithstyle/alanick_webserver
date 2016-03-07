/*
	Author:
		Nicholas Siow | nick@siow.me
		Alani Douglas | fresh@alani.style
	Description:
		Core webserver for http://alanick.us
*/

package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

var site_root string = "/opt/go/src/github.com/compilewithstyle/alanick_webserver/site_root"

/*
	start listening and serving requests
*/
func main() {

	// serve static css files at /static
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// set up the main page handler
	http.HandleFunc("/", serve)
	http.ListenAndServe(":8080", nil)

}

/*
	base handler for serving html pages and directories
*/
func serve(w http.ResponseWriter, r *http.Request) {

	// pull out the requested path
	requested_path := r.URL.Path
	if strings.Contains(requested_path, ".") {
		http.Error(w, "Don't be a dick...", 500)
		return
	}

	// put the requested path in the context of the local file system
	path_to_html := path.Join(site_root, requested_path)

	// try both options for file vs directory
	final_path := path.Join(path_to_html, "this.html")
	if _, err := os.Stat(final_path); os.IsNotExist(err) {
		final_path = path_to_html + ".html"
	}

	// verify that final_path is valid
	if _, err := os.Stat(final_path); os.IsNotExist(err) {
		http.Error(w, "This page does not exist, sorry!", 404)
		return
	}

	// serve up page
	file, err := ioutil.ReadFile(final_path)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}

	w.Write(file)
}
