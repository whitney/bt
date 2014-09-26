package main

import (
    //"errors"
    //"flag"
    "crypto/sha1"
    "encoding/base64"
    "fmt"
    "io"
    "log"
    "net/url"
    "os"
    "strconv"

    "github.com/marksamman/bencode"
)

const (
    peerId = "-WZ0001-373bff40fe0e"
    port   = "6881"
)

type client struct {
    torrent io.Reader
    announce string
    infoHash map[string]interface{}
    totalBytes int64
    uploadedBytes int64
    downloadedBytes int64
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

    //files := infoHash["files"].()

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
    sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

    params := url.Values{}
    
    params.Set("info_hash", sha)
    params.Set("peer_id", peerId)
    params.Set("port", port)
    params.Set("port", port)

    params.Set("compact", "1")
    params.Set("event", "started")
    
    params.Set("uploaded", strconv.FormatInt(c.uploadedBytes, 10))
    params.Set("downloaded", strconv.FormatInt(c.downloadedBytes, 10))
    params.Set("left", strconv.FormatInt(remainingBytes, 10))

    return c.announce + "?" + params.Encode()
}

func main() {
    file, err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    c := New(file)

    fmt.Printf("announce (tracker base URL): %s\n", c.announce)

    fmt.Printf("tracker URL: %s\n", c.trackerURL())
}
