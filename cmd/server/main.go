package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Viswesh934/gotei/internal/engine"
)

type RenderRequest struct {
	HTML string `json:"html"`
}

func renderHandler(w http.ResponseWriter, r *http.Request) {
	var req RenderRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	pdf, _ := engine.Render(req.HTML)

	w.Header().Set("Content-Type", "application/pdf")
	w.Write(pdf)
}

func main() {
	http.HandleFunc("/render", renderHandler)
	log.Println("gotei running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
