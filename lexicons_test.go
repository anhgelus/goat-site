package site_test

import (
	"crypto/sha256"
	"net"
	"net/http"
	"time"

	"pgregory.net/rapid"
	"tangled.org/anhgelus.world/xrpc"
	"tangled.org/anhgelus.world/xrpc/atproto"
)

var (
	rapidLowerRunes = rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyz"))
)

func genCID(t *rapid.T, label string) *atproto.CID {
	cid := &atproto.CID{
		Version:  atproto.CIDVersion,
		Codec:    atproto.CIDCodecRaw,
		HashType: atproto.CIDHashSha256,
		HashSize: 32,
	}
	str := rapid.StringN(64, -1, -1).Draw(t, label)
	cp := make([]byte, 32)
	for i, v := range sha256.Sum256([]byte(str)) {
		cp[i] = v
	}
	cid.Digest = cp
	return cid
}

var dir *atproto.Directory

func getClient() xrpc.Client {
	if dir == nil {
		dir = atproto.NewDirectory(http.DefaultClient, net.DefaultResolver, 5*time.Minute)
	}
	return xrpc.NewClient(http.DefaultClient, dir)
}

func genTime(t *rapid.T, label string) time.Time {
	return time.Unix(int64(rapid.Uint32().Draw(t, label)), 0)
}
