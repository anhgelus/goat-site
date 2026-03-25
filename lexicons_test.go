package site_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/bluesky-social/indigo/atproto/atclient"
	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"pgregory.net/rapid"
	site "tangled.org/anhgelus.world/goat-site"
)

var (
	rapidLowerRunes = rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyz"))
)

func genBlob(t *rapid.T, baseMime, label string) (*site.Blob, map[string]any) {
	blob := &site.Blob{
		CID: rapid.StringN(2, -1, 128).Draw(t, label+"_cid"),
		MimeType: baseMime + "/" +
			rapid.StringOfN(rapidLowerRunes, 2, 20, -1).Draw(t, label+"_mimeType"),
		Size: rapid.UintMin(1).Draw(t, label+"_size"),
	}
	return blob, map[string]any{
		"$type":    blob.Type(),
		"ref":      map[string]any{"$link": blob.CID},
		"mimeType": blob.MimeType,
		"size":     blob.Size,
	}
}

func genURL(t *rapid.T, label string) string {
	scheme := "http"
	if rapid.Bool().Draw(t, label+"_secure?") {
		scheme += "s"
	}
	base := rapid.StringOfN(rapidLowerRunes, 1, -1, 64).Draw(t, label+"_base")
	tld := rapid.StringOfN(rapidLowerRunes, 2, -1, 10).Draw(t, label+"_tld")
	sub := rapid.StringOfN(rapidLowerRunes, -1, -1, 32).Draw(t, label+"_sub")
	var sb strings.Builder
	sb.Grow(len(base) + len(tld) + len(sub) + len(scheme) + 5)
	sb.WriteString(scheme)
	sb.WriteString("://")
	if sub != "" {
		sb.WriteString(sub)
		sb.WriteRune('.')
	}
	sb.WriteString(base)
	sb.WriteRune('.')
	sb.WriteString(tld)
	path := genPath(t, label+"_path")
	if path != "/" {
		sb.Grow(len(path))
		sb.WriteString(path)
	}
	return sb.String()
}

func genPath(t *rapid.T, label string) string {
	valid := rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	return "/" + rapid.StringOfN(valid, -1, -1, 64).Draw(t, label)
}

func genRecordKey(t *rapid.T, label string) string {
	valid := rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyz0123456789"))
	return rapid.StringOfN(valid, 1, -1, 128).Draw(t, label)
}

func getClient(t *testing.T, test string, uri *syntax.ATURI, client **atclient.APIClient) (syntax.ATURI, *atclient.APIClient) {
	var err error
	defer func() {
		if err == nil {
			t.Log(uri.Authority().String(), uri.RecordKey())
		}
	}()
	if *client != nil {
		return *uri, *client
	}
	dir := identity.DefaultDirectory()
	*uri, err = syntax.ParseATURI(test)
	if err != nil {
		t.Fatal(err)
	}
	var id *identity.Identity
	id, err = dir.Lookup(context.Background(), uri.Authority())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("using", id.PDSEndpoint(), "for", test)
	*client = atclient.NewAPIClient(id.PDSEndpoint())
	return *uri, *client
}

func genDid(t *rapid.T, label string) string {
	return "did:plc:" + rapid.StringOfN(rapidLowerRunes, 24, -1, 24).Draw(t, label)
}

func genTime(t *rapid.T, label string) time.Time {
	return time.Unix(int64(rapid.Uint32().Draw(t, label)), 0)
}
