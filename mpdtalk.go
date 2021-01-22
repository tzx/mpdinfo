package main

import (
	"fmt"
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
  Elapsed   string    `json:"elapsed"`
  Duration  string    `json:"duration"`
  State     string    `json:"state"`
}

func GetMPDInfo() (Info, error) {
  conn, err := textproto.Dial("tcp", config.MPDTCPURL)
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

  forJSON := Info {
    Artist: info["Artist"],
    Title: info["Title"],
    Album: info["Album"],
    State:  info["state"],
    Elapsed: secondsToMinutes(elapsed),
    Duration: secondsToMinutes(duration),
  }
  return forJSON
}

func secondsToMinutes(toConvert float64) string {
  rounded :=  uint(math.Round(toConvert))
  minutes := rounded / 60
  seconds := rounded % 60
  res := fmt.Sprintf("%d:%02d", minutes, seconds)
  return res
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
