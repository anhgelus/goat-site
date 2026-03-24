package site_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/bluesky-social/indigo/atproto/atclient"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"pgregory.net/rapid"
	site "tangled.org/anhgelus.world/goat-site"
)

func genBasicTheme(t *rapid.T) (*site.Theme, map[string]any) {
	theme := new(site.Theme)
	colors := func(base string, rgbP **site.RGB) map[string]any {
		*rgbP = new(site.RGB)
		rgb := *rgbP
		rgb.Red = rapid.Uint8().Draw(t, base+"_r")
		rgb.Green = rapid.Uint8().Draw(t, base+"_g")
		rgb.Blue = rapid.Uint8().Draw(t, base+"_b")
		return map[string]any{
			"$type": site.CollectionThemeColorRGB,
			"r":     rgb.Red,
			"g":     rgb.Green,
			"b":     rgb.Blue,
		}
	}
	return theme, map[string]any{
		"$type":            theme.Type(),
		"accent":           colors("accent", &theme.Accent),
		"accentForeground": colors("accentForeground", &theme.AccentForeground),
		"foreground":       colors("foreground", &theme.Foreground),
		"background":       colors("background", &theme.Background),
	}
}

func TestPublication_JSON(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		theme, themeRaw := genBasicTheme(t)
		icon, iconRaw := genBlob(t, "image")
		description := rapid.StringN(2, 3_000, 30_000).Draw(t, "description")
		name := rapid.StringN(2, 500, 5_000).Draw(t, "name")
		url := genURL(t)
		showInDiscover := rapid.Bool().Draw(t, "showInDiscover")
		input := map[string]any{
			"$type":       site.CollectionPublication,
			"basicTheme":  themeRaw,
			"icon":        iconRaw,
			"description": description,
			"name":        name,
			"url":         url,
			"preferences": map[string]any{"showInDiscover": showInDiscover},
		}
		b, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))
		var v *site.RecordJSON
		err = json.Unmarshal(b, &v)
		if err != nil {
			t.Fatal(err)
		}
		pub := v.Record.(*site.Publication)
		if pub.Name != name {
			t.Errorf("invalid name: %s, wanted %s", pub.Name, name)
		}
		if pub.URL.String() != url {
			t.Errorf("invalid url: %s, wanted %s", pub.URL, url)
		}
		if *pub.Description != description {
			t.Errorf("invalid description: %s, wanted %s", *pub.Description, description)
		}
		if pub.Icon.CID != icon.CID {
			t.Errorf("invalid Icon CID: %s, wanted %s", pub.Icon.CID, icon.CID)
		}
		if pub.Icon.MimeType != icon.MimeType {
			t.Errorf("invalid Icon MimeType: %s, wanted %s", pub.Icon.MimeType, icon.MimeType)
		}
		if pub.Icon.Size != icon.Size {
			t.Errorf("invalid Icon Size: %d, wanted %d", pub.Icon.Size, icon.Size)
		}
		if pub.Preferences.ShowInDiscover != showInDiscover {
			t.Errorf("invalid Preferences ShowInDiscover: %v", pub.Preferences.ShowInDiscover)
		}
		th := pub.BasicTheme
		if *th.Accent != *theme.Accent {
			t.Errorf("invalid theme accent color: %s, wanted %s", th.Accent, theme.Accent)
		}
		if *th.AccentForeground != *theme.AccentForeground {
			t.Errorf("invalid theme accent foreground color: %s, wanted %s", th.AccentForeground, theme.AccentForeground)
		}
		if *th.Background != *theme.Background {
			t.Errorf("invalid theme background color: %s, wanted %s", th.Background, theme.Background)
		}
		if *th.Foreground != *theme.Foreground {
			t.Errorf("invalid theme foreground color: %s, wanted %s", th.Foreground, theme.Foreground)
		}
		b, err = json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
	})
}

// leaflet publication
// const testPub = "at://did:plc:yk4dd2qkboz2yv6tpubpc6co/site.standard.publication/3m6zrpzbs3s2y"

// pckt publication (because they do not use the preferred time format!)
const testPub = "at://did:plc:revjuqmkvrw6fnkxppqtszpv/site.standard.publication/3md4kftpfxs2z"

var (
	pubURI    syntax.ATURI
	pubClient *atclient.APIClient
)

func TestGetPublication(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	uri, client := getClient(t, testPub, &pubURI, &pubClient)
	pub, err := site.GetRecord[*site.Publication](context.Background(), client, uri.Authority(), uri.RecordKey())
	if err != nil {
		t.Fatal(err)
	}
	if pub == nil {
		t.Errorf("pub is nil")
	}
}

func TestListPublications(t *testing.T) {
	if testing.Short() {
		t.Skip("not doing http requests in short")
	}
	uri, client := getClient(t, testPub, &pubURI, &pubClient)
	pubs, _, err := site.ListRecords[*site.Publication](context.Background(), client, uri.Authority(), "", false)
	if err != nil {
		t.Fatal(err)
	}
	if pubs == nil {
		t.Errorf("pubs is nil")
	}
	for i, pub := range pubs {
		if pub == nil {
			t.Errorf("pub %d is nil", i)
		}
	}
}

func TestPublicationVerification(t *testing.T) {
	uri := site.GetPublicationVerificationURI("")
	if uri != "/.well-known/site.standard.publication" {
		t.Errorf("invalid uri: %s", uri)
	}
	uri = site.GetPublicationVerificationURI("/path/to/publication")
	if uri != "/.well-known/site.standard.publication/path/to/publication" {
		t.Errorf("invalid uri: %s", uri)
	}
	uri = site.GetPublicationVerificationURI("path/to/publication")
	if uri != "/.well-known/site.standard.publication/path/to/publication" {
		t.Errorf("invalid uri: %s", uri)
	}
}

func TestPublication_Verify(t *testing.T) {
	if testing.Short() {
		t.Skip("not doing http requests in short")
	}
	id, client := getClient(t, testPub, &pubURI, &pubClient)
	pub, err := site.GetRecord[*site.Publication](context.Background(), client, id.Authority(), id.RecordKey())
	if err != nil {
		t.Fatal(err)
	}
	v, err := pub.Verify(context.Background(), client.Client, id.Authority(), id.RecordKey())
	if err != nil {
		t.Fatal(err)
	}
	if !v {
		t.Errorf("cannot verify %s", id)
	}
}
