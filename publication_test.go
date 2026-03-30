package site_test

import (
	"context"
	"encoding/json"
	"testing"

	"pgregory.net/rapid"
	site "tangled.org/anhgelus.world/goat-site"
	"tangled.org/anhgelus.world/xrpc"
	"tangled.org/anhgelus.world/xrpc/atproto"
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
		"$type":            theme.Collection(),
		"accent":           colors("accent", &theme.Accent),
		"accentForeground": colors("accentForeground", &theme.AccentForeground),
		"foreground":       colors("foreground", &theme.Foreground),
		"background":       colors("background", &theme.Background),
	}
}

func TestPublication_JSON(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		theme, themeRaw := genBasicTheme(t)
		description := rapid.StringN(2, 3_000, 30_000).Draw(t, "description")
		name := rapid.StringN(2, 500, 5_000).Draw(t, "name")
		url := genURL(t, "url")
		showInDiscover := rapid.Bool().Draw(t, "showInDiscover")
		input := map[string]any{
			"$type":       site.CollectionPublication,
			"basicTheme":  themeRaw,
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
		var pub *site.Publication
		err = json.Unmarshal(b, &pub)
		if err != nil {
			t.Fatal(err)
		}
		if pub.Name != name {
			t.Errorf("invalid name: %s, wanted %s", pub.Name, name)
		}
		if pub.URL.String() != url {
			t.Errorf("invalid url: %s, wanted %s", pub.URL, url)
		}
		if *pub.Description != description {
			t.Errorf("invalid description: %s, wanted %s", *pub.Description, description)
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
		b, err = xrpc.Marshal(pub)
		if err != nil {
			t.Fatal(err)
		}
	})
}

var genPubAt = []string{
	"at://did:plc:revjuqmkvrw6fnkxppqtszpv/site.standard.publication/3md4kftpfxs2z", // leaflet pub
	"at://did:plc:yk4dd2qkboz2yv6tpubpc6co/site.standard.publication/3m6zrpzbs3s2y", // pckt pub
}

func TestGetPublication(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	for _, uri := range genPubAt {
		client := getClient()
		u, err := atproto.ParseURI(context.Background(), client.Directory(), uri)
		if err != nil {
			t.Fatal(err)
		}
		union, err := client.FetchURI(context.Background(), u)
		if err != nil {
			t.Fatal(err)
		}
		pub := new(site.Publication)
		if !union.Value.As(pub) {
			t.Fatalf("cannot convert union to publication: %s", union.Value.Raw)
		}
		v, err := pub.Verify(context.Background(), client.HTTP(), u.Authority(), *u.RecordKey())
		if err != nil {
			t.Fatal(err)
		}
		if !v {
			t.Errorf("cannot verify %s", uri)
		}
	}
}

func TestListPublications(t *testing.T) {
	if testing.Short() {
		t.Skip("not doing http requests in short")
	}
	for _, uri := range genPubAt {
		client := getClient()
		u, err := atproto.ParseURI(context.Background(), client.Directory(), uri)
		if err != nil {
			t.Fatal(err)
		}
		pubs, _, err := xrpc.ListRecords[*site.Document](context.Background(), client, u.Authority(), 0, "", false)
		if err != nil {
			t.Fatal(err)
		}
		if pubs == nil {
			t.Errorf("pubs is nil")
		}
		for i, pub := range pubs {
			if pub.Value == nil {
				t.Errorf("pub %d is nil", i)
			}
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
