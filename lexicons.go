package site

import (
	"encoding/json"
	"errors"

	"github.com/bluesky-social/indigo/api/agnostic"
)

// Record represents an ATProto record.
type Record interface {
	Type() string
}

const (
	// CollectionBase is the base NSID for Standard.site.
	CollectionBase = "site.standard"
	CollectionBlob = "blob"

	TimeFormat = "2006-01-02T15:04:05.000Z"
)

// RecordJSON is used to encode and to decode [Record] from JSON.
type RecordJSON struct {
	// Record parsed.
	// Nil if [Record] is unknown.
	Record Record
	// Type stored if [Record] is unknown.
	// Set after [json.Unmarshal].
	Type string
	// Raw returns bytes stored if [Record] is unknown.
	// Set after [json.Unmarshal].
	Raw []byte
}

func (l *RecordJSON) MarshalJSON() ([]byte, error) {
	if l.Record == nil {
		return l.Raw, nil
	}
	mp, err := l.MarshalMap()
	if err != nil {
		return nil, err
	}
	mp["$type"] = l.Record.Type()
	return json.Marshal(mp)
}

func (l *RecordJSON) MarshalMap() (mp map[string]any, err error) {
	if l.Record == nil {
		err = json.Unmarshal(l.Raw, &mp)
		return
	}
	mp, err = MarshalToMap(l.Record)
	return
}

func (l *RecordJSON) UnmarshalJSON(b []byte) error {
	var v struct {
		Type string `json:"$type"`
	}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	switch v.Type {
	case CollectionPublication:
		l.Record = &Publication{}
	case CollectionDocument:
		l.Record = &Document{}
	case CollectionSubscription:
		l.Record = &Subscription{}
	case CollectionThemeBasic:
		l.Record = &Theme{}
	case CollectionThemeColorRGB:
		l.Record = &RGB{}
	case CollectionThemeColorRGBA:
		l.Record = &RGBA{}
	case CollectionBlob:
		l.Record = &Blob{}
	default:
		l.Raw = b
		l.Type = v.Type
		return nil
	}
	return json.Unmarshal(b, l.Record)
}

// Blob represents an ATProto `blob` type.
type Blob struct {
	CID      string `json:"-"`
	MimeType string `json:"mimeType"`
	Size     uint   `json:"size"`
}

func (b *Blob) Type() string {
	return CollectionBlob
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

var (
	ErrInvalidType = errors.New("invalid collection type")
)

// MaxItemsPerList is the number of items per list call.
const MaxItemsPerList = 25

// Result is returned when after creating a record.
type Result struct {
	URI              string
	CID              string
	ValidationStatus *string
	Commit           *agnostic.RepoDefs_CommitMeta
}
