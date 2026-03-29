package site

import (
	"fmt"

	"tangled.org/anhgelus.world/xrpc/atproto"
)

var (
	// CollectionBase is the base NSID for Standard.site.
	CollectionBase = atproto.NewNSIDBuilder("site.standard")
	CollectionBlob = "blob"
)

type ErrInvalidCollection struct {
	expected, got *atproto.NSID
}

func (err ErrInvalidCollection) Error() string {
	return fmt.Sprintf("invalid collection: expected %s, got %s", err.expected, err.got)
}
