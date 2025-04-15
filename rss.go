package main

import (
	"encoding/xml"
	"net/http"
)

// RSSFeed represents the top-level RSS element with namespace attributes.
type RSSFeed struct {
	XMLName  xml.Name   `xml:"rss"`
	Version  string     `xml:"version,attr"`
	ITunesNS string     `xml:"xmlns:itunes,attr,omitempty"` // e.g., "http://www.itunes.com/dtds/podcast-1.0.dtd"
	MediaNS  string     `xml:"xmlns:media,attr,omitempty"`  // e.g., "http://search.yahoo.com/mrss/"
	Channel  RSSChannel `xml:"channel"`
}

// RSSChannel holds channel-level metadata and a slice of items.
type RSSChannel struct {
	Title          string      `xml:"title"`
	Description    string      `xml:"description"`
	Link           string      `xml:"link"`
	Language       string      `xml:"language,omitempty"`
	ITunesAuthor   string      `xml:"itunes:author,omitempty"`
	ITunesSummary  string      `xml:"itunes:summary,omitempty"`
	ITunesExplicit string      `xml:"itunes:explicit,omitempty"`
	ITunesImage    ITunesImage `xml:"itunes:image,omitempty"`
	Items          []RSSItem   `xml:"item"`
}

// ITunesImage is used for specifying the podcast artwork.
type ITunesImage struct {
	Href string `xml:"href,attr"`
}

// RSSItem represents a single podcast episode.
type RSSItem struct {
	Title          string        `xml:"title"`
	Description    string        `xml:"description"`
	Link           string        `xml:"link"`
	PubDate        string        `xml:"pubDate"`
	GUID           GUID          `xml:"guid"`
	Enclosure      Enclosure     `xml:"enclosure"`
	ITunesDuration string        `xml:"itunes:duration,omitempty"`
	ITunesExplicit string        `xml:"itunes:explicit,omitempty"`
	MediaContent   *MediaContent `xml:"media:content,omitempty"`
}

// GUID can include a flag to indicate if it is a permalink.
type GUID struct {
	Value       string `xml:",chardata"`
	IsPermaLink string `xml:"isPermaLink,attr,omitempty"`
}

// Enclosure holds information about the media file for the episode.
type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

// MediaContent represents additional media info, often from the media namespace.
type MediaContent struct {
	URL    string `xml:"url,attr"`
	Medium string `xml:"medium,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
}

// urlToFeed is a placeholder function where you would
// implement the fetching and unmarshaling of the RSS feed.
func urlToFeed(url string) (*RSSFeed, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var feed RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}
	return &feed, nil
}
