package site

import (
	"encoding/json"
)

type Lexicon interface {
	Type() string
}

const (
	LexiconBase = "site.standard"
	LexiconBlob = "blob"

	TimeFormat = "2006-01-02T15:04:05.000Z"
)

// LexiconJSON is used to encode and decode [Lexicon] from JSON.
type LexiconJSON struct {
	// Lexicon parsed.
	// Nil if [Lexicon] is unknown.
	Lexicon Lexicon
	// Type stored if [Lexicon] is unknown.
	// Set after [json.Unmarshal].
	Type string
	// Raw returns bytes stored if [Lexicon] is unknown.
	// Set after [json.Unmarshal].
	Raw []byte
}

func (l *LexiconJSON) MarshalJSON() ([]byte, error) {
	if l.Lexicon == nil {
		return l.Raw, nil
	}
	mp, err := l.MarshalMap()
	if err != nil {
		return nil, err
	}
	mp["$type"] = l.Lexicon.Type()
	return json.Marshal(mp)
}

func (l *LexiconJSON) MarshalMap() (mp map[string]any, err error) {
	if l.Lexicon == nil {
		err = json.Unmarshal(l.Raw, &mp)
		return
	}
	mp, err = MarshalToMap(l.Lexicon)
	return
}

func (l *LexiconJSON) UnmarshalJSON(b []byte) error {
	var v struct {
		Type string `json:"$type"`
	}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	switch v.Type {
	case LexiconPublication:
		l.Lexicon = &Publication{}
	case LexiconDocument:
		l.Lexicon = &Document{}
	case LexiconSubscription:
		l.Lexicon = &Subscription{}
	case LexiconThemeBasic:
		l.Lexicon = &Theme{}
	case LexiconThemeColorRGB:
		l.Lexicon = &RGB{}
	case LexiconThemeColorRGBA:
		l.Lexicon = &RGBA{}
	case LexiconBlob:
		l.Lexicon = &Blob{}
	default:
		l.Raw = b
		l.Type = v.Type
		return nil
	}
	return json.Unmarshal(b, l.Lexicon)
}

// Blob represents an ATProto `blob` type.
type Blob struct {
	CID      string `json:"-"`
	MimeType string `json:"mimeType"`
	Size     uint   `json:"size"`
}

func (b *Blob) Type() string {
	return LexiconBlob
}

func (b *Blob) MarshalMap() (map[string]any, error) {
	mp := make(map[string]any, 3)
	mp["mimeType"] = b.MimeType
	mp["size"] = b.Size
	mp["ref"] = map[string]any{"$link": b.CID}
	return mp, nil
}

func (b *Blob) UnmarshalJSON(data []byte) error {
	type t Blob
	var v struct {
		t
		Ref struct {
			Link string `json:"$link"`
		} `json:"ref"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*b = Blob(v.t)
	b.CID = v.Ref.Link
	return nil
}
