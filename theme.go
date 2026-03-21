package site

import (
	"image/color"
	"math"
)

const (
	LexiconTheme          = LexiconBase + ".theme"
	LexiconThemeColor     = LexiconTheme + ".color"
	LexiconThemeColorRGB  = LexiconThemeColor + "#rgb"
	LexiconThemeColorRGBA = LexiconThemeColor + "#rgba"
)

// Theme ensures [Publication]s maintain their visual identity across different reading applications and platforms by
// defining core colors for content display.
type Theme struct {
	// Background is the color used for content background.
	Background any `json:"background"`
	// Foreground is the color used for content text.
	Foreground any `json:"foreground"`
	// Accent is the color used for links and button backgrounds.
	Accent any `json:"accent"`
	// AccentForeground is the color used for button text.
	AccentForeground any `json:"accentForeground"`
}

func (t *Theme) Type() string {
	return LexiconTheme
}

// RGB represents a RGB color.
type RGB struct {
	Red   uint8 `json:"r"`
	Green uint8 `json:"g"`
	Blue  uint8 `json:"b"`
}

func (r *RGB) Type() string {
	return LexiconThemeColorRGB
}

func (r *RGB) RGBA() *RGBA {
	return &RGBA{*r, 100}
}

// RGBA represents a [color.RGBA].
type RGBA struct {
	RGB
	// Alpha is the alpha channel where 0 is transparent and 100 is opaque.
	Alpha uint8 `json:"a"`
}

func (r *RGBA) Type() string {
	return LexiconThemeColorRGBA
}

func NewRGBA(r *color.RGBA) *RGBA {
	red, g, b, a := r.RGBA()
	return &RGBA{
		RGB:   RGB{uint8(red), uint8(g), uint8(b)},
		Alpha: uint8(math.Floor(float64(a)/float64(255)) * 100),
	}
}
