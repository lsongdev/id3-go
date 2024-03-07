package v2

import (
	"bytes"
	"encoding/binary"
)

func ParseID3v23FrameSize(buf []byte) int {
	var size int32
	bufr := bytes.NewBuffer(buf)
	binary.Read(bufr, binary.BigEndian, &size)
	return int(size)
}

// V23FrameTypeMap specifies the frame IDs and constructors allowed in ID3v2.3
var V23FrameTypeMap = map[string]FrameType{
	"AENC": {id: "AENC", description: "Audio encryption", constructor: ParseDataFrame},
	"APIC": {id: "APIC", description: "Attached picture", constructor: ParseImageFrame},
	"COMM": {id: "COMM", description: "Comments", constructor: ParseUnsynchTextFrame},
	"COMR": {id: "COMR", description: "Commercial frame", constructor: ParseDataFrame},
	"ENCR": {id: "ENCR", description: "Encryption method registration", constructor: ParseDataFrame},
	"EQUA": {id: "EQUA", description: "Equalization", constructor: ParseDataFrame},
	"ETCO": {id: "ETCO", description: "Event timing codes", constructor: ParseDataFrame},
	"GEOB": {id: "GEOB", description: "General encapsulated object", constructor: ParseDataFrame},
	"GRID": {id: "GRID", description: "Group identification registration", constructor: ParseDataFrame},
	"IPLS": {id: "IPLS", description: "Involved people list", constructor: ParseDataFrame},
	"LINK": {id: "LINK", description: "Linked information", constructor: ParseDataFrame},
	"MCDI": {id: "MCDI", description: "Music CD identifier", constructor: ParseDataFrame},
	"MLLT": {id: "MLLT", description: "MPEG location lookup table", constructor: ParseDataFrame},
	"OWNE": {id: "OWNE", description: "Ownership frame", constructor: ParseDataFrame},
	"PRIV": {id: "PRIV", description: "Private frame", constructor: ParseDataFrame},
	"PCNT": {id: "PCNT", description: "Play counter", constructor: ParseDataFrame},
	"POPM": {id: "POPM", description: "Popularimeter", constructor: ParseDataFrame},
	"POSS": {id: "POSS", description: "Position synchronisation frame", constructor: ParseDataFrame},
	"RBUF": {id: "RBUF", description: "Recommended buffer size", constructor: ParseDataFrame},
	"RVAD": {id: "RVAD", description: "Relative volume adjustment", constructor: ParseDataFrame},
	"RVRB": {id: "RVRB", description: "Reverb", constructor: ParseDataFrame},
	"SYLT": {id: "SYLT", description: "Synchronized lyric/text", constructor: ParseDataFrame},
	"SYTC": {id: "SYTC", description: "Synchronized tempo codes", constructor: ParseDataFrame},
	"TALB": {id: "TALB", description: "Album/Movie/Show title", constructor: ParseTextFrame},
	"TBPM": {id: "TBPM", description: "BPM (beats per minute)", constructor: ParseTextFrame},
	"TCOM": {id: "TCOM", description: "Composer", constructor: ParseTextFrame},
	"TCON": {id: "TCON", description: "Content type", constructor: ParseTextFrame},
	"TCOP": {id: "TCOP", description: "Copyright message", constructor: ParseTextFrame},
	"TDAT": {id: "TDAT", description: "Date", constructor: ParseTextFrame},
	"TDLY": {id: "TDLY", description: "Playlist delay", constructor: ParseTextFrame},
	"TENC": {id: "TENC", description: "Encoded by", constructor: ParseTextFrame},
	"TEXT": {id: "TEXT", description: "Lyricist/Text writer", constructor: ParseTextFrame},
	"TFLT": {id: "TFLT", description: "File type", constructor: ParseTextFrame},
	"TIME": {id: "TIME", description: "Time", constructor: ParseTextFrame},
	"TIT1": {id: "TIT1", description: "Content group description", constructor: ParseTextFrame},
	"TIT2": {id: "TIT2", description: "Title/songname/content description", constructor: ParseTextFrame},
	"TIT3": {id: "TIT3", description: "Subtitle/Description refinement", constructor: ParseTextFrame},
	"TKEY": {id: "TKEY", description: "Initial key", constructor: ParseTextFrame},
	"TLAN": {id: "TLAN", description: "Language(s)", constructor: ParseTextFrame},
	"TLEN": {id: "TLEN", description: "Length", constructor: ParseTextFrame},
	"TMED": {id: "TMED", description: "Media type", constructor: ParseTextFrame},
	"TOAL": {id: "TOAL", description: "Original album/movie/show title", constructor: ParseTextFrame},
	"TOFN": {id: "TOFN", description: "Original filename", constructor: ParseTextFrame},
	"TOLY": {id: "TOLY", description: "Original lyricist(s)/text writer(s)", constructor: ParseTextFrame},
	"TOPE": {id: "TOPE", description: "Original artist(s)/performer(s)", constructor: ParseTextFrame},
	"TORY": {id: "TORY", description: "Original release year", constructor: ParseTextFrame},
	"TOWN": {id: "TOWN", description: "File owner/licensee", constructor: ParseTextFrame},
	"TPE1": {id: "TPE1", description: "Lead performer(s)/Soloist(s)", constructor: ParseTextFrame},
	"TPE2": {id: "TPE2", description: "Band/orchestra/accompaniment", constructor: ParseTextFrame},
	"TPE3": {id: "TPE3", description: "Conductor/performer refinement", constructor: ParseTextFrame},
	"TPE4": {id: "TPE4", description: "Interpreted, remixed, or otherwise modified by", constructor: ParseTextFrame},
	"TPOS": {id: "TPOS", description: "Part of a set", constructor: ParseTextFrame},
	"TPUB": {id: "TPUB", description: "Publisher", constructor: ParseTextFrame},
	"TRCK": {id: "TRCK", description: "Track number/Position in set", constructor: ParseTextFrame},
	"TRDA": {id: "TRDA", description: "Recording dates", constructor: ParseTextFrame},
	"TRSN": {id: "TRSN", description: "Internet radio station name", constructor: ParseTextFrame},
	"TRSO": {id: "TRSO", description: "Internet radio station owner", constructor: ParseTextFrame},
	"TSIZ": {id: "TSIZ", description: "Size", constructor: ParseTextFrame},
	"TSRC": {id: "TSRC", description: "ISRC (international standard recording code)", constructor: ParseTextFrame},
	"TSSE": {id: "TSSE", description: "Software/Hardware and settings used for encoding", constructor: ParseTextFrame},
	"TYER": {id: "TYER", description: "Year", constructor: ParseTextFrame},
	"TXXX": {id: "TXXX", description: "User defined text information frame", constructor: ParseDescTextFrame},
	"UFID": {id: "UFID", description: "Unique file identifier", constructor: ParseIdFrame},
	"USER": {id: "USER", description: "Terms of use", constructor: ParseDataFrame},
	"TCMP": {id: "TCMP", description: "Part of a compilation (iTunes extension)", constructor: ParseTextFrame},
	"USLT": {id: "USLT", description: "Unsychronized lyric/text transcription", constructor: ParseUnsynchTextFrame},
	"WCOM": {id: "WCOM", description: "Commercial information", constructor: ParseDataFrame},
	"WCOP": {id: "WCOP", description: "Copyright/Legal information", constructor: ParseDataFrame},
	"WOAF": {id: "WOAF", description: "Official audio file webpage", constructor: ParseDataFrame},
	"WOAR": {id: "WOAR", description: "Official artist/performer webpage", constructor: ParseDataFrame},
	"WOAS": {id: "WOAS", description: "Official audio source webpage", constructor: ParseDataFrame},
	"WORS": {id: "WORS", description: "Official internet radio station homepage", constructor: ParseDataFrame},
	"WPAY": {id: "WPAY", description: "Payment", constructor: ParseDataFrame},
	"WPUB": {id: "WPUB", description: "Publishers official webpage", constructor: ParseDataFrame},
	"WXXX": {id: "WXXX", description: "User defined URL link frame", constructor: ParseDataFrame},
	"TDRC": {id: "TDRC", description: "Recording date", constructor: ParseTextFrame},
}

var V23FrameMapping = map[string]string{
	"title":   "TIT2",
	"artist":  "TPE1",
	"album":   "TALB",
	"year":    "TYER",
	"comment": "COMM",
	"track":   "TRCK",
	"genre":   "TCON",
}

type DataFrame struct {
}

func (d *DataFrame) String() string {
	return ""
}

func ParseDataFrame(data []byte) (ID3v2Framer, error) {
	return &DataFrame{}, nil
}

func ParseImageFrame(data []byte) (ID3v2Framer, error) {
	return &DataFrame{}, nil
}

func ParseUnsynchTextFrame(data []byte) (ID3v2Framer, error) {
	return &DataFrame{}, nil
}

func ParseDescTextFrame(data []byte) (ID3v2Framer, error) {
	return &DataFrame{}, nil
}

func ParseIdFrame(data []byte) (ID3v2Framer, error) {
	return &DataFrame{}, nil
}

type TextFrame struct {
	Text string
}

func (t *TextFrame) String() string {
	return t.Text
}

func ParseTextFrame(data []byte) (ID3v2Framer, error) {
	str, err := parseString(data)
	if err != nil {
		return nil, err
	}
	text := &TextFrame{
		Text: str,
	}
	return text, nil
}
