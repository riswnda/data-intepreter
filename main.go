package main

import (
	"data-inteprenter/entity"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	data       entity.InputData
	targetData map[string]any
)

type FormatType string

const (
	Empty  = FormatType("EMPTY")
	Date   = FormatType("DATE")
	String = FormatType("STRING")
)

func getData(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		if err == io.EOF {
			http.Error(w, "Request body is empty", http.StatusUnprocessableEntity)
		} else {
			http.Error(w, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	defer r.Body.Close()

	ParsingFormat := data.ParsingFormat
	targetData = checkFormat(ParsingFormat)

	object, err := json.Marshal(targetData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error Marshal Response %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(object)

}

func Flatten(m map[string]any) map[string]any {
	o := make(map[string]any)
	for k, v := range m {
		switch child := v.(type) {
		case map[string]any:
			nm := Flatten(child)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		default:
			o[k] = v
		}
	}
	return o
}

func checkFormat(f []entity.Format) map[string]any {
	par := make(map[string]any)
	originFlatten := Flatten(data.OriginData)

	for _, parV := range f {
		format, ft := getFormatType(parV.Format)
		for key, value := range originFlatten {
			if strings.Contains(key, parV.Origin) {
				switch ft {
				case Empty:
					par[parV.Target] = value
				case String:
					address, err := addressFormat(&format, key, fmt.Sprintf("%s", value))
					if err == nil {
						par[parV.Target] = address
					}
				case Date:
					t, err := dateparse.ParseAny(fmt.Sprint(value))
					if err != nil {
						log.Fatal(err)
					}
					par[parV.Target] = t.Format(format)

				}
			}
		}
	}
	return par
}

func getFormatType(f string) (format string, formatType FormatType) {
	// Conditional to check empty format
	if len(f) <= 0 {
		return f, Empty
	}
	// Conditional to check date format
	dateFormat, err := time.Parse(f, f)
	if err == nil {
		if !dateFormat.IsZero() && dateFormat.Year() != 0 {
			return f, Date
		}
	}
	// If none of format match with 2 conditional above, return format string
	return f, String
}

func addressFormat(pattern *string, format ...string) (string, error) {
	if strings.Contains(*pattern, format[0]) {
		*pattern = strings.ReplaceAll(strings.Replace(*pattern, format[0], format[1], -1), "$", "")
	}

	return *pattern, nil
}

func main() {

	// Initialize the router
	route := mux.NewRouter()

	//route Handlers
	route.HandleFunc("/api/data", getData).Methods("POST")

	fmt.Println("Successfully connect to port :6000")
	log.Fatal(http.ListenAndServe(":6000", route))

}
