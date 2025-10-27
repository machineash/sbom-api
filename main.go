package main

import (
	"log"
	"net/http"
	"sbom-api/api"
)

func main() {
	h := api.NewHandlers()

	http.HandleFunc("/components", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			// if id is present, fetch by id; else.. list
			if r.URL.Query().Get("id") != "" {
				h.GetByID(w, r)
				return
			}
			h.List(w, r)
		case http.MethodPut:
			h.Put(w, r)
		case http.MethodPatch:
			h.Patch(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
