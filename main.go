package main

import (
	"ascii-art/functions"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func main() {
	// handle the homepae request
	http.HandleFunc("/", homepage)
	// ensure the css will be executed upon request
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// handle the result page request
	http.HandleFunc("/result", resultpage)
	// listens for incoming requests on the port mentioned below then handles those requests
	http.ListenAndServe(":8080", nil)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	// If it's nnot the homepage error handle
	if r.URL.Path != "/" {
		renderErrorPage(w, 404)
		return
	}

	// Parse the HTML file
	t, err := template.ParseFiles("index.html")
	if err != nil {
		// http.Error(w, "Error parsing html", http.StatusInternalServerError)
		renderErrorPage(w, 500)
		return
	}

	// execute the HTML template
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func resultpage(w http.ResponseWriter, r *http.Request) {
	// For resultpage the request is always POST not GET
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	// if url is not for result page error handle
	if r.URL.Path != "/result" {
		renderErrorPage(w, 404)
		return
	}

	// Get the form values
	inputString := r.FormValue("inputString")
	style := r.FormValue("style")

	// Validate the input string
	for _, ch := range inputString {
		if ch != 10 && ch != 13 && (ch < 32 || ch > 126) {
			renderErrorPage(w, 400)
			return
		}
	}

	// Process the ASCII art
	fileLines := functions.Read(style)
	asciiRep := functions.AsciiRep(fileLines)

	var res strings.Builder

	// Split the input string into lines
	inputLines := strings.Split(inputString, "\r\n")

	for _, line := range inputLines {
		if strings.TrimSpace(line) == "" {
			res.WriteString("\n")
			continue
		} else {
			asciiArt := functions.PrintStr(line, asciiRep)
			for _, asciiLine := range asciiArt {
				res.WriteString(strings.Join(asciiLine, ""))
				res.WriteString("\n")
			}
		}
	}

	// Parse the HTML template again to render the result
	t, err := template.ParseFiles("index.html")
	if err != nil {
		// http.Error(w, "Error parsing html", http.StatusInternalServerError)
		renderErrorPage(w, 500)
		return
	}


	// Render the template with the result
	err = t.Execute(w, res.String())
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

// Custom error page handler
func renderErrorPage(w http.ResponseWriter, code int) {
	// Set the status code
	w.WriteHeader(http.StatusNotFound)

	// Generate ASCII art for the error code with the "Standard" style
	inputString := strconv.Itoa(code)
	style := "standard"

	fileLines := functions.Read(style)
	asciiRep := functions.AsciiRep(fileLines)

	var res strings.Builder

	// Process the ASCII art for error code
	asciiArt := functions.PrintStr(inputString, asciiRep)
	for _, asciiLine := range asciiArt {
		res.WriteString(strings.Join(asciiLine, ""))
		res.WriteString("\n")
	}

	// Parse and render the custom 404 template
	t, err := template.ParseFiles(fmt.Sprintf("static/%d.html", code))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing %d HTML", code), http.StatusInternalServerError)
		return
	}

	// Render the template with the result
	err = t.Execute(w, res.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing %d template", code), http.StatusInternalServerError)
	}
}
