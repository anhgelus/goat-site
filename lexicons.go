package site

import (
	"encoding/json"
	"errors"
)

type Lexicon interface {
	Type() string
}

const LexiconBase = "site.standard"

// LexiconJSON is used to convert [Lexicon] into JSON.
type LexiconJSON struct {
	Lexicon
}

func (l LexiconJSON) MarshalJSON() ([]byte, error) {
	v := struct {
		any
		Type string `json:"$type"`
	}{nil, l.Type()}
	switch l.Type() {
	case LexiconPublication:
		v.any = l.Lexicon.(*Publication)
	case LexiconDocument:
		v.any = l.Lexicon.(*Document)
	case LexiconSubscription:
		v.any = l.Lexicon.(*Subscription)
	case LexiconTheme:
		v.any = l.Lexicon.(*Theme)
	case LexiconThemeColorRGB:
		v.any = l.Lexicon.(*RGB)
	case LexiconThemeColorRGBA:
		v.any = l.Lexicon.(*RGBA)
	default:
		return nil, errors.New("unsupported lexicon type")
	}
	return json.Marshal(v)
}

func (l LexiconJSON) UnmarshalJSON(b []byte) error {
	var v struct {
		Type string `json:"$type"`
	}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	switch v.Type {
	case LexiconPublication:
		l.Lexicon = l.Lexicon.(*Publication)
	case LexiconDocument:
		l.Lexicon = l.Lexicon.(*Document)
	case LexiconSubscription:
		l.Lexicon = l.Lexicon.(*Subscription)
	case LexiconTheme:
		l.Lexicon = l.Lexicon.(*Theme)
	case LexiconThemeColorRGB:
		l.Lexicon = l.Lexicon.(*RGB)
	case LexiconThemeColorRGBA:
		l.Lexicon = l.Lexicon.(*RGBA)
	default:
		return errors.New("unsupported lexicon type")
	}
	return json.Unmarshal(b, l.Lexicon)
}
