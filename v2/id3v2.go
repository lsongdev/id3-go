package v2

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"unicode/utf16"

	v1 "github.com/song940/id3-go/v1"
)

// A parsed ID3v2 header as defined in Section 3 of
// http://id3.org/id3v2.4.0-structure
type ID3v2Header struct {
	Version           int
	MinorVersion      int
	Unsynchronization bool
	Extended          bool
	Experimental      bool
	Footer            bool
	Size              int32
}

// A parsed ID3 file with common fields exposed.
type ID3v2Tag struct {
	Header ID3v2Header

	Title  string
	Artist string
	Album  string
	Year   string
	Track  string
	Disc   string
	Genre  string
	Length string
}

var skipBuffer []byte = make([]byte, 1024*4)

func ISO8859_1ToUTF8(data []byte) string {
	p := make([]rune, len(data))
	for i, b := range data {
		p[i] = rune(b)
	}
	return string(p)
}

func toUTF16(data []byte) []uint16 {
	if len(data) < 2 {
		panic("Sequence is too short too contain a UTF-16 BOM")
	}
	if len(data)%2 > 0 {
		// TODO: if this is UTF-16 BE then this is likely encoded wrong
		data = append(data, 0)
	}

	var shift0, shift1 uint
	if data[0] == 0xFF && data[1] == 0xFE {
		// UTF-16 LE
		shift0 = 0
		shift1 = 8
	} else if data[0] == 0xFE && data[1] == 0xFF {
		// UTF-16 BE
		shift0 = 8
		shift1 = 0
		panic("UTF-16 BE found!")
	} else {
		panic(fmt.Sprintf("Unrecognized UTF-16 BOM: 0x%02X%02X", data[0], data[1]))
	}

	s := make([]uint16, 0, len(data)/2)
	for i := 2; i < len(data); i += 2 {
		s = append(s, uint16(data[i])<<shift0|uint16(data[i+1])<<shift1)
	}
	return s
}

// Peeks at the buffer to see if there is a valid frame.
func hasFrame(reader *bufio.Reader, frameSize int) bool {
	data, err := reader.Peek(frameSize)
	if err != nil {
		return false
	}

	for _, c := range data {
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

// Sizes are stored big endian but with the first bit set to 0 and always ignored.
//
// Refer to section 3.1 of http://id3.org/id3v2.4.0-structure
func parseSize(data []byte) int32 {
	size := int32(0)
	for i, b := range data {
		if b&0x80 > 0 {
			fmt.Println("Size byte had non-zero first bit")
		}

		shift := uint32(len(data)-i-1) * 7
		size |= int32(b&0x7f) << shift
	}
	return size
}

// Parses a string from frame data. The first byte represents the encoding:
//
//	0x01  ISO-8859-1
//	0x02  UTF-16 w/ BOM
//	0x03  UTF-16BE w/o BOM
//	0x04  UTF-8
//
// Refer to section 4 of http://id3.org/id3v2.4.0-structure
func parseString(data []byte) string {
	var s string
	switch data[0] {
	case 0: // ISO-8859-1 text.
		s = ISO8859_1ToUTF8(data[1:])
		break
	case 1: // UTF-16 with BOM.
		s = string(utf16.Decode(toUTF16(data[1:])))
		break
	case 2: // UTF-16BE without BOM.
		panic("Unsupported text encoding UTF-16BE.")
	case 3: // UTF-8 text.
		s = string(data[1:])
		break
	default:
		// No encoding, assume ISO-8859-1 text.
		s = ISO8859_1ToUTF8(data)
	}
	return strings.TrimRight(s, "\u0000")
}

func readBytes(reader *bufio.Reader, c int) []byte {
	b := make([]byte, c)
	pos := 0
	for pos < c {
		i, err := reader.Read(b[pos:])
		pos += i
		if err != nil {
			panic(err)
		}
	}
	return b
}

func readString(reader *bufio.Reader, c int) string {
	return parseString(readBytes(reader, c))
}

func readGenre(reader *bufio.Reader, c int) string {
	genre := parseString(readBytes(reader, c))
	return convertID3v1Genre(genre)
}

func skipBytes(reader *bufio.Reader, c int) {
	pos := 0
	for pos < c {
		end := c - pos
		if end > len(skipBuffer) {
			end = len(skipBuffer)
		}

		i, err := reader.Read(skipBuffer[0:end])
		pos += i
		if err != nil {
			panic(err)
		}
	}
}

// ID3v2.2 and ID3v2.3 use "(NN)" where as ID3v2.4 simply uses "NN" when
// referring to ID3v1 genres. The "(NN)" format is allowed to have trailing
// information.
//
// RX and CR are shorthand for Remix and Cover, respectively.
//
// Refer to the following documentation:
//
//	http://id3.org/id3v2-00          TCO frame
//	http://id3.org/id3v2.3.0         TCON frame
//	http://id3.org/id3v2.4.0-frames  TCON frame
func convertID3v1Genre(genre string) string {
	if genre == "RX" || strings.HasPrefix(genre, "(RX)") {
		return "Remix"
	}
	if genre == "CR" || strings.HasPrefix(genre, "(CR)") {
		return "Cover"
	}

	var id3v1Genres = v1.ID3v1Genres
	// Try to parse "NN" format.
	index, err := strconv.Atoi(genre)
	if err == nil {
		if index >= 0 && index < len(id3v1Genres) {
			return id3v1Genres[index]
		}
		return "Unknown"
	}

	// Try to parse "(NN)" format.
	index = 0
	_, err = fmt.Sscanf(genre, "(%d)", &index)
	if err == nil {
		if index >= 0 && index < len(id3v1Genres) {
			return id3v1Genres[index]
		}
		return "Unknown"
	}

	// Couldn't parse so it's likely not an ID3v1 genre.
	return genre
}

// Parse the input for ID3 information. Returns nil if parsing failed or the
// input didn't contain ID3 information.
func Read(reader io.Reader) *ID3v2Tag {
	file := new(ID3v2Tag)
	bufReader := bufio.NewReader(reader)
	if !isID3Tag(bufReader) {
		log.Println("No ID3v2 tag found.")
		return nil
	}

	parseID3v2Header(bufReader, file)
	limitReader := bufio.NewReader(io.LimitReader(bufReader, int64(file.Header.Size)))
	if file.Header.Version == 2 {
		parseID3v22File(limitReader, file)
	} else if file.Header.Version == 3 {
		parseID3v23File(limitReader, file)
	} else if file.Header.Version == 4 {
		parseID3v24File(limitReader, file)
	} else {
		panic(fmt.Sprintf("Unrecognized ID3v2 version: %d", file.Header.Version))
	}

	return file
}

func isID3Tag(reader *bufio.Reader) bool {
	data, err := reader.Peek(3)
	if len(data) < 3 || err != nil {
		return false
	}
	return data[0] == 'I' && data[1] == 'D' && data[2] == '3'
}

func parseID3v2Header(reader *bufio.Reader, file *ID3v2Tag) {
	data := readBytes(reader, 10)
	file.Header.Version = int(data[3])
	file.Header.MinorVersion = int(data[4])
	file.Header.Unsynchronization = data[5]&1<<7 != 0
	file.Header.Extended = data[5]&1<<6 != 0
	file.Header.Experimental = data[5]&1<<5 != 0
	file.Header.Footer = data[5]&1<<4 != 0
	file.Header.Size = parseSize(data[6:])
}
