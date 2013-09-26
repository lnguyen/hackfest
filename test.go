package main

import(
    "log"
    "github.com/GeertJohan/go.hid"
    "fmt"
    "time"
    "os/exec"
    "net/http"
)

var coinFound bool
var coinEject bool

func readCoin() {
    // open the MSI GT leds device
    leds, err := hid.Open(0x2341, 0x8036, "")
    if err != nil {
      log.Fatalf("Could not open leds device: %s", err)
    }
    defer leds.Close()
    output, _ := leds.ManufacturerString()
    fmt.Println(output)
    coinDetected := false

    for i := 0; i < 10000000000; i++ {
      buf := make([]byte, 1024)
      _, err = leds.Read(buf)
      if err != nil {
        log.Fatalf("%s", err)
      }
      if buf[3] != 0 {
        if coinDetected == false {
          coinDetected = true
          coinFound = true
          fmt.Println("coin inserted")
          lsCmd := exec.Command("say", "coin inserted")
          _, err = lsCmd.Output()
          if err != nil {
            panic(err)
          }
        }
      } else {
        if coinDetected == true {
          coinDetected = false
          coinEject = true
          fmt.Println("coin ejected")
          lsCmd := exec.Command("say", "coin removed")
          _, err = lsCmd.Output()
          if err != nil {
            panic(err)
          }
        }
      }
      time.Sleep(100 * time.Millisecond)
    }
}

func CoinHandler(w http.ResponseWriter, r *http.Request) {
  if coinFound {
    fmt.Fprint(w,`{ "coin": "Found" }`)
    coinFound = false
    return
  } else if coinEject {
    fmt.Fprint(w,`{ "coin": "Ejected" }`)
    coinEject = false
    return
  } else {
    fmt.Fprint(w,`{ "coin": "no data" }`)
    return
  }

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "index.html")
  return
}

func main() {
  go readCoin()
  http.HandleFunc("/", HomeHandler)
  http.HandleFunc("/coin", CoinHandler)
  http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:])
  })
  http.ListenAndServe(":5000",nil)
}
