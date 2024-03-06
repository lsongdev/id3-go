package v2

import (
	"bufio"
)

// ID3 v2.4 uses sync-safe frame sizes similar to those found in the header.
func parseID3v24Size(reader *bufio.Reader) int {
	return int(parseSize(readBytes(reader, 4)))
}

func parseID3v24File(reader *bufio.Reader, file *ID3v2Tag) {
	for hasFrame(reader, 4) {
		id := string(readBytes(reader, 4))
		size := parseID3v24Size(reader)

		// Skip over frame flags.
		skipBytes(reader, 2)

		switch id {
		case "TALB":
			file.Album = readString(reader, size)
		case "TRCK":
			file.Track = readString(reader, size)
		case "TPE1":
			file.Artist = readString(reader, size)
		case "TCON":
			file.Genre = readGenre(reader, size)
		case "TIT2":
			file.Title = readString(reader, size)
		case "TDRC":
			// TODO: implement timestamp parsing
			file.Year = readString(reader, size)
		case "TPOS":
			file.Disc = readString(reader, size)
		case "TLEN":
			file.Length = readString(reader, size)
		default:
			skipBytes(reader, size)
		}
	}
}
