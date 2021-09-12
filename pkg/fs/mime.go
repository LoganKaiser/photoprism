package fs

import (
	"os"

	"github.com/h2non/filetype"
)

const (
	MimeTypeJpeg   = "image/jpeg"
	MimeTypePng    = "image/png"
	MimeTypeGif    = "image/gif"
	MimeTypeBitmap = "image/bmp"
	MimeTypeTiff   = "image/tiff"
	MimeTypeHEIF   = "image/heif"
	MimeTypeMP4    = "video/mp4"
)

// MimeType returns the mime type of a file, empty string if unknown.
func MimeType(filename string) string {
	handle, err := os.Open(filename)

	if err != nil {
		return ""
	}

	defer handle.Close()

	// Only the first 261 bytes are used to sniff the content type.
	buffer := make([]byte, 261)

	if _, err := handle.Read(buffer); err != nil {
		return ""
	} else if t, err := filetype.Get(buffer); err == nil && t != filetype.Unknown {
		return t.MIME.Value
	} else if t := filetype.GetType(NormalizedExt(filename)); t != filetype.Unknown {
		return t.MIME.Value
	} else {
		return ""
	}
}
