package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type TemplateInfo struct {
  MPDURL string
  ServerURL string
}

var URLs TemplateInfo

func getMPDInfo(w http.ResponseWriter, r *http.Request) {
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

func index(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("./index.html")
  if err != nil {
    log.Fatal("Template Error:", err)
  }
  t.Execute(w, URLs)
}

func main() {
  http.HandleFunc("/status", getMPDInfo)
  http.HandleFunc("/", index)
  http.ListenAndServe(":8080", nil)
}

func init() {
  argsWithoutProg := os.Args[1:]
  if len(argsWithoutProg) != 2 {
    log.Fatal("Usage: mpdinfo mpdURL serverURL")
  }
  URLs = TemplateInfo{MPDURL: argsWithoutProg[0], ServerURL: argsWithoutProg[1]}
}
