package site_test

import (
	"context"
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

func getClient(t rapid.TB, test string) (syntax.ATURI, *atclient.APIClient) {
	dir := identity.DefaultDirectory()
	uri, err := syntax.ParseATURI(test)
	if err != nil {
		t.Fatal(err)
	}
	var id *identity.Identity
	id, err = dir.Lookup(context.Background(), uri.Authority())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("using", id.PDSEndpoint(), "for", test)
	client := atclient.NewAPIClient(id.PDSEndpoint())
	t.Log(uri.Authority().String(), uri.RecordKey())
	return uri, client
}

func genTime(t *rapid.T, label string) time.Time {
	return time.Unix(int64(rapid.Uint32().Draw(t, label)), 0)
}
