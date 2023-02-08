package main

import (
	"data-inteprenter/entity"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	data          entity.InputData
	originFlatten map[string]interface{}
	targetData    map[string]interface{}
)

func getData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var ()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		if err == io.EOF {
			http.Error(w, "Request body is empty", http.StatusUnprocessableEntity)
		} else {
			http.Error(w, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	originFlatten = Flatten(data.OriginData)
	ParsingFormat := data.ParsingFormat
	targetData = compareData(ParsingFormat)

	object, err := json.Marshal(targetData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error Marshal Response")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(object)

}

func Flatten(m map[string]interface{}) map[string]interface{} {
	o := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
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

func compareData(f []entity.Format) map[string]interface{} {
	par := make(map[string]interface{})

	for _, parV := range f {
		for key, value := range originFlatten {
			if key == parV.Origin {
				par[parV.Target] = value
			}
			//else {
			//	fmt.Println("Target doesn't match", key, value, parV.Origin)
			//}
			if address, err := matchFormat(&parV.Format, key) {
				fmt.Println("== Address:", address)
				par[parV.Target] = address
			} else {
				fmt.Println("alamat kosong")
			}
		}

	}

	return par
}

func matchFormat(pattern *string, format ...string) (string, error) {
	if *pattern == "" {
		return "", errors.New("No Format Detected")
	}

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
