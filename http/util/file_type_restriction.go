package util

import (
	"mime"
	"path/filepath"
	"strings"

	"github.com/gw-gong/gwkit-go/log"
)

/*
xxx:
  enable: true
  white_list:
    - image/*
  black_list:
    - image/svg+xml
*/

type FileTypeRestriction struct {
	Enable    bool     `yaml:"enable"`
	WhiteList []string `yaml:"white_list"`
	BlackList []string `yaml:"black_list"`
}

// IsFileTypeAllowed Check if the file type is allowed (preferred to use Content-Type in the HTTP header)
func (f *FileTypeRestriction) IsFileTypeAllowed(contentType string, filename string) (isAllowed bool) {
	if !f.Enable {
		return true
	}

	mimeType := contentType
	if mimeType == "" {
		ext := strings.ToLower(filepath.Ext(filename))
		mimeType = mime.TypeByExtension(ext)
	}
	defer func() {
		log.Info("FileTypeRestriction.IsFileTypeAllowed",
			log.Any("FileTypeRestriction", f),
			log.Str("contentType", contentType),
			log.Str("filename", filename),
			log.Str("mimeType", mimeType),
			log.Bool("isAllowed", isAllowed),
		)
	}()

	if !f.matchWhiteList(mimeType) {
		return false
	}
	if f.matchBlackList(mimeType) {
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
