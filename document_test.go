package site_test

import (
	"context"
	"encoding/json"
	"slices"
	"testing"

	"pgregory.net/rapid"
	site "tangled.org/anhgelus.world/goat-site"
	"tangled.org/anhgelus.world/xrpc"
	"tangled.org/anhgelus.world/xrpc/atproto"
)

type content struct {
	Pages any `json:"pages"`
}

func (c *content) Collection() *atproto.NSID {
	return atproto.NewNSIDBuilder(`pub.leaflet`).Name(`content`).Build()
}

func TestDocument_JSON(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		var pubUrl string
		if rapid.Bool().Draw(t, "url_at?") {
			pubUrl = "at://" + genDid(t, "url_did") +
				"/" + site.CollectionPublication.String() +
				"/" + genRecordKey(t, "url_record_key")
		} else {
			pubUrl = genURL(t, "url")
		}
		title := rapid.StringN(1, 500, 5_000).Draw(t, "title")
		publishedAt := genTime(t, "published_at")
		path := genPath(t, "path")
		description := rapid.StringN(0, 3_000, 30_000).Draw(t, "description")
		coverImage, coverImageRaw := genBlob(t, "image", "cover_image")
		textContent := rapid.String().Draw(t, "text_content")
		tags := rapid.SliceOfN(rapid.String(), 0, 1280).Draw(t, "tags")
		updatedAt := genTime(t, "updated_at")
		input := map[string]any{
			"$type":       site.CollectionDocument,
			"site":        pubUrl,
			"title":       title,
			"publishedAt": publishedAt.Format(atproto.TimeFormat),
			"path":        path,
			"description": description,
			"coverImage":  coverImageRaw,
			"content":     json.RawMessage(`{"$type":"pub.leaflet.content","pages":[{"$type":"pub.leaflet.pages.linearDocument","blocks":[{"$type":"pub.leaflet.pages.linearDocument#block","block":{"$type":"pub.leaflet.blocks.text","plaintext":"hiiiiiiiii"}}],"id":"019d1297-2fdd-733b-9837-911e1758f300"}]}`),
			"textContent": textContent,
			"bskyPostRef": json.RawMessage(`{"cid":"bafyreidepvhssy3zglq3bo4nauszqhqmbk6lzzfay3r2nskvijyiewlr2u","commit":{"cid":"bafyreickwfv4p2jr6zvbdk6mldmddag2m6grpkbbvkvz57mvaqso5dpf5e","rev":"3mhm4oeyyzi2g"},"uri":"at://did:plc:jdhpqeb4cb4mng533dx56cbc/app.bsky.feed.post/3mhm4oevhmk2d","validationStatus":"valid"}`),
			"tags":        tags,
			"updatedAt":   updatedAt.Format(atproto.TimeFormat),
		}
		b, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))
		var doc *site.Document
		err = json.Unmarshal(b, &doc)
		if err != nil {
			t.Fatal(err)
		}
		var site string
		if doc.Site.IsAT() {
			site = doc.Site.AT().String()
		} else {
			site = doc.Site.URL().String()
		}
		if site != pubUrl {
			t.Errorf("invalid site: %s, wanted %s", doc.Site.URL().String(), pubUrl)
		}
		if doc.Title != title {
			t.Errorf("invalid title: %s, wanted %s", doc.Title, title)
		}
		if !doc.PublishedAt.Equal(publishedAt) {
			t.Errorf("invalid publishedAt: %s, wanted %s", doc.PublishedAt, publishedAt)
		}
		if *doc.CoverImage != *coverImage {
			t.Errorf("invalid cover image: %v, wanted %v", *doc.CoverImage, *coverImage)
		}
		if *doc.Path != path {
			t.Errorf("invalid path: %s, wanted %s", *doc.Path, path)
		}
		if doc.Content == nil {
			t.Errorf("invalid content: is nil")
		} else {
			if doc.Content.Collection().String() != `pub.leaflet.content` {
				t.Errorf("invalid content type: %s", doc.Content.Collection())
			}
			if !slices.Equal(doc.Content.Raw, []byte(`{"$type":"pub.leaflet.content","pages":[{"$type":"pub.leaflet.pages.linearDocument","blocks":[{"$type":"pub.leaflet.pages.linearDocument#block","block":{"$type":"pub.leaflet.blocks.text","plaintext":"hiiiiiiiii"}}],"id":"019d1297-2fdd-733b-9837-911e1758f300"}]}`)) {
				t.Errorf("invalid content raw: %s", string(doc.Content.Raw))
			}
		}
		if !slices.Equal(doc.Tags, tags) {
			t.Errorf("invalid tags: %v, wanted %v", doc.Tags, tags)
		}

		b, err = xrpc.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))

		c := new(content)
		ok := doc.Content.As(c)
		if !ok {
			t.Fatal("expected content type to be", c.Collection().String())
		}
		if c.Pages == nil {
			t.Errorf("invalid content pages: nil")
		}
		t.Logf("%v", c.Pages)
	})
}

var genDocAt = []string{
	"at://did:plc:zcanytzlaumjwgaopolw6wes/site.standard.document/3mhmdp3qobs2o", // leaflet doc
	"at://did:plc:revjuqmkvrw6fnkxppqtszpv/site.standard.document/3mbfqhezge25u", // pckt doc
}

func TestGetDocument(t *testing.T) {
	if testing.Short() {
		t.Skip("not doing http requests in short")
	}
	for _, uri := range genDocAt {
		client := getClient()
		u, err := atproto.ParseURI(context.Background(), client.Directory(), uri)
		if err != nil {
			t.Fatal(err)
		}
		union, err := client.FetchURI(context.Background(), u)
		if err != nil {
			t.Fatal(err)
		}
		doc := new(site.Document)
		if !union.Value.As(doc) {
			t.Fatalf("cannot convert union to document: %s", union.Value.Raw)
		}
		pubURL, err := doc.PublicationURL(context.Background(), client)
		if err != nil {
			t.Fatal(err)
		}
		valid, err := doc.Verify(
			context.Background(),
			client.HTTP(),
			pubURL,
			u.Authority(),
			*u.RecordKey(),
		)
		if err != nil {
			t.Errorf("cannot verify %s: %v", uri, err)
		} else if !valid {
			t.Errorf("cannot verify %s", uri)
		}
		if doc.BlueskyPostRef != nil {
			_, err := doc.BlueskyPostRef.GetRef(context.Background(), client)
			if err != nil {
				t.Errorf("cannot get bluesky post ref %s: %v", doc.BlueskyPostRef.URI, err)
			}
		}
	}
}

func TestListDocuments(t *testing.T) {
	if testing.Short() {
		t.Skip("not doing http requests in short")
	}
	for _, uri := range genDocAt {
		client := getClient()
		u, err := atproto.ParseURI(context.Background(), client.Directory(), uri)
		if err != nil {
			t.Fatal(err)
		}
		docs, _, err := xrpc.ListRecords[*site.Document](context.Background(), client, u.Authority(), 0, "", false)
		if err != nil {
			t.Fatal(err)
		}
		if docs == nil {
			t.Errorf("docs is nil")
		}
		for i, doc := range docs {
			if doc.Value == nil {
				t.Errorf("doc %d is nil", i)
			}
		}
	}
}
