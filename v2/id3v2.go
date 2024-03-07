package v2

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf16"

	v1 "github.com/song940/id3-go/v1"
)

type ID3v2Framer interface {
	String() string
}

type ID3v2Frame struct {
	Id          string
	Description string
	Raw         []byte
	Data        ID3v2Framer
}

// FrameType holds frame id metadata and constructor method
// A set number of these are created in the version specific files
type FrameType struct {
	id          string
	description string
	constructor func([]byte) (ID3v2Framer, error)
}

// A parsed ID3 file with common fields exposed.
type ID3v2Tag struct {
	Header *ID3v2Header
	Frames []*ID3v2Frame
}

func (h *ID3v2Tag) Version() string {
	return fmt.Sprintf("2.%d.%d", h.Header.Version, h.Header.Revision)
}

var skipBuffer []byte = make([]byte, 1024*4)

// A parsed ID3v2 header as defined in Section 3 of
// http://id3.org/id3v2.4.0-structure
type ID3v2Header struct {
	Version           int
	Revision          int
	Unsynchronization bool
	Extended          bool
	Experimental      bool
	Footer            bool
	Size              int32
}

func isID3Tag(reader *bufio.Reader) bool {
	data, err := reader.Peek(3)
	if len(data) < 3 || err != nil {
		return false
	}
	return string(data[0:3]) == "ID3"
}

func ParseID3v2Header(reader *bufio.Reader) (*ID3v2Header, error) {
	if !isID3Tag(reader) {
		return nil, fmt.Errorf("invalid ID3 header")
	}
	data, err := readBytes(reader, 10)
	if err != nil {
		return nil, err
	}
	h := new(ID3v2Header)
	h.Version = int(data[3])
	h.Revision = int(data[4])
	h.Unsynchronization = data[5]&1<<7 != 0
	h.Extended = data[5]&1<<6 != 0
	h.Experimental = data[5]&1<<5 != 0
	h.Footer = data[5]&1<<4 != 0
	h.Size = parseSize(data[6:])
	return h, nil
}

// Parse the input for ID3 information. Returns nil if parsing failed or the
// input didn't contain ID3 information.
func Read(reader io.Reader) (tag *ID3v2Tag, err error) {
	bufReader := bufio.NewReader(reader)
	h, err := ParseID3v2Header(bufReader)
	if err != nil {
		return
	}
	tag = &ID3v2Tag{
		Header: h,
	}
	frameReader := NewID3v2FrameParser(h.Version)
	rd := io.LimitReader(bufReader, int64(h.Size))
	limitReader := bufio.NewReader(rd)
	for frameReader.hasFrame(limitReader) {
		frame, err := frameReader.readFrame(limitReader)
		if err != nil {
			return nil, err
		}
		t, ok := frameReader.FrameTypeMap[frame.Id]
		if !ok {
			return nil, fmt.Errorf("unknown frame type: %s", frame.Id)
		}
		frame.Description = t.description
		frame.Data, err = t.constructor(frame.Raw)
		if err != nil {
			return nil, err
		}
		tag.Frames = append(tag.Frames, frame)
	}
	return
}

type ID3v2FrameParser struct {
	HeaderLen    int
	IdLen        int
	SizeLen      int
	SizeParser   func([]byte) int
	FrameTypeMap map[string]FrameType
}

// NewID3v2FrameParser creates a new frame parser initialized for the proper
// minor version settings.
func NewID3v2FrameParser(version int) *ID3v2FrameParser {
	p := new(ID3v2FrameParser)
	switch version {
	case 2:
		p.HeaderLen = 6
		p.IdLen = 3
		p.SizeLen = 3
		p.SizeParser = ParseID3v22FrameSize
		p.FrameTypeMap = V22FrameTypeMap
	case 3:
		p.HeaderLen = 10
		p.IdLen = 4
		p.SizeLen = 4
		p.SizeParser = ParseID3v23FrameSize
		p.FrameTypeMap = V23FrameTypeMap
	case 4:
		p.HeaderLen = 10
		p.IdLen = 4
		p.SizeLen = 4
		p.SizeParser = ParseID3v24FrameSize
		p.FrameTypeMap = V23FrameTypeMap
	}
	return p
}

// Peeks at the buffer to see if there is a valid frame.
func (p *ID3v2FrameParser) hasFrame(reader *bufio.Reader) bool {
	data, err := reader.Peek(p.SizeLen)
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

func (parser *ID3v2FrameParser) readFrame(reader *bufio.Reader) (frame *ID3v2Frame, err error) {
	idBytes, err := readBytes(reader, parser.IdLen)
	if err != nil {
		return nil, err
	}
	id := string(idBytes)
	sizeBytes, err := readBytes(reader, parser.SizeLen)
	if err != nil {
		return nil, err
	}
	size := parser.SizeParser(sizeBytes)
	skipBytes(reader, 2)
	data, err := readBytes(reader, size)
	if err != nil {
		return nil, err
	}
	frame = &ID3v2Frame{
		Id:  id,
		Raw: data,
	}
	return
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

// Parses a string from frame data. The first byte represents the encoding:
//
//	0x01  ISO-8859-1
//	0x02  UTF-16 w/ BOM
//	0x03  UTF-16BE w/o BOM
//	0x04  UTF-8
//
// Refer to section 4 of http://id3.org/id3v2.4.0-structure
func parseString(data []byte) (s string, err error) {
	switch data[0] {
	case 0x00: // ISO-8859-1 text.
		s = ISO8859_1ToUTF8(data[1:])
		break
	case 0x01: // UTF-16 with BOM.
		s = string(utf16.Decode(toUTF16(data[1:])))
		break
	case 0x02: // UTF-16BE without BOM.
		err = fmt.Errorf("unsupported text encoding UTF-16BE")
		return
	case 0x03: // UTF-8 text.
		s = string(data[1:])
		break
	default:
		// No encoding, assume ISO-8859-1 text.
		// s = ISO8859_1ToUTF8(data)
		err = fmt.Errorf("Unsupported text encoding: 0x" + fmt.Sprintf("%02X", data[0]))
		return
	}
	return strings.TrimRight(s, "\u0000"), nil
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

	// Try to parse "NN" format.
	index, err := strconv.Atoi(genre)
	if err == nil {
		return v1.GetGenre(index)
	}
	// Try to parse "(NN)" format.
	_, err = fmt.Sscanf(genre, "(%d)", &index)
	if err == nil {
		return v1.GetGenre(index)
	}
	// Couldn't parse so it's likely not an ID3v1 genre.
	return genre
}

func readGenre(reader *bufio.Reader, c int) (string, error) {
	data, err := readBytes(reader, c)
	if err != nil {
		return "", err
	}
	genre, err := parseString(data)
	if err != nil {
		return "", err
	}
	return convertID3v1Genre(genre), nil
}

func readBytes(reader *bufio.Reader, c int) ([]byte, error) {
	b := make([]byte, c)
	pos := 0
	for pos < c {
		i, err := reader.Read(b[pos:])
		pos += i
		if err != nil {
			return nil, err
		}
	}
	return b, nil
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

func (tag *ID3v2Tag) Title() string {
	for _, frame := range tag.Frames {
		if frame.Id == "TIT2" {
			return frame.Data.String()
		}
	}
	return ""
}
