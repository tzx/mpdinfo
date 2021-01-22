package main

import (
	"encoding/json"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  info, err := GetMPDInfo()
  if err != nil {
    w.WriteHeader(http.StatusServiceUnavailable)
    res, _ := json.Marshal("MPD is not on or having issues!")
    w.Write(res)
    return
  }
  json.NewEncoder(w).Encode(info)
}

func main() {
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}
