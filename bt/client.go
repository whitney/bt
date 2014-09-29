package bt

import (
    "crypto/sha1"
    "fmt"
    "io"
    //"io/ioutil"
    "log"
    "net/http"
    "net/url"
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

func (c *client) TrackerRequest() string {
    resp, err := http.Get(c.trackerURL())
    if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close()

    /*
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    */

    fmt.Printf("tracker request status code: %s\n", resp.Status)

    dict, err := bencode.Decode(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("TrackerRequest dict:", dict)

    failureReason := dict["failure reason"]
    if failureReason != nil {
        failureReason = failureReason.(string)
    }
    fmt.Printf("failureReason: %s\n", failureReason)

    interval := dict["interval"]
    if interval != nil {
        interval = interval.(int64)
    }
    fmt.Printf("interval: %s\n", interval)

    minInterval := dict["min interval"]
    if minInterval != nil {
        minInterval = minInterval.(int64)
    }
    fmt.Printf("minInterval: %s\n", minInterval)

    complete := dict["complete"]
    if complete != nil {
        complete = complete.(int64)
    }
    fmt.Printf("complete: %s\n", complete)

    incomplete := dict["incomplete"]
    if incomplete != nil {
        incomplete = incomplete.(int64)
    }
    fmt.Printf("incomplete: %s\n", incomplete)

    var binPeers string
    var isBinaryPeers bool

    // determine if peers field is disctionary model or binary model
    dictPeers, isDictPeers := dict["peers"].(map[string]interface{})
    if !isDictPeers {
        binPeers, isBinaryPeers = dict["peers"].(string)
        if !isBinaryPeers {
            log.Fatal("invalid peers field")
        }
    } 

    fmt.Println("dictPeers:", dictPeers)
    fmt.Println("binPeers:", binPeers)

    //return string(body)
    return "pants"
}
