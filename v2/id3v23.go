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

func parseID3v23File(reader *bufio.Reader, tag *ID3v2Tag) {
	for hasFrame(reader, 4) {
		id := string(readBytes(reader, 4))
		size := parseID3v23Size(reader)

		// Skip over frame flags.
		skipBytes(reader, 2)

		frame := &ID3v2Frame{
			Id:   id,
			Data: readBytes(reader, size),
		}
		tag.Frames = append(tag.Frames, frame)
	}
	for _, frame := range tag.Frames {
		switch frame.Id {
		case "TALB":
			tag.Album = parseString(frame.Data)
		case "TRCK":
			tag.Track = parseString(frame.Data)
		case "TPE1":
			tag.Artist = parseString(frame.Data)
		case "TCON":
			tag.Genre = parseString(frame.Data)
		case "TIT2":
			tag.Title = parseString(frame.Data)
		case "TYER":
			tag.Year = parseString(frame.Data)
		case "TPOS":
			tag.Disc = parseString(frame.Data)
		case "TLEN":
			tag.Length = parseString(frame.Data)
		}
	}
}
