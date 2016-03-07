/*
	Author:
		Nicholas Siow | nick@siow.me
		Alani Douglas | fresh@alani.style
	Description:
		Core webserver for http://alanick.us
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

//------------------------------------------------------------
//	CONFIGURATION
//------------------------------------------------------------

var config_path = "/etc/alanick.conf"

type Config struct {
	SiteRoot string
	LogFile  string
	Debug    bool
}

var config Config

/*
	functions to run on initialization
*/
func init() {
	configSetup()
	loggingSetup()
}

/*
	function to read in and validate configuration
*/
func configSetup() {
	configFile, err := os.Open(config_path)
	if err != nil {
		fmt.Println("Error opening config file", err.Error())
		os.Exit(1)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		fmt.Println("Error parsing config file", err.Error())
		os.Exit(1)
	}
}

//------------------------------------------------------------
//	LOGGING
//------------------------------------------------------------

/*
	set up server logging
*/
func loggingSetup() {
	f, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file", err.Error())
		os.Exit(1)
	}

	log.SetOutput(f)
}

/*
	helper function for debug logging
*/
func debug(message string) {
	if config.Debug {
		log.Printf("[DEBUG] %s", message)
	}
}

/*
	helper function to standardize error logging
*/
func logError(request string, message string, statusCode int) {
	log.Printf("[ERROR] %d for %s :: %s", statusCode, request, message)
}

//------------------------------------------------------------
//	CLEANUP
//------------------------------------------------------------

/*
	clean up server loose ends
*/
func cleanup() {
}

/*
	start listening and serving requests
*/
func main() {

	// serve static css files at /static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// set up the main page handler
	http.HandleFunc("/", serve)
	http.ListenAndServe(":8080", nil)

}

/*
	base handler for serving html pages and directories
*/
func serve(w http.ResponseWriter, r *http.Request) {

	debug(fmt.Sprintf("Received request: %+v", r))

	// pull out the requested path
	requested_path := r.URL.Path
	if strings.Contains(requested_path, ".") {
		http.Error(w, "Don't be a dick...", 500)
		logError(requested_path, "Directory traversal", 500)
		return
	}

	// put the requested path in the context of the local file system
	path_to_html := path.Join(config.SiteRoot, requested_path)

	// try both options for file vs directory
	final_path := path.Join(path_to_html, "this.html")
	if _, err := os.Stat(final_path); os.IsNotExist(err) {
		final_path = path_to_html + ".html"
	}

	// verify that final_path is valid
	if _, err := os.Stat(final_path); os.IsNotExist(err) {
		http.Error(w, "This page does not exist, sorry!", 404)
		logError(requested_path, "Page not found", 404)
		return
	}

	// serve up page
	file, err := ioutil.ReadFile(final_path)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		logError(requested_path, err.Error(), 500)
		return
	}

	w.Write(file)
}
