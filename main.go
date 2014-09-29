package main

import (
    "fmt"
    "log"
    "os"

    "github.com/whitney/bt/bt"
)

func main() {
    file, err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    c := bt.New(file)

    fmt.Printf("tracker response: %s\n", c.TrackerRequest())
}
