package main

import (
  "fmt"
  "log"
  "net/http"

  "github.com/moorara/konfig"
)

var config = struct {
  AuthToken string
} {}

func main() {
  konfig.Pick(&config, konfig.Telepresence(), konfig.Debug(3))
  log.Printf("making service-to-service calls using this token: %s", config.AuthToken)

  http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
  })

  log.Fatal(http.ListenAndServe(":8080", nil))
}
