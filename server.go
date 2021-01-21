package main

import (
  "fmt"
  "net/textproto"
  "log"
)

func main() {
  // TODO: localhost 6600 variable
  fmt.Println("Hello World")
  conn, err := textproto.Dial("tcp", "localhost:6600")
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  okmsg, err :=  conn.ReadLine()
  if err != nil || okmsg[:6] != "OK MPD" {
    log.Fatal("Error reading", err)
  }

  command := "command_list_begin\nstatus\ncurrentsong\ncommand_list_end"
  talk(conn, command)
}

func talk(c *textproto.Conn, s string) {
  id, err := c.Cmd(s)
  if err != nil {
    log.Fatal("Error sending Cmd", err)
  }
  c.StartResponse(id)
  defer c.EndResponse(id)
  read(c, "OK")
}

func read(c *textproto.Conn, end string) {
  for {
    line, err := c.ReadLine()
    if err != nil {
      log.Fatal("Error reading", err)
    }
    if line == end {
      break
    }
    fmt.Println("line: ", line)
  }
}
