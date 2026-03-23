package site

import (
	"encoding/json"

	"github.com/bluesky-social/indigo/atproto/syntax"
)

const CollectionSubscription = CollectionBase + ".graph.subscription"

// Subscription enable users to follow publications and receive updates about new content.
// They represent the social connection between readers and the publications they're interested in.
type Subscription struct {
	// Publication is an AT-URI reference to the publication record being subscribed to.
	// E.g., `at://did:plc:abc123/site.standard.publication/xyz789`.
	Publication syntax.ATURI `json:"publication,string"`
}

func (s *Subscription) UnmarshalJSON(b []byte) error {
	var v struct {
		Publication string `json:"publication"`
	}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	s.Publication, err = syntax.ParseATURI(v.Publication)
	return err
}

func (s *Subscription) Type() string {
	return CollectionSubscription
}
