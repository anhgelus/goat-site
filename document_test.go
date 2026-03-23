package site_test

import (
	"context"
	"encoding/json"
	"slices"
	"testing"
	"time"

	"github.com/bluesky-social/indigo/atproto/syntax"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	site "tangled.org/anhgelus.world/goat-site"
)

const sampleDoc = `
{"$type":"site.standard.document","bskyPostRef":{"cid":"bafyreidepvhssy3zglq3bo4nauszqhqmbk6lzzfay3r2nskvijyiewlr2u","commit":{"cid":"bafyreickwfv4p2jr6zvbdk6mldmddag2m6grpkbbvkvz57mvaqso5dpf5e","rev":"3mhm4oeyyzi2g"},"uri":"at://did:plc:jdhpqeb4cb4mng533dx56cbc/app.bsky.feed.post/3mhm4oevhmk2d","validationStatus":"valid"},"content":{"$type":"pub.leaflet.content","pages":[{"$type":"pub.leaflet.pages.linearDocument","blocks":[{"$type":"pub.leaflet.pages.linearDocument#block","block":{"$type":"pub.leaflet.blocks.text","plaintext":"hiiiiiiiii"}}],"id":"019d1297-2fdd-733b-9837-911e1758f300"}]},"path":"/3mhm4obhnx22y","publishedAt":"2026-03-21T22:52:35.182Z","site":"at://did:plc:jdhpqeb4cb4mng533dx56cbc/site.standard.publication/3mhm4m2tets2y","tags":[],"title":"hello world"}
`

func TestDocument_JSON(t *testing.T) {
	var v *site.RecordJSON
	err := json.Unmarshal([]byte(sampleDoc), &v)
	if err != nil {
		t.Fatal(err)
	}
	doc := v.Record.(*site.Document)
	if doc.Site != `at://did:plc:jdhpqeb4cb4mng533dx56cbc/site.standard.publication/3mhm4m2tets2y` {
		t.Errorf("invalid site: %s", doc.Site)
	}
	if doc.Title != `hello world` {
		t.Errorf("invalid title: %s", doc.Title)
	}
	tt, _ := time.Parse(site.TimeFormat, "2026-03-21T22:52:35.182Z")
	if !doc.PublishedAt.Equal(tt) {
		t.Errorf("invalid publishedAt: %s", doc.PublishedAt.Format(site.TimeFormat))
	}
	if *doc.Path != `/3mhm4obhnx22y` {
		t.Errorf("invalid path: %s", *doc.Path)
	}

	if doc.Content.Record != nil {
		t.Errorf("invalid content lexicon: %v", doc.Content.Record)
	}
	if doc.Content.Type != `pub.leaflet.content` {
		t.Errorf("invalid content type: %s", doc.Content.Type)
	}
	if !slices.Equal(doc.Content.Raw, []byte(`{"$type":"pub.leaflet.content","pages":[{"$type":"pub.leaflet.pages.linearDocument","blocks":[{"$type":"pub.leaflet.pages.linearDocument#block","block":{"$type":"pub.leaflet.blocks.text","plaintext":"hiiiiiiiii"}}],"id":"019d1297-2fdd-733b-9837-911e1758f300"}]}`)) {
		t.Errorf("invalid content raw: %s", doc.Content.Raw)
	}
	if len(doc.Tags) > 0 {
		t.Errorf("invalid tags: %v", doc.Tags)
	}

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

const testDoc = "at://did:plc:zcanytzlaumjwgaopolw6wes/site.standard.document/3mhmdp3qobs2o"

var (
	docURI    syntax.ATURI
	docClient *lexutil.LexClient
)

func TestGetDocument(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	uri, client := getClient(t, testDoc, &docURI, &docClient)
	doc, err := site.GetDocument(context.Background(), client, uri.Authority(), uri.RecordKey())
	if err != nil {
		t.Fatal(err)
	}
	if doc == nil {
		t.Errorf("doc is nil")
	}
}

func TestListDocuments(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	uri, client := getClient(t, testDoc, &docURI, &docClient)
	docs, _, err := site.ListDocuments(context.Background(), client, uri.Authority(), "", false)
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

func TestDocumentVerification(t *testing.T) {
	tag := site.GetDocumentVerificationTag("did:plc:xyz789", "rkey")
	if tag != `<link rel="site.standard.document" href="at://did:plc:xyz789/site.standard.document/rkey">` {
		t.Errorf("invalid tag: %s", tag)
	}
}
