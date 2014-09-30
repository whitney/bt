package bt

import (
    "testing"
)

func TestParsePeer(t *testing.T) {
    expectedIp := "209.131.52.148:6881"
    peerBytes := []byte{0xd1, 0x83, 0x34, 0x94, 0x1a, 0xe1}
    peerIp := ParsePeer(peerBytes)
    if peerIp != expectedIp {
        t.Error("peerIp " + peerIp)
    }
}
