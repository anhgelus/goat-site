package site

const LexiconSubscription = LexiconBase + ".graph.subscription"

// Subscription enable users to follow publications and receive updates about new content.
// They represent the social connection between readers and the publications they're interested in.
type Subscription struct {
	// Publication is an AT-URI reference to the publication record being subscribed to.
	// E.g., `at://did:plc:abc123/site.standard.publication/xyz789`.
	Publication string `json:"publication"`
}

func (s *Subscription) Type() string {
	return LexiconSubscription
}
