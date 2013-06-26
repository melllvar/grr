/*****************************************************************************
 **
 ** grr >:(
 ** https://github.com/melllvar/grr
 ** Copyright (C) 2013 Akop Karapetyan
 **
 ** This program is free software; you can redistribute it and/or modify
 ** it under the terms of the GNU General Public License as published by
 ** the Free Software Foundation; either version 2 of the License, or
 ** (at your option) any later version.
 **
 ** This program is distributed in the hope that it will be useful,
 ** but WITHOUT ANY WARRANTY; without even the implied warranty of
 ** MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 ** GNU General Public License for more details.
 **
 ** You should have received a copy of the GNU General Public License
 ** along with this program; if not, write to the Free Software
 ** Foundation, Inc., 675 Mass Ave, Cambridge, MA 02139, USA.
 **
 ******************************************************************************
 */
 
package parser

import (
  "time"
  "encoding/xml"
  "errors"
)

type RSS2Feed struct {
  XMLName xml.Name `xml:"rss"`
  Title string `xml:"channel>title"`
  Description string `xml:"channel>description"`
  Updated string `xml:"channel>lastBuildDate"`
  Link string `xml:"channel>link"`
  Entry []*RSS2Entry `xml:"channel>item"`
}

type RSS2Entry struct {
  Id string `xml:"guid"`
  Published string `xml:"pubDate"`
  EntryTitle string `xml:"title"`
  Link string `xml:"link"`
  Author string `xml:"creator"`
  EncodedContent string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
  Content string `xml:"description"`
}

func (rss2Feed *RSS2Feed) Marshal() (feed Feed, err error) {
  updated := time.Time {}
  if rss2Feed.Updated != "" {
    updated, err = AsTime(rss2Feed.Updated)
  }

  feed = Feed {
    Title: rss2Feed.Title,
    Description: rss2Feed.Description,
    Updated: updated,
    WWWURL: rss2Feed.Link,
  }

  if rss2Feed.Entry != nil {
    feed.Entry = make([]*Entry, len(rss2Feed.Entry))
    for i, v := range rss2Feed.Entry {
      var entryError error
      feed.Entry[i], entryError = v.Marshal()

      if entryError != nil && err == nil {
        err = entryError
      }
    }
  }

  return feed, err
}

func (rss2Entry *RSS2Entry) Marshal() (entry *Entry, err error) {
  guid := rss2Entry.Id
  if guid == "" {
    guid = rss2Entry.Link
  }

  content := rss2Entry.EncodedContent
  if content == "" {
    content = rss2Entry.Content
  }

  published := time.Time {}
  if rss2Entry.Published != "" {
    published, err = AsTime(rss2Entry.Published)
  }

  entry = &Entry {
    GUID: guid,
    Author: rss2Entry.Author,
    Title: rss2Entry.EntryTitle,
    Content: content,
    Published: published,
    WWWURL: rss2Entry.Link,
  }

  return entry, err
}

func AsTime(timeSpec string) (time.Time, error) {
  if timeSpec != "" {
    if parsedTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", timeSpec); err == nil {
      return parsedTime, nil
    } else if parsedTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", timeSpec); err == nil {
      return parsedTime, nil
    } else {
      return time.Time {}, errors.New("Unrecognized time format: " + timeSpec)
    }
  }

  return time.Time {}, nil
} 