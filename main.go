package main

import (
    //"errors"
    //"flag"
    "crypto/sha1"
    "encoding/base64"
    "fmt"
    "io"
    "log"
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
}

/*
func (c *client) area() int {
    return r.width * r.height
}
*/

func New(torrent io.Reader) *client {
    c := new(client)
    c.torrent = torrent
    /*
    s.CurrentVersion = defaultVersion
    s.Root = newDir(s, "/", s.CurrentIndex, nil, "", Permanent)
    s.Stats = newStats()
    s.WatcherHub = newWatchHub(1000)
    s.ttlKeyHeap = newTtlKeyHeap()
    */
    return c
}

// tracker request params described here:
// https://wiki.theory.org/BitTorrentSpecification#Tracker_Request_Parameters
func trackerURL(baseURL string, infoHash []byte) (string, error) {

    downloadedBytes := 0
    totalBytes := 0
    left := totalBytes - downloadedBytes 

    hasher := sha1.New()
    hasher.Write(infoHash)
    sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
    
    baseURL += "?info_hash=" + sha
    baseURL += "&peer_id=" + peerId
    baseURL += "&port=" + port
    baseURL += "&uploaded=0"
    baseURL += "&downloaded=" + strconv.Itoa(downloadedBytes)
    baseURL += "&left=" + strconv.Itoa(left)
    baseURL += "&compact=1"
    baseURL += "&event=started"

    return baseURL, nil
}

func main() {
    file, err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    dict, err := bencode.Decode(file)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("announce (tracker base URL): %s\n", dict["announce"].(string))

    infoHash := bencode.Encode(dict["info"])

    trackerURL, err := trackerURL(dict["announce"].(string), infoHash)
    fmt.Printf("tracker URL: %s\n", trackerURL)
}
