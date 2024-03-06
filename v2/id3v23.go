package v2

import (
	"bufio"
	"encoding/binary"
)

// ID3 v2.3 doesn't use sync-safe frame sizes: read in as a regular big endian number.
func parseID3v23Size(reader *bufio.Reader) int {
	var size int32
	binary.Read(reader, binary.BigEndian, &size)
	return int(size)
}

func parseID3v23File(reader *bufio.Reader, file *ID3v2Tag) {
	for hasFrame(reader, 4) {
		id := string(readBytes(reader, 4))
		size := parseID3v23Size(reader)

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
		case "TYER":
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
