package common

import (
	"github.com/kennygrant/sanitize"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"regexp"
)

var fileNameRegex = regexp.MustCompile(`(.*)[\.^-](.*)$`)
var extFilename = regexp.MustCompile(`\..*`)

func SafeFilename(s string) string {
	san := sanitize.BaseName(s)
	if extFilename.MatchString(s) {
		return fileNameRegex.ReplaceAllString(san, "$1.$2")
	}
	return san
}

func GetFileContentType(out *os.File, l logrus.FieldLogger) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	if contentType == "application/octet-stream" {
		l.WithField("fileName", out.Name()).Warn("contentType was not detected properly for file")
	}

	return contentType, nil
}
