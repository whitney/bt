package main

import (
    "encoding/gob"
    "fmt"
    "log"
    "net"
)

func main() {
    //host := "209.131.52.148:6881"
    host := "localhost:6881"

    fmt.Println("initiating testTCP conn with ");

    conn, err := net.Dial("tcp", host)
    if err != nil {
        log.Fatal("Connection error", err)
    }
    encoder := gob.NewEncoder(conn)

    msg := "SWEET PANTS MSG"

    encoder.Encode(msg)

    handleConn(conn)

    conn.Close()
    fmt.Println("done");
}

type P struct {
    M, N int64
}
func handleConn(conn net.Conn) {
    dec := gob.NewDecoder(conn)
    p := &P{}
    dec.Decode(p)
    fmt.Printf("Received : %+v", p);
    conn.Write([]byte("Message received."))
}
