package main

import (
    "log"
    "os"
)

func main() {
    file, err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    c := NewClient(file)

    c.Start()
}
