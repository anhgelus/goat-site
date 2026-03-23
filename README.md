# GoAT Site

GoAT Site implements [Standard.site](https://standard.site/) in Go.

Use official [`bluesky-social/indigo`](https://github.com/bluesky-social/indigo/).

Main repository is hosted on [Tangled](https://tangled.org/anhgelus.world/goat-site/), an ATProto forge.

## Usage

Get the module with:
```bash
go get -u tangled.org/anhgelus.world/goat-site
```

Each Standard.site lexicon is implemented:
- `Publication` is for `site.standard.publication`;
- `Document` is for `site.standard.document`;
- `Subscription` is for `site.standard.graph.subscription`.

These types implement `Record`, an interface describing records.

You can get, list, create, update or delete them with functions.
Each function starts with the action followed by the lexicon's name, e.g.,
- `GetPublication` to get a publication;
- `ListDocuments` to list documents;
- `CreateDocument` to create a new document.

Currently, functions related to `Subscription` are not implemented.

You can [verify](https://standard.site/docs/verification/) a publication with `Publication.Verify`.

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
const CollectionContent = `tld.example.content`

func (c *Content) Type() string {
    return CollectionContent
}
```

But, if you use `site.GetDocument` to retrieve one, it will return a simple `site.Document` without your custom content!
The `Document.Content` field is a `site.RecordJSON`, a wrapper.
You can get the type of the content with `RecordJSON.Type` and the raw bytes with `RecordJSON.Raw`.
You can also directly parse your `Content` with `RecordJSON.As`:
```go
var doc *site.Document
var c *Content
// returns an error if it cannot parse or if the type is invalid
err := doc.Content.As(c)
if err != nil {
    panic(err)
}
```

### Marshal/Unmarshal

When your record is sent, it is firstly marshaled to a map.
We provide `site.MarshalToMap` which works like the JSON API:
```go
var c *Content
// mp is the map[string]any created
mp, err := site.MarshalToMap(c)
if err != nil {
    panic(err)
}
/*
mp = map[string]any{"content": []string{}}
*/
```

It uses the `json` tag to determine how to marshal the content.
It supports `omitempty`, `string` and embedded type.

If you are using complexe types, you may have to implement `json.Unmarshaler` to unmarshal from JSON and
`site.MarshalerMap` to marshal to a map.
```go
func (c *Content) MarshalMap() (map[string]any, error) {
    mp := make(map[string]any, 1)
    mp["foo"] = "bar"
    return mp, nil
}
// the future call to site.MarshalToMap on *Content will return map[string]any{"foo":"bar"}.
```
