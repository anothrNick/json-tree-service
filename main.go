package main

import (
	"encoding/json"
	"net/http"
)

var byt = []byte(`{"age": "test", "job": {"title": "engineer"}, "jage": 27, "name": "bob", "friends": ["one", "two", "three", "five"]}`)

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	var obj interface{}
	json.Unmarshal(byt, &obj)
	js, err := json.MarshalIndent(obj, "", "   ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
