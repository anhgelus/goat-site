package site

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Record represents an ATProto record.
type Record interface {
	Type() string
}

const (
	// CollectionBase is the base NSID for Standard.site.
	CollectionBase = "site.standard"
	CollectionBlob = "blob"

	// TimeFormat is the standard time format specified by the ATProto.
	//
	// See [ParseTime]
	TimeFormat = "2006-01-02T15:04:05.000Z07:00"
)

// ParseTime returns a [time.Time] if it follows the standard time format specified by the ATProto.
//
// See [TimeFormat].
// Fallback to [time.RFC3339] if it doesn't work.
func ParseTime(raw string) (t time.Time, err error) {
	t, err = time.Parse(TimeFormat, raw)
	if err != nil {
		t, err = time.Parse(time.RFC3339, raw)
	}
	return
}

type ErrInvalidType struct {
	expected, got string
}

func (err ErrInvalidType) Error() string {
	return fmt.Sprintf("invalid collection type: expected %s, got %s", err.expected, err.got)
}

func (err ErrInvalidType) As(target any) bool {
	it, ok := target.(*ErrInvalidType)
	if !ok {
		return false
	}
	*it = ErrInvalidType{err.expected, err.got}
	return true
}

func (err ErrInvalidType) Is(e error) bool {
	var it ErrInvalidType
	ok := errors.As(e, &it)
	if !ok {
		return false
	}
	return it.expected == err.expected && it.got == err.got
}

var (
	ErrRecordAlreadyParsed = errors.New("record already parsed")
	ErrNoContent           = errors.New("no content")
)

// RecordJSON is used to encode and to decode [Record] from JSON.
//
// If the [Record] is known by the library, it is decoded in [RecordJSON.Record].
// Else its type is filled in [RecordJSON.Type] and its raw bytes are placed in [RecordJSON.Raw].
//
// See [AsJSON] to create a [RecordJSON] from a [Record].
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

// AsJSON wraps a [Record] as a [RecordJSON].
func AsJSON(r Record) *RecordJSON {
	return &RecordJSON{Record: r}
}

// GetType returns the type associated with the [RecordJSON].
func (r *RecordJSON) GetType() string {
	if r.Record != nil {
		return r.Record.Type()
	}
	return r.Type
}

// As unmarshals the [RecordJSON] as the provided [Record].
//
// [ErrRecordAlreadyParsed] if the [Record] was already parsed (stored in [RecordJSON.Record]).
// [ErrNoContent] if [RecordJSON.Raw] is nil.
func (r *RecordJSON) As(rec Record) error {
	if r.Record != nil {
		return ErrRecordAlreadyParsed
	}
	if r.Raw == nil {
		return ErrNoContent
	}
	if r.Type != rec.Type() {
		return ErrInvalidType{r.Type, rec.Type()}
	}
	return json.Unmarshal(r.Raw, rec)
}

func (r *RecordJSON) MarshalJSON() ([]byte, error) {
	if r.Record == nil {
		return r.Raw, nil
	}
	mp, err := r.MarshalMap()
	if err != nil {
		return nil, err
	}
	mp["$type"] = r.Record.Type()
	return json.Marshal(mp)
}

func (r *RecordJSON) MarshalMap() (mp map[string]any, err error) {
	if r.Record == nil {
		err = json.Unmarshal(r.Raw, &mp)
		return
	}
	mp, err = MarshalToMap(r.Record)
	return
}

func (r *RecordJSON) UnmarshalJSON(b []byte) error {
	var v struct {
		Type string `json:"$type"`
	}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	switch v.Type {
	case CollectionPublication:
		r.Record = &Publication{}
	case CollectionDocument:
		r.Record = &Document{}
	case CollectionSubscription:
		r.Record = &Subscription{}
	case CollectionThemeBasic:
		r.Record = &Theme{}
	case CollectionThemeColorRGB:
		r.Record = &RGB{}
	case CollectionThemeColorRGBA:
		r.Record = &RGBA{}
	case CollectionBlob:
		r.Record = &Blob{}
	default:
		r.Raw = b
		r.Type = v.Type
		return nil
	}
	return json.Unmarshal(b, r.Record)
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
