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
  "strings"
  "strconv"
  "encoding/json"
  "go.bug.st/serial"
)

const (
  ResetReason  = "?R\r\n"
  PartNumber   = "?P\r\n"
  SerialNumber = "?N\r\n"
  MilesPerHour = "US\r\n"
  
  Reset = "{\"Reset\" : \"Board was reset."
  USBActive = "{\"USB\" : \"USB Interface Active\""
)

func B(s string) []byte {
  return []byte(s)
}


func readPort(port serial.Port) string {
  buff := make([]byte, 128)
  n, _ := port.Read(buff)
  opsStr   := strings.TrimSpace(string(buff[:n]))
  
  return opsStr  
}

func readPortJSON(port serial.Port) string {
  buff := make([]byte, 128)
  n    := 0

  buff[0] = 0

  for buff[0] != '{' {
    n, _ = port.Read(buff)
  }
  opsStr   := string(buff[:n])
  return opsStr  
}

var OPS243 struct {
  product string
  serial string
}

var partNumber struct {
  Product string `json:"Product"`
}

var serialNumber struct {
  SerialNumber string `json:"SerialNumber"`
}


func initOPS243(port serial.Port) {

  port.Write(B(PartNumber))

  response := readPortJSON(port)
  fmt.Print(response)
  if err := json.Unmarshal([]byte(response), &partNumber); err != nil {
    log.Fatal("Fatal:  ", err)
  }

  OPS243.product = partNumber.Product

  log.Print("Get serial number")

  port.Write(B(SerialNumber))
  response = readPortJSON(port)
  fmt.Print(response)
  if err := json.Unmarshal([]byte(response), &serialNumber); err != nil {
    log.Fatal("Fatal:  ", err)
  }

  OPS243.serial = serialNumber.SerialNumber

  // Set output units to miles per hour
  //port.Write(B(MilesPerHour))

  //milesPerHour := readPort(port)
  fmt.Printf("Product:  %s, Serial:  %s",  OPS243.product, OPS243.serial)
  
  //fmt.Println(milesPerHour)
 
}

func main() {

  mode := &serial.Mode {
    BaudRate:  115200,
  }

  port, err := serial.Open("/dev/ttyACM0", mode)

  if err != nil {
    log.Fatal(err)
  }

  ready := false
  for !ready {
    reading := readPort(port)
    _, err := strconv.ParseFloat(reading, 64)
    if err == nil {
      fmt.Println("Receiving OPS243 readings, ready!")
      ready = true 
    } else {
      fmt.Println("Synchronizing")
    }
  }

  initOPS243(port)

  // Get speed
  buff := make([]byte, 128)
  for {
    n, _ := port.Read(buff)
    fmt.Printf("%s", string(buff[:n]))
  }

  fmt.Println("Exiting")
}
