package site

import (
	"time"

	"anhgelus.world/xrpc/atproto"
)

var (
	collectionGraph        = CollectionBase.SubAuthority("graph")
	CollectionSubscription = collectionGraph.Name("subscription").Build()
	CollectionRecommend    = collectionGraph.Name("recommend").Build()
)

// Subscription enable users to follow publications and receive updates about new content.
// They represent the social connection between readers and the publications they're interested in.
type Subscription struct {
	// Publication is an AT-URI reference to the publication record being subscribed to.
	// E.g., `at://did:plc:abc123/site.standard.publication/xyz789`.
	Publication atproto.RawURI `json:"publication"`
}

func (s *Subscription) Collection() *atproto.NSID {
	return CollectionSubscription
}

type Recommend struct {
	// Document is an AT-URI reference to the recommended document.
	Document  atproto.RawURI `json:"document"`
	CreatedAt time.Time      `json:"createdAt"`
}

func (r *Recommend) Collection() *atproto.NSID {
	return CollectionRecommend
}
