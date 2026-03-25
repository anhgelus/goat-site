package site

import (
	"context"
	"encoding/json"

	"github.com/bluesky-social/indigo/api/agnostic"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/atproto/syntax"
	lexutil "github.com/bluesky-social/indigo/lex/util"
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

// GetRecord returns the [Record] in the repo associated with the rkey.
// Automatically uses the latest CID.
//
// Returns [ErrInvalidType] if the [Record] got doesn't have a valid type.
func GetRecord[T Record](ctx context.Context, client lexutil.LexClient, repo syntax.AtIdentifier, rkey syntax.RecordKey) (t T, err error) {
	var rec *agnostic.RepoGetRecord_Output
	rec, err = agnostic.RepoGetRecord(ctx, client, "", t.Type(), repo.String(), rkey.String())
	if err != nil {
		return
	}
	var v *RecordJSON
	err = json.Unmarshal(*rec.Value, &v)
	if err != nil {
		return
	}
	if v.GetType() != t.Type() {
		err = ErrInvalidType{t.Type(), v.Type}
		return
	}
	return v.Record.(T), nil
}

// ListRecords returns all the [Record]s stored in the repo and the cursor.
//
// Returns [ErrInvalidType] if a [Record] got doesn't have a valid type.
//
// See [MaxItemsPerList].
func ListRecords[T Record](ctx context.Context, client lexutil.LexClient, repo syntax.AtIdentifier, cursor string, reverse bool) ([]T, *string, error) {
	var t T
	rec, err := agnostic.RepoListRecords(ctx, client, t.Type(), cursor, MaxItemsPerList, repo.String(), reverse)
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
		if v.GetType() != t.Type() {
			return nil, nil, ErrInvalidType{t.Type(), v.Type}
		}
		docs[i] = v.Record.(T)
	}
	return docs[:i], rec.Cursor, nil
}

// CreateRecord in a repo with the given rkey.
// Always tries to validate the [Record] against the lexicon saved.
//
// Rkey can be nil.
func CreateRecord[T Record](ctx context.Context, client lexutil.LexClient, repo syntax.AtIdentifier, rkey *syntax.RecordKey, v T) (*Result, error) {
	mp, err := MarshalToMap(AsJSON(v))
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
		Collection: v.Type(),
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

// UpdateRecord in a repo with the given rkey.
// Always tries to validate the [Record] against the lexicon saved.
func UpdateRecord[T Record](ctx context.Context, client lexutil.LexClient, repo syntax.AtIdentifier, rkey syntax.RecordKey, v T) (*Result, error) {
	mp, err := MarshalToMap(AsJSON(v))
	if err != nil {
		return nil, err
	}
	t := true
	out, err := agnostic.RepoPutRecord(ctx, client, &agnostic.RepoPutRecord_Input{
		Collection: v.Type(),
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

// DeleteRecord in a repo with the given rkey.
func DeleteRecord[T Record](ctx context.Context, client lexutil.LexClient, repo syntax.AtIdentifier, rkey syntax.RecordKey) error {
	var t T
	_, err := atproto.RepoDeleteRecord(ctx, client, &atproto.RepoDeleteRecord_Input{
		Collection: t.Type(),
		Repo:       repo.String(),
		Rkey:       rkey.String(),
	})
	return err
}

// createAtURL returns a valid [syntax.ATURI].
func createAtURL(repo syntax.AtIdentifier, collection string, rkey syntax.RecordKey) *ATURL {
	return &ATURL{syntax.ATURI("at://" + repo.String() + "/" + collection + "/" + rkey.String())}
}
