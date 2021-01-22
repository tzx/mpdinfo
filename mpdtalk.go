package main

import (
	"log"
	"math"
	"net/textproto"
	"strconv"
	"strings"
)

type Info struct {
  Artist    string    `json:"artist"`
  Title     string    `json:"title"`
  Album     string    `json:"album"`
  Elapsed   float64   `json:"elapsed"`
  Duration  float64   `json:"duration"`
  State     string    `json:"state"`
}

func GetMPDInfo() (Info, error) {
  // TODO: localhost 6600 variable
  conn, err := textproto.Dial("tcp", "localhost:6600")
  if err != nil {
    var res Info
    return res, err
  }
  defer conn.Close()

  okmsg, err :=  conn.ReadLine()
  if err != nil || okmsg[:6] != "OK MPD" {
    log.Fatal("Error reading", err)
  }

  cmd := "command_list_begin\nstatus\ncurrentsong\ncommand_list_end"
  return sendInfoCmd(conn, cmd), nil
}

func sendInfoCmd(c *textproto.Conn, s string) Info {
  id, err := c.Cmd(s)
  if err != nil {
    log.Fatal("Error sending Cmd", err)
  }
  c.StartResponse(id)
  defer c.EndResponse(id)
  info := read(c, "OK")

  if info["state"] == "stop" {
    var empty Info
    return empty
  }

  elapsed, err := strconv.ParseFloat(info["elapsed"], 64)
  duration, err := strconv.ParseFloat(info["duration"], 64)
  if err != nil {
    log.Fatal("Uhhh mpd didn't give a number?", err)
  }
  elapsed = math.Round(elapsed)
  duration = math.Round(duration)

  forJSON := Info {
    Artist: info["Artist"],
    Title: info["Title"],
    Album: info["Album"],
    State:  info["state"],
    Elapsed: elapsed,
    Duration: duration,
  }
  return forJSON
}

func read(c *textproto.Conn, end string) map[string]string {
  m := make(map[string]string)
  for {
    line, err := c.ReadLine()
    if err != nil || line == end {
      break
    }
    splitted := strings.Split(line, ": ")
    k, v := splitted[0], splitted[1]
    m[k] = v
  }
  return m
}
