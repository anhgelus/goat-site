package site_test

import (
	"encoding/json"
	"testing"

	site "tangled.org/anhgelus.world/goat-site"
)

const samplePub = `{
  "$type": "site.standard.publication",
  "basicTheme": {
    "$type": "site.standard.theme.basic",
    "accent": {
      "$type": "site.standard.theme.color#rgb",
      "b": 20,
      "g": 105,
      "r": 139
    },
    "accentForeground": {
      "$type": "site.standard.theme.color#rgb",
      "b": 204,
      "g": 243,
      "r": 255
    },
    "background": {
      "$type": "site.standard.theme.color#rgb",
      "b": 225,
      "g": 249,
      "r": 255
    },
    "foreground": {
      "$type": "site.standard.theme.color#rgb",
      "b": 32,
      "g": 53,
      "r": 74
    }
  },
  "description": "the latest and greatest from pckt !",
  "icon": {
    "$type": "blob",
	"ref": {
      "$link": "bafkreia3gaejwdadslicpqbgtzitcysop7lhuyry6bjf6xlf5fe7jvvcdy"
    },
    "mimeType": "image/png",
    "size": 8535
  },
  "name": "pckt - Dev Journal",
  "preferences": {
    "showInDiscover": true
  },
  "theme": {
    "$type": "blog.pckt.theme",
    "dark": {
      "accent": "#ffc947",
      "background": "#3d2a1a",
      "link": "#ffe082",
      "surfaceHover": "#4d3822",
      "text": "#fff9e0"
    },
    "font": "sans",
    "light": {
      "accent": "#8b6914",
      "background": "#fff9e1",
      "link": "#d68910",
      "surfaceHover": "#fff3cc",
      "text": "#4a3520"
    },
    "transparency": 100
  },
  "url": "https://devlog.pckt.blog"
}`

func TestPublication_JSON(t *testing.T) {
	var v *site.LexiconJSON
	err := json.Unmarshal([]byte(samplePub), &v)
	if err != nil {
		t.Fatal(err)
	}
	pub := v.Lexicon.(*site.Publication)
	if pub.Name != "pckt - Dev Journal" {
		t.Errorf("invalid name: %s", pub.Name)
	}
	if pub.URL != "https://devlog.pckt.blog" {
		t.Errorf("invalid url: %s", pub.URL)
	}
	if *pub.Description != "the latest and greatest from pckt !" {
		t.Errorf("invalid description: %s", *pub.Description)
	}
	if pub.Icon.CID != "bafkreia3gaejwdadslicpqbgtzitcysop7lhuyry6bjf6xlf5fe7jvvcdy" {
		t.Errorf("invalid Icon CID: %s", pub.Icon.CID)
	}
	if pub.Icon.MimeType != "image/png" {
		t.Errorf("invalid Icon MimeType: %s", pub.Icon.MimeType)
	}
	if pub.Icon.Size != 8535 {
		t.Errorf("invalid Icon Size: %d", pub.Icon.Size)
	}
	if !pub.Preferences.ShowInDiscover {
		t.Errorf("invalid Preferences ShowInDiscover: %v", pub.Preferences.ShowInDiscover)
	}
	theme := pub.BasicTheme
	if *theme.Accent != *site.NewRGB(139, 105, 20) {
		t.Errorf("invalid theme accent color: %s", theme.Accent)
	}
	if *theme.AccentForeground != *site.NewRGB(255, 243, 204) {
		t.Errorf("invalid theme accent foreground color: %s", theme.AccentForeground)
	}
	if *theme.Background != *site.NewRGB(255, 249, 225) {
		t.Errorf("invalid theme background color: %s", theme.Background)
	}
	if *theme.Foreground != *site.NewRGB(74, 53, 32) {
		t.Errorf("invalid theme foreground color: %s", theme.Foreground)
	}

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
