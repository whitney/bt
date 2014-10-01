package bt

import (
    "crypto/sha1"
    "encoding/gob"
    "fmt"
    "io"
    "log"
    "net"
    "net/http"
    "net/url"
    "strconv"

    "github.com/marksamman/bencode"
)

const (
    peerId = "-WZ0001-373bff40fe0e"
    host = "localhost"
    port   = "6881"
)

type client struct {
    torrent io.Reader
    announce string
    infoHash map[string]interface{}
    totalBytes int64
    uploadedBytes int64
    downloadedBytes int64
    server net.Listener
}

func New(torrent io.Reader) *client {
    c := new(client)
    c.torrent = torrent

    dict, err := bencode.Decode(torrent)
    if err != nil {
        log.Fatal(err)
    }

    c.announce = dict["announce"].(string)
    c.infoHash = dict["info"].(map[string]interface{})
    c.totalBytes = torrentSize(c.infoHash)
    c.uploadedBytes = 0
    c.downloadedBytes = 0

    return c
}

// in bytes
func torrentSize(infoHash map[string]interface{}) int64 {
    length, ok := infoHash["length"]

    // single file mode
    if ok {
        return length.(int64)
    }

    var numBytes int64 = 0

    // multi file mode
    for _, file := range infoHash["files"].([]interface{}) {
        numBytes += file.(map[string]interface{})["length"].(int64)
    }  

    return numBytes
}

// tracker request params described here:
// https://wiki.theory.org/BitTorrentSpecification#Tracker_Request_Parameters
func (c *client) trackerURL() string {
    infoHash := bencode.Encode(c.infoHash)

    remainingBytes := c.totalBytes - c.downloadedBytes 

    hasher := sha1.New()
    hasher.Write(infoHash)
    sha := string(hasher.Sum(nil))

    params := url.Values{}
    
    params.Set("info_hash", sha)
    params.Set("peer_id", peerId)
    params.Set("port", port)

    params.Set("compact", "1")
    params.Set("event", "started")
    
    params.Set("uploaded", strconv.FormatInt(c.uploadedBytes, 10))
    params.Set("downloaded", strconv.FormatInt(c.downloadedBytes, 10))
    params.Set("left", strconv.FormatInt(remainingBytes, 10))

    return c.announce + "?" + params.Encode()
}

func (c *client) Start() {

    // start server
    c.startServer()

    // connect to peers

    // do handshakes

    // run until:
    // c.downloadedBytes == c.totalBytes

    fmt.Println(c.trackerRequest()) 
}

func (c *client) startServer() {
    ln, err := net.Listen("tcp", host + ":" + port)
    if err != nil {
        // handle error
    }

    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            // handle error
            continue
        }
        go handleConn(conn)
    }
}

func (c *client) stopServer() {
    c.server.Close()
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

func sendToPeer() {
    fmt.Println("start client");
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatal("Connection error", err)
    }
    encoder := gob.NewEncoder(conn)
    p := &P{1, 2}
    encoder.Encode(p)
    conn.Close()
    fmt.Println("done");
}

type trackerResponse struct {
    FailureReason string
    Interval int64
    MinInterval int64
    Complete int64
    Incomplete int64
    Peers []string
}

func (c *client) trackerRequest() *trackerResponse {
    resp, err := http.Get(c.trackerURL())
    if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close()

    fmt.Printf("tracker request status code: %s\n", resp.Status)

    dict, err := bencode.Decode(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    tr := new(trackerResponse)

    failureReason := dict["failure reason"]
    if failureReason != nil {
        tr.FailureReason = failureReason.(string)
    }

    interval := dict["interval"]
    if interval != nil {
        tr.Interval = interval.(int64)
    }

    minInterval := dict["min interval"]
    if minInterval != nil {
        tr.MinInterval = minInterval.(int64)
    }

    complete := dict["complete"]
    if complete != nil {
        tr.Complete = complete.(int64)
    }

    incomplete := dict["incomplete"]
    if incomplete != nil {
        tr.Incomplete = incomplete.(int64)
    }

    var binPeers string
    var isBinaryPeers bool
    var peers []string

    // determine if peers field is disctionary model or binary model
    dictPeers, isDictPeers := dict["peers"].(map[string]interface{})
    if isDictPeers {
        // TODO handle dictionary peers
        fmt.Println("dictPeers:", dictPeers)
    } else {
        binPeers, isBinaryPeers = dict["peers"].(string)
        if !isBinaryPeers {
            log.Fatal("invalid peers field")
        }

        low := 0
        high := 6
        
        for high <= len(binPeers) {
            peers = append(peers, ParsePeer([]byte(binPeers[low:high])))
            low += 6
            high += 6
        }
    } 

    tr.Peers = peers

    return tr
}

// returns <ip address>:<port>
func ParsePeer(peer []byte) string {
    peerIp := ""
    ipBytes := peer[0:4]

    for i := 0; i < len(ipBytes); i++ {
        ipComponent := int64(ipBytes[i])
        peerIp += strconv.FormatInt(ipComponent, 10)

        if i < len(ipBytes) - 1 {
            peerIp += "."
        }
    }

    portBytes := peer[4:6]
    port := strconv.FormatInt(256 * int64(portBytes[0]) + int64(portBytes[1]), 10)

    return peerIp + ":" + port
}
