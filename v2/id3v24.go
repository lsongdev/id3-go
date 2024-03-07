package v2

func ParseID3v24FrameSize(data []byte) int {
	return int(parseSize(data))
}

var V24FrameTypeMap = V23FrameTypeMap
