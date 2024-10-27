/* 
Copyright 2024 iAchieved.it LLC

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above 
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY 
AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE
OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THIS SOFTWARE.

SPDX-License-Identifier: ISC
*/

package main

import (
  "fmt"
  "log"
  //"encoding/json"
  "go.bug.st/serial"
)

const (
  PartNumber   = "?P\r\n"
  SerialNumber = "?N\r\n"
  MilesPerHour = "US\r\n"
  
)

func B(s string) []byte {
  return []byte(s)
}

func initOPS243(port serial.Port) {

  buff := make([]byte, 512)
  n    := 0

  port.Write(B(PartNumber))
  buff[0] = 0
  
  for buff[0] != '{' {
    n, _ = port.Read(buff)
  }

  fmt.Printf("%s", string(buff[:n]))

  port.Write(B(SerialNumber))
  buff[0] = 0
  
  for buff[0] != '{' {
    n, _ = port.Read(buff)
  }

  fmt.Printf("%s", string(buff[:n]))

  // Set output units to miles per hour
  port.Write(B(MilesPerHour))
  buff[0] = 0
  
  for buff[0] != '{' {
    n, _ = port.Read(buff)
  }

  fmt.Printf("%s", string(buff[:n]))

  
}

func main() {

  mode := &serial.Mode {
    BaudRate:  115200,
  }

  port, err := serial.Open("/dev/ttyACM0", mode)

  if err != nil {
    log.Fatal(err)
  }

  initOPS243(port)

  // Get speed
  buff := make([]byte, 512)
  for {
    n, _ := port.Read(buff)
    fmt.Printf("%s", string(buff[:n]))
  }

  fmt.Println("Exiting")
}
