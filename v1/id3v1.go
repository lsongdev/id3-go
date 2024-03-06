package v1

import (
	"fmt"
	"strings"
)

type ID3v1Tag struct {
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Album   string `json:"album"`
	Year    string `json:"year"`
	Track   string `json:"track"`
	Genre   string `json:"genre"`
	Comment string `json:"comment"`
}

func trimString(data []byte) string {
	return strings.TrimRight(string(data), "\u0000")
}

func GetGenre(i int) string {
	if i > len(ID3v1Genres)-1 {
		return "Unspecified"
	}
	return ID3v1Genres[i]
}

// ParseID3v1Tag parses the ID3v1 tag provided in the data argument and returns
// an ID3v1Tag struct with parsed strings from the tag for each field.
func ParseID3v1Tag(data []byte) (*ID3v1Tag, error) {
	if string(data[0:3]) != "TAG" {
		return nil, fmt.Errorf("invalid ID3v1 header: %s", string(data[0:3]))
	}
	tag := new(ID3v1Tag)
	tag.Title = trimString(data[3:33])
	tag.Artist = trimString(data[33:63])
	tag.Album = trimString(data[63:93])
	tag.Year = trimString(data[93:97])
	if data[125] == 0 && data[126] != 0 {
		tag.Track = fmt.Sprint(data[126])
		tag.Comment = trimString(data[97:125])
	} else {
		tag.Comment = trimString(data[97:127])
	}
	tag.Genre = GetGenre(int(data[127]))
	return tag, nil
}
