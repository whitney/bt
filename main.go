package main

import (
    //"errors"
    //"flag"
    "crypto/sha1"
    "encoding/base64"
    "fmt"
    "log"
    "os"
    "strconv"

    "github.com/marksamman/bencode"
)

const (
    peerId = "-WZ0001-373bff40fe0e"
    port   = "6881"
)

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
