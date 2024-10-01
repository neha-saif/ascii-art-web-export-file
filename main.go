package main

import (
	"ascii-art/functions"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type ResultData struct {
	AsciiArt    string
	InputString string
	Style       string
}

func main() {
	// handle the homepae request
	http.HandleFunc("/", homepage)
	// ensure the css will be executed upon request
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// handle the result page request
	http.HandleFunc("/result", resultpage)
	http.HandleFunc("/download", downloadfile)
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

	AsciiArt := renderasciires(inputString, style)

	data := ResultData{
		AsciiArt:    AsciiArt,
		InputString: inputString,
		Style:       style,
	}
	// Parse the HTML template again for the resultpage
	t, err := template.ParseFiles("index.html")
	if err != nil {
		// http.Error(w, "Error parsing html", http.StatusInternalServerError)
		renderErrorPage(w, 500)
		return
	}

	// Render the template with the result
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func renderasciires(inputString, style string) string {
	// Process the ASCII art
	fileLines := functions.Read(style)
	asciiRep := functions.AsciiRep(fileLines)
	var res strings.Builder
	var content [][]string

	inputString = strings.ReplaceAll(inputString, "\\n", "\n")
	inputString = strings.ReplaceAll(inputString, "\r", "")
	inputLines := strings.Split(inputString, "\n")

	for _, line := range inputLines {
		if strings.TrimSpace(line) == "" {
			res.WriteString("\n")
			continue
		}
		content = functions.PrintStr(line, asciiRep)
		for i, asciiLine := range content {
			res.WriteString(strings.Join(asciiLine, ""))
			if i < len(content) {

				res.WriteString("\n")

			}
		}
		//res.WriteString("\n")

	}
	AsciiArt := res.String()
	return AsciiArt
}

func downloadfile(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/download" {
		renderErrorPage(w, 404)
		return
	}

	// Get the form values for inputstring and style hidden in the html page
	inputString := r.FormValue("inputString")
	style := r.FormValue("style")

	// Generate ASCII art with output in string form
	AsciiArt := renderasciires(inputString, style)

	fileName := "outputFile.txt"

	// Create the file and write the ASCII art to it
	file, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}

	// Write the ASCII art to the file created abnove
	_, err = file.WriteString(AsciiArt)
	if err != nil {
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	file.Close()

	// open the file for reading
	file, err = os.Open(fileName)
	if err != nil {
		http.Error(w, "Unable to open file: "+fileName, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get file information for content length
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not get the file info", http.StatusInternalServerError)
		return
	}

	// Set header parameters for the file to nbe downloaded
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Copy the file content to the outputfile in bbrowser
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error copying file to response", http.StatusInternalServerError)
		return
	}

}

func renderErrorPage(w http.ResponseWriter, code int) {
	w.WriteHeader(http.StatusNotFound)

	inputString := strconv.Itoa(code)
	style := "standard"

	fileLines := functions.Read(style)
	asciiRep := functions.AsciiRep(fileLines)

	var res strings.Builder

	asciiArt := functions.PrintStr(inputString, asciiRep)
	for _, asciiLine := range asciiArt {
		res.WriteString(strings.Join(asciiLine, ""))
		res.WriteString("\n")
	}

	t, err := template.ParseFiles(fmt.Sprintf("static/%d.html", code))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing %d HTML", code), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, res.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing %d template", code), http.StatusInternalServerError)
	}
}
