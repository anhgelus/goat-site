# GoAT Site

GoAT Site implements [Standard.site](https://standard.site/) in Go.

Use [`anhgelus.world/xrpc`](https://tangled.org/anhgelus.world/xrpc/), a lightweight XRPC client.

Main repository is hosted on [Tangled](https://tangled.org/anhgelus.world/goat-site/), an ATProto forge.

## Usage

> [!NOTE] Check [anhgelus.world/xrpc's documentation](https://tangled.org/anhgelus.world/xrpc) first!

Get the module with:
```bash
go get -u tangled.org/anhgelus.world/goat-site
```

Each [Standard.site lexicon](https://standard.site/#definitions) is implemented:
- `Publication` is for `site.standard.publication`;
- `Document` is for `site.standard.document`;
- `Subscription` is for `site.standard.graph.subscription`.

These types implement `xrpc.Record`, an interface describing records.

You can get, list, create, update or delete them with functions:
- `xrpc.GetRecord[*site.Publication]` to get a publication;
- `xrpc.ListRecords[*site.Document]` to list documents;
- `xrpc.CreateRecord[*site.Document]` to create a new document;
- `xrpc.UpdateRecord[*site.Subscription]` to update a subscription;
- `xrpc.DeleteRecord[*site.Publication]` to delete a publication.

You can [verify](https://standard.site/docs/verification/) a publication with `Publication.Verify` and a document with
`Document.Verify`:
```go
var pub *site.Publication
var did *atproto.DID
var client xrpc.Client
valid, err := pub.Verify(context.Background(), client, did, "pub_rkey")
if err != nil {
    panic(err)
}
if !valid {
    println("invalid publication :(")
}

var doc *site.Document
pubUrl, err := doc.PublicationURL(context.Background(), client)
if err != nil {
    panic(err)
}
valid, err = doc.Verify(context.Background(), client, pubUrl, did, "doc_rkey")
if err != nil {
    panic(err)
}
if !valid {
    panic("invalid document :(")
}
```

## Creating custom records

`Document.Content` is an open union: you can create your own lexicon to use it.

If the NSID of your lexicon is `tld.example.content` and its definition in Go is:
```go
type Content struct {
    // Pars represents the paragraphs in [Content].
    Pars []string `json:"pars"`
}
```

To use it, you have to implement `site.Record`:
```go
var CollectionContent = atproto.NewNSIDBuilder(`tld.example`).Name("content").Build()

func (c *Content) Collection() *atproto.NSID {
    return CollectionContent
}
```

But if you use `xrpc.GetRecord[*site.Document]` to retrieve one, it will return a simple `site.Document` without your 
custom content!
The `Document.Content` field is a `xrpc.Union`, a type representing an open union.
You can get the collection of the content with `Union.Collection()` and the raw bytes with `Union.Raw`.
You can also directly parse your `Content` with `Union.As`:
```go
var doc *site.Document
c := new(Content)
// returns an error if it cannot parse or if the type is invalid
if !doc.Content.As(c) {
    panic("not a Content :(")
}
```

### Marshal/Unmarshal

See [anhgelus.world/xrpc documentation](https://tangled.org/anhgelus.world/xrpc/#complexe-records).

## Extending lexicons

Lexicons defined by Standard.site [can be extended](https://standard.site/docs/introduction/#design-philosophy).

To extend a lexicon, you can create a new type and embed the base lexicon:
```go
type CustomPublication struct {
    site.Publication
    // your custom fields
}
```

You can call any functions with this new lexicon: the embedded base lexicon already implements the `xrpc.Record`
interface!
