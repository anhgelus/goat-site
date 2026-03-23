package site

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bluesky-social/indigo/api/agnostic"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/atproto/syntax"
	lexutil "github.com/bluesky-social/indigo/lex/util"
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

var (
	ErrRecordAlreadyParsed = errors.New("record already parsed")
	ErrNoContent           = errors.New("no content")
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

// MaxItemsPerList is the number of items per list call.
const MaxItemsPerList = 25

// Result is returned when after creating a record.
type Result struct {
	URI              string
	CID              string
	ValidationStatus *string
	Commit           *agnostic.RepoDefs_CommitMeta
}

// get returns the T in the repo associated with the rkey.
// Automatically uses the latest CID.
func get[T Record](ctx context.Context, client lexutil.LexClient, collection string, repo syntax.AtIdentifier, rkey syntax.RecordKey) (t T, err error) {
	var rec *agnostic.RepoGetRecord_Output
	rec, err = agnostic.RepoGetRecord(ctx, client, "", collection, repo.String(), rkey.String())
	if err != nil {
		return
	}
	var v *RecordJSON
	err = json.Unmarshal(*rec.Value, &v)
	if err != nil {
		return
	}
	if v.Record == nil {
		err = ErrInvalidType{collection, v.Type}
		return
	}
	if v.Record.Type() != collection {
		err = ErrInvalidType{collection, v.Record.Type()}
		return
	}
	return v.Record.(T), nil
}

// listRecord returns all the Ts stored in the repo and the cursor.
//
// See [MaxItemsPerList].
func listRecord[T Record](ctx context.Context, client lexutil.LexClient, collection string, repo syntax.AtIdentifier, cursor string, reverse bool) ([]T, *string, error) {
	rec, err := agnostic.RepoListRecords(ctx, client, collection, cursor, MaxItemsPerList, repo.String(), reverse)
	if err != nil {
		return nil, nil, err
	}
	docs := make([]T, MaxItemsPerList)
	i := 0
	for i = range len(rec.Records) {
		r := rec.Records[i]
		var v *RecordJSON
		err = json.Unmarshal(*r.Value, &v)
		if err != nil {
			return nil, nil, err
		}
		if v.Record == nil {
			return nil, nil, ErrInvalidType{collection, v.Type}
		}
		if v.Record.Type() != collection {
			return nil, nil, ErrInvalidType{collection, v.Record.Type()}
		}
		docs[i] = v.Record.(T)
	}
	return docs[:i], rec.Cursor, nil
}

// createRecord a T in a repo with the given rkey.
// Always tries to validate the [Document] against the lexicon saved.
//
// Rkey can be nil.
func createRecord[T Record](ctx context.Context, client lexutil.LexClient, collection string, repo syntax.AtIdentifier, rkey *syntax.RecordKey, v T) (*Result, error) {
	mp, err := MarshalToMap(&RecordJSON{Record: v})
	if err != nil {
		return nil, err
	}
	var cv *string
	if rkey != nil {
		t := rkey.String()
		cv = &t
	}
	t := true
	out, err := agnostic.RepoCreateRecord(ctx, client, &agnostic.RepoCreateRecord_Input{
		Collection: collection,
		Record:     mp,
		Repo:       repo.String(),
		Rkey:       cv,
		Validate:   &t,
	})
	if err != nil {
		return nil, err
	}
	return &Result{out.Uri, out.Cid, out.ValidationStatus, out.Commit}, nil
}

// updateRecord T in a repo with the given rkey.
// Always tries to validate the [Document] against the lexicon saved.
func updateRecord[T Record](ctx context.Context, client lexutil.LexClient, collection string, repo syntax.AtIdentifier, rkey syntax.RecordKey, v T) (*Result, error) {
	mp, err := MarshalToMap(&RecordJSON{Record: v})
	if err != nil {
		return nil, err
	}
	t := true
	out, err := agnostic.RepoPutRecord(ctx, client, &agnostic.RepoPutRecord_Input{
		Collection: collection,
		Record:     mp,
		Repo:       repo.String(),
		Rkey:       rkey.String(),
		Validate:   &t,
		//SwapRecord: &cid,
	})
	if err != nil {
		return nil, err
	}
	return &Result{out.Uri, out.Cid, out.ValidationStatus, out.Commit}, nil
}

// delete in a repo with the given rkey.
func deleteRecord(ctx context.Context, client lexutil.LexClient, collection string, repo syntax.AtIdentifier, rkey syntax.RecordKey) error {
	_, err := atproto.RepoDeleteRecord(ctx, client, &atproto.RepoDeleteRecord_Input{
		Collection: collection,
		Repo:       repo.String(),
		Rkey:       rkey.String(),
	})
	return err
}

// createAtURI returns a valid [syntax.ATURI].
func createAtURI(repo syntax.AtIdentifier, collection string, rkey syntax.RecordKey) string {
	return fmt.Sprintf("at://%s/%s/%s", repo, collection, rkey)
}
