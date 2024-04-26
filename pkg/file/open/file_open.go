package gfopen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var suffixMatcher = regexp.MustCompile(`_(\d+)(\.[^.\s\n]+)$`)

func createParentDirs(path string) error {
	// get parent directory
	if i := strings.LastIndexByte(path, filepath.Separator); i > 1 {
		p := path[:i]

		// create the parent directory if it does not exist
		return os.MkdirAll(p, 0755)
	}
	return nil
}

type OpenOption func(*openOption)
type openOption struct {
	// string: path, bool: isFileExists, fs.FileInfo: fileInfo
	incrementalSuffixIfExistsValidator                                         func(string, fs.FileInfo) bool
	perm                                                                       fs.FileMode
	disablePathCreation, isReadOnly, isWriteOnly, isTruncate, mustExistsBefore bool
}

func WithDisablePathCreation() OpenOption {
	return func(o *openOption) {
		o.disablePathCreation = true
	}
}

func WithMustExistsBefore() OpenOption {
	return func(o *openOption) {
		o.mustExistsBefore = true
	}
}

func WithIncrementalSuffixIfExists(fn func(string, fs.FileInfo) bool) OpenOption {
	return func(o *openOption) {
		o.incrementalSuffixIfExistsValidator = fn
	}
}

func WithReadOnly() OpenOption {
	return func(o *openOption) {
		o.isReadOnly = true
	}
}

func WithWriteOnly() OpenOption {
	return func(o *openOption) {
		o.isWriteOnly = true
	}
}

func WithTruncate() OpenOption {
	return func(o *openOption) {
		o.isTruncate = true
	}
}

func WithPerm(perm fs.FileMode) OpenOption {
	return func(o *openOption) {
		o.perm = perm
	}
}

func OpenFile(path string, opts ...OpenOption) (*os.File, error) {
	var err error
	// default options
	o := &openOption{
		perm: 0644,
	}

	// update the options
	for _, opt := range opts {
		opt(o)
	}

	// create parent directory if it does not exist
	if !o.disablePathCreation {
		err = createParentDirs(path)
		if err != nil {
			return nil, err
		}
	}

	// get the flags
	flag := os.O_RDWR
	if o.isReadOnly {
		flag = os.O_RDONLY
	}
	if o.isWriteOnly {
		flag = os.O_WRONLY
	}

	// add create flag if it must exists before
	flag |= os.O_CREATE
	if o.mustExistsBefore {
		flag |= os.O_EXCL
	}

	// append to the file if it exists
	if !o.isTruncate {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	// create incremental suffixed file if it exists
	if o.incrementalSuffixIfExistsValidator != nil {
		path, err = getIncrementalSuffixedPath(path, o.incrementalSuffixIfExistsValidator)
		if err != nil {
			return nil, err
		}
	}

	// create the file if it does not exist
	return os.OpenFile(path, flag, o.perm)
}

// FileStatIfExists returns the file info if the file exists, otherwise nil
func FileStatIfExists(path string) fs.FileInfo {
	s, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return s
}

func getIncrementalSuffixedPath(path string, validatorFn func(string, fs.FileInfo) bool) (string, error) {
	s := FileStatIfExists(path)
	if !validatorFn(path, s) {
		return path, nil
	}

	m := suffixMatcher.FindStringSubmatch(path)
	n := 1
	p1 := ""
	p2 := ""
	if len(m) > 0 {
		// get the incremental number and parse to int
		n, _ = strconv.Atoi(m[1])

		// get the path without the incremental number
		p1 = path[:len(path)-len(m[0])]

		// increment the number
		n++

		// get extension
		p2 = m[2]
	} else {
		p := strings.LastIndexByte(path, '.')
		p1 = path[:p] + "_"
		p2 = path[p:]
	}

	// create the new path
	path = p1 + fmt.Sprintf("%d", n) + p2
	for validatorFn(path, FileStatIfExists(path)) {
		n++
		path = p1 + fmt.Sprintf("%d", n) + p2
	}
	return path, nil
}
