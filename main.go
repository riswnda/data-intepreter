package main

import (
	"data-inteprenter/entity"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
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
	originFlatten := Flatten(data.OriginData)
	targetData = checkFormat(ParsingFormat, originFlatten)
	out, err := Unflatten(targetData)
	spew.Dump(out)
	if err != nil {
		fmt.Println(err)
		return
	}

	object, err := json.Marshal(out)
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
		case []interface{}:
			for i := 0; i < len(child); i++ {
				o[k+"."+strconv.Itoa(i)] = child[i]
			}
		default:
			o[k] = v
		}
	}
	return o
}

func Unflatten(flat map[string]any) (map[string]any, error) {
	unflat := map[string]any{}

	for key, value := range flat {
		keyParts := strings.Split(key, ".")

		// Walk the keys until we get to a leaf node.
		m := unflat
		for i, k := range keyParts[:len(keyParts)-1] {
			v, exists := m[k]
			if !exists {
				newMap := map[string]any{}
				m[k] = newMap
				m = newMap
				continue
			}

			innerMap, ok := v.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("key=%v is not an object", strings.Join(keyParts[0:i+1], "."))
			}
			m = innerMap
		}

		leafKey := keyParts[len(keyParts)-1]
		if _, exists := m[leafKey]; exists {
			return nil, fmt.Errorf("key=%v already exists", key)
		}
		m[keyParts[len(keyParts)-1]] = value
	}

	return unflat, nil

}

func checkFormat(f []entity.Format, originFlatten map[string]any) map[string]any {
	par := make(map[string]any)

	for _, parV := range f {
		format, ft := getFormatType(parV.Format)
		for key, value := range originFlatten {
			//fmt.Printf("key dari originFlatten >>> %v\nparV.Origin            >>> %v\n\n", key, parV.Origin)
			if strings.Contains(key, parV.Origin) {
				switch ft {
				case Empty:
					par[parV.Target] = value
				case Date:
					t, err := dateparse.ParseAny(fmt.Sprint(value))
					if err != nil {
						log.Fatal(err)
					}
					par[parV.Target] = t.Format(format)
				case String:
					address, err := addressFormat(&format, key, fmt.Sprintf("%s", value))
					if err == nil {
						par[parV.Target] = address
					}
					log.Fatal(err)

				}
			}

		}
	}
	return par
}

func getFormatType(f string) (format string, formatType FormatType) {
	// Conditional to check empty format
	if len(f) == 0 {
		return "", Empty
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
