package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ResultAscii struct {
	TextAscii string
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 NOT FOUND")
	}
	if status == http.StatusBadRequest {
		fmt.Fprint(w, "ERROR 400 BAD REQUEST")
	}
	if status == http.StatusInternalServerError {
		fmt.Fprint(w, "ERORR 500 INTERNAL SERVER")
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	namePolice := r.FormValue("police")
	if namePolice == "" {
		namePolice = "standard"
	}
	file, err := os.Open("assets/" + namePolice + ".txt")
	fmt.Println(r.FormValue("police"))
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	content, _ := ioutil.ReadAll(file)

	table := strings.Split(string(content), "\n")

	var result []string

	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	} else {
		textChecker := r.FormValue("asciitext")
		for i := range textChecker {
			if textChecker[i] == '\n' || textChecker[i] == '\r' {
				continue
			}
			if textChecker[i] < 32 || textChecker[i] > 127 {
				fmt.Fprintf(w, "non")
				// fmt.Println(int(textChecker[i]))
				return
			}
		}
	}
	line := strings.Split(r.FormValue("asciitext"), "\\n")
	strTemp := strings.Join(line, "\n")
	strTemp = strings.ReplaceAll(strTemp, "\\n", string([]byte{0x0D, 0x0A}))
	jump := strings.ReplaceAll(strTemp, string([]byte{0x0D, 0x0A}), "\n")
	contentAll := strings.Split(jump, "\n")
	for i := 0; i < len(contentAll); i++ {
		if len(contentAll[i]) > 0 {
			chars := []rune(contentAll[i])
			for n := 0; n < 8; n++ {
				for v := 0; v < len(chars); v++ {
					group := int(chars[v]) - 32
					adress := group * 9
					charLine := table[adress+1+n]
					result = append(result, charLine)
				}
				result = append(result, string(rune('\n')))
			}
		} else {
			result = append(result, string(rune('\n')))
		}
	}
	sresult := ""
	for i := range result {
		sresult += result[i]
	}
	t, err := template.ParseFiles("./templates/home.html")

	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	resFinal := ResultAscii{sresult}
	t.Execute(w, resFinal)
}
