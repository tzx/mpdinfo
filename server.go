package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Config struct {
  MPDHTTPURL  string
  MPDTCPURL   string
  ServerPort  string
}

type TemplateInfo struct {
  MPDHTTPURL  string
  StatusURL   string
}

var config Config
const statusEndpoint = "/status"

func getMPDInfo(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  info, err := GetMPDInfo()
  if err != nil {
    log.Println(err)
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

  templating := TemplateInfo {
    MPDHTTPURL: config.MPDHTTPURL,
    StatusURL: "http://localhost:" + config.ServerPort + statusEndpoint,
  }
  t.Execute(w, templating)
}

func main() {
  http.HandleFunc(statusEndpoint, getMPDInfo)
  http.HandleFunc("/", index)
  http.ListenAndServe(":" + config.ServerPort, nil)
}

func init() {
  file, err := os.Open("config.json")
  if err != nil {
    log.Fatal(err)
  }
  decoder := json.NewDecoder(file)
  err = decoder.Decode(&config)
  if err != nil {
    log.Fatal(err)
  }
  log.Println(config)
}
