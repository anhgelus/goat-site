package site_test

import (
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

func genBlob(t *rapid.T, baseMime, label string) (*xrpc.Blob, map[string]any) {
	blob := &xrpc.Blob{
		CID: rapid.StringN(2, -1, 128).Draw(t, label+"_cid"),
		MimeType: baseMime + "/" +
			rapid.StringOfN(rapidLowerRunes, 2, 20, -1).Draw(t, label+"_mimeType"),
		Size: rapid.UintMin(1).Draw(t, label+"_size"),
	}
	return blob, map[string]any{
		"$type":    blob.Collection(),
		"ref":      map[string]any{"$link": blob.CID},
		"mimeType": blob.MimeType,
		"size":     blob.Size,
	}
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
