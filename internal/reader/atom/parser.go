// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package atom // import "miniflux.app/v2/internal/reader/atom"

import (
	"encoding/xml"
	"fmt"
	"io"

	"miniflux.app/v2/internal/model"
	xml_decoder "miniflux.app/v2/internal/reader/xml"
)

type atomFeed interface {
	Transform(baseURL string) *model.Feed
}

// Parse returns a normalized feed struct from a Atom feed.
func Parse(baseURL string, r io.ReadSeeker) (*model.Feed, error) {
	var rawFeed atomFeed
	if getAtomFeedVersion(r) == "0.3" {
		rawFeed = new(atom03Feed)
	} else {
		rawFeed = new(atom10Feed)
	}
	r.Seek(0, io.SeekStart)

	if err := xml_decoder.NewXMLDecoder(r).Decode(rawFeed); err != nil {
		return nil, fmt.Errorf("atom: unable to parse feed: %w", err)
	}

	return rawFeed.Transform(baseURL), nil
}

func getAtomFeedVersion(data io.ReadSeeker) string {
	decoder := xml_decoder.NewXMLDecoder(data)
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		if element, ok := token.(xml.StartElement); ok {
			if element.Name.Local == "feed" {
				for _, attr := range element.Attr {
					if attr.Name.Local == "version" && attr.Value == "0.3" {
						return "0.3"
					}
				}
				return "1.0"
			}
		}
	}
	return "1.0"
}
