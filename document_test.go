package site_test

import (
	"context"
	"encoding/json"
	"slices"
	"testing"

	"pgregory.net/rapid"
	site "tangled.org/anhgelus.world/goat-site"
)

type content struct {
	Pages any `json:"pages"`
}

func (c *content) Type() string {
	return `pub.leaflet.content`
}

func TestDocument_JSON(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		var pubUrl string
		if rapid.Bool().Draw(t, "url_at?") {
			pubUrl = "at://" + genDid(t, "url_did") +
				"/" + site.CollectionPublication +
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
			"publishedAt": publishedAt,
			"path":        path,
			"description": description,
			"coverImage":  coverImageRaw,
			"content":     json.RawMessage(`{"$type":"pub.leaflet.content","pages":[{"$type":"pub.leaflet.pages.linearDocument","blocks":[{"$type":"pub.leaflet.pages.linearDocument#block","block":{"$type":"pub.leaflet.blocks.text","plaintext":"hiiiiiiiii"}}],"id":"019d1297-2fdd-733b-9837-911e1758f300"}]}`),
			"textContent": textContent,
			"bskyPostRef": json.RawMessage(`{"cid":"bafyreidepvhssy3zglq3bo4nauszqhqmbk6lzzfay3r2nskvijyiewlr2u","commit":{"cid":"bafyreickwfv4p2jr6zvbdk6mldmddag2m6grpkbbvkvz57mvaqso5dpf5e","rev":"3mhm4oeyyzi2g"},"uri":"at://did:plc:jdhpqeb4cb4mng533dx56cbc/app.bsky.feed.post/3mhm4oevhmk2d","validationStatus":"valid"}`),
			"tags":        tags,
			"updatedAt":   updatedAt,
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
		doc := v.Record.(*site.Document)
		if doc.Site.String() != pubUrl {
			t.Errorf("invalid site: %s, wanted %s", doc.Site, pubUrl)
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
		if doc.Content.Record != nil {
			t.Errorf("invalid content lexicon: %v", doc.Content.Record)
		} else {
			if doc.Content.Type != `pub.leaflet.content` {
				t.Errorf("invalid content type: %s", doc.Content.Type)
			}
			if !slices.Equal(doc.Content.Raw, []byte(`{"$type":"pub.leaflet.content","pages":[{"$type":"pub.leaflet.pages.linearDocument","blocks":[{"$type":"pub.leaflet.pages.linearDocument#block","block":{"$type":"pub.leaflet.blocks.text","plaintext":"hiiiiiiiii"}}],"id":"019d1297-2fdd-733b-9837-911e1758f300"}]}`)) {
				t.Errorf("invalid content raw: %s", doc.Content.Raw)
			}
		}
		if !slices.Equal(doc.Tags, tags) {
			t.Errorf("invalid tags: %v, wanted %v", doc.Tags, tags)
		}

		b, err = json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))

		c := new(content)
		err = doc.Content.As(c)
		if err != nil {
			t.Fatal(err)
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
		uri, client := getClient(t, uri)
		doc, err := site.GetRecord[*site.Document](context.Background(), client, uri.Authority(), uri.RecordKey())
		if err != nil {
			t.Fatal(err)
		}
		if doc == nil {
			t.Errorf("doc is nil")
		}
	}

}

func TestListDocuments(t *testing.T) {
	if testing.Short() {
		t.Skip("not doing http requests in short")
	}
	for _, uri := range genDocAt {
		uri, client := getClient(t, uri)
		docs, _, err := site.ListRecords[*site.Document](context.Background(), client, uri.Authority(), "", false)
		if err != nil {
			t.Fatal(err)
		}
		if docs == nil {
			t.Errorf("docs is nil")
		}
		for i, doc := range docs {
			if doc == nil {
				t.Errorf("doc %d is nil", i)
			}
		}
	}
}

func TestDocumentVerification(t *testing.T) {
	tag := site.GetDocumentVerificationTag("did:plc:xyz789", "rkey")
	if tag != `<link rel="site.standard.document" href="at://did:plc:xyz789/site.standard.document/rkey">` {
		t.Errorf("invalid tag: %s", tag)
	}
}

func TestDocument_Verify(t *testing.T) {
	if testing.Short() {
		t.Skip("not doing http requests in short")
	}
	for _, uri := range genDocAt {
		uri, client := getClient(t, uri)
		doc, err := site.GetRecord[*site.Document](context.Background(), client, uri.Authority(), uri.RecordKey())
		if err != nil {
			t.Fatal(err)
		}
		pubURL, err := doc.PublicationURL(context.Background(), client)
		if err != nil {
			t.Fatal(err)
		}
		valid, err := doc.Verify(
			context.Background(),
			client.Client,
			pubURL,
			uri.Authority(),
			uri.RecordKey(),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !valid {
			t.Errorf("cannot verify %s", uri)
		}
	}
}
