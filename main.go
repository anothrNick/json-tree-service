package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var byt = []byte(`{"age": "test", "job": {"title": "engineer"}, "jage": 27, "name": "bob", "friends": ["one", "two", "three", "five"]}`)

func main() {
	router := gin.Default()

	router.GET("/api/*keys", func(c *gin.Context) {
		keys := c.Param("keys")
		c.String(http.StatusOK, keys)
	})

	router.Run(":5000")
}

// func main() {
// 	http.HandleFunc("/", foo)
// 	http.ListenAndServe(":3000", nil)
// }

// func foo(w http.ResponseWriter, r *http.Request) {
// 	var obj interface{}
// 	json.Unmarshal(byt, &obj)
// 	js, err := json.MarshalIndent(obj, "", "   ")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }
