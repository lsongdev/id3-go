package v2

import (
	"bufio"
)

// ID3 v2.2 uses 24-bit big endian frame sizes.
func parseID3v22FrameSize(reader *bufio.Reader) int {
	size := readBytes(reader, 3)
	return int(size[0])<<16 | int(size[1])<<8 | int(size[2])
}

func parseID3v22File(reader *bufio.Reader, file *ID3v2Tag) {
	for hasFrame(reader, 3) {
		id := string(readBytes(reader, 3))
		size := parseID3v22FrameSize(reader)

		switch id {
		case "TAL":
			file.Album = readString(reader, size)
		case "TRK":
			file.Track = readString(reader, size)
		case "TP1":
			file.Artist = readString(reader, size)
		case "TT2":
			file.Title = readString(reader, size)
		case "TYE":
			file.Year = readString(reader, size)
		case "TPA":
			file.Disc = readString(reader, size)
		case "TCO":
			file.Genre = readGenre(reader, size)
		default:
			skipBytes(reader, size)
		}
	}
}
