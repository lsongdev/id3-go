package v2

func ParseID3v24FrameSize(data []byte) int {
	return int(parseSize(data))
}

var V24FrameTypeMap = V23FrameTypeMap

var V24FrameMapping = map[string]string{
	"title":   "TIT2",
	"artist":  "TPE1",
	"album":   "TALB",
	"year":    "TDRC",
	"comment": "COMM",
	"track":   "TRCK",
	"genre":   "TCON",
}
