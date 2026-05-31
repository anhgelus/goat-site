package site_test

import (
	"net"
	"net/http"
	"time"

	"anhgelus.world/xrpc"
	"anhgelus.world/xrpc/atproto"
	"pgregory.net/rapid"
)

var (
	rapidLowerRunes = rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyz"))
)

var dir atproto.Directory

func getClient() xrpc.Client {
	if dir == nil {
		dir = atproto.NewDirectory(http.DefaultClient, net.DefaultResolver)
	}
	return xrpc.NewClient(http.DefaultClient, dir, "GoAT Site tests v0.1.2 (Linux; Tangled Spindle; +https://anhgelus.world/goat-site/)")
}

func genTime(t *rapid.T, label string) time.Time {
	return time.Unix(int64(rapid.Uint32().Draw(t, label)), 0)
}
