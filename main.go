package main

import (
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

    c.Start()
}
