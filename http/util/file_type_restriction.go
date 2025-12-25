package util

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gw-gong/gwkit-go/log"
)

type FileTypeRestriction struct {
	Enable    bool     `json:"enable" yaml:"enable" mapstructure:"enable"`
	WhiteList []string `json:"white_list" yaml:"white_list" mapstructure:"white_list"`
	BlackList []string `json:"black_list" yaml:"black_list" mapstructure:"black_list"`
}

func (f *FileTypeRestriction) IsFileTypeAllowedWithReader(ctx context.Context, contentType string, reader io.ReadSeeker) (isAllowed bool) {
	var (
		detectedContentType string
		err                 error
	)
	defer func() {
		log.Infoc(ctx, "File type restriction is file type allowed with reader done",
			log.Str("contentType", contentType),
			log.Str("detectedContentType", detectedContentType),
			log.Bool("isAllowed", isAllowed),
			log.Err(err),
		)
	}()

	if !f.isFileTypeAllowed(contentType) {
		return false
	}

	if reader == nil {
		err = errors.New("reader is nil")
		return false
	}

	buf := make([]byte, 512)
	_, err = reader.Read(buf)
	if err != nil {
		err = errors.New("failed to read reader")
		return false
	}
	defer func() { // !important
		if _, err = reader.Seek(0, io.SeekStart); err != nil {
			err = errors.New("failed to seek reader")
			isAllowed = false
		}
	}()
	detectedContentType = http.DetectContentType(buf)

	// files with special processing formats
	if detectedContentType == "application/zip" {
		detectedContentType = contentType
	}

	return f.isFileTypeAllowed(detectedContentType)
}

func (f *FileTypeRestriction) isFileTypeAllowed(contentType string) (isAllowed bool) {
	if !f.Enable {
		return true
	}

	if !f.matchWhiteList(contentType) {
		return false
	}
	if f.matchBlackList(contentType) {
		return false
	}
	return true
}

func (f *FileTypeRestriction) matchWhiteList(mimeType string) bool {
	for _, rule := range f.WhiteList {
		if f.isPatternMatch(rule, mimeType) {
			return true
		}
	}
	return false
}

func (f *FileTypeRestriction) matchBlackList(mimeType string) bool {
	for _, rule := range f.BlackList {
		if f.isPatternMatch(rule, mimeType) {
			return true
		}
	}
	return false
}

// isPatternMatch Check whether the rule matches the MIME type (supports wildcard characters such as image/*)
func (f *FileTypeRestriction) isPatternMatch(rule string, mimeType string) bool {
	if rule == "*/*" { // all types are allowed / disabled
		return true
	}
	if strings.HasSuffix(rule, "/*") {
		baseType := strings.TrimSuffix(rule, "/*")
		return strings.HasPrefix(mimeType, baseType+"/")
	}
	return rule == mimeType
}
