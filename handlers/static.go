package handlers

import (
	"net/http"
	"log"
)



func ServeWebapp(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving static files")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "webapp/dist/index.html")
}
