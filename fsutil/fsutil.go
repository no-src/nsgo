package fsutil

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	isNotExist = os.IsNotExist
	abs        = filepath.Abs
	rel        = filepath.Rel

	errFileSysInfoIsNil = errors.New("file sys info is nil")
)

// StatFunc the function prototype of os.Stat
type StatFunc func(name string) (fs.FileInfo, error)

// GetFileTimeFunc the function prototype of GetFileTime
type GetFileTimeFunc func(path string) (cTime time.Time, aTime time.Time, mTime time.Time, err error)

// IsDirFunc the function prototype of IsDir
type IsDirFunc func(path string) (bool, error)

// FileExist is file Exist
func FileExist(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err != nil && isNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CreateFile create a file without truncate
func CreateFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
}

// OpenRWFile open a file with read write mode
func OpenRWFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR, 0666)
}

// IsDir the path is directory or not
func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

// IsEOF whether the error is io.EOF
func IsEOF(err error) bool {
	return err != nil && errors.Is(err, io.EOF)
}

// IsNonEOF whether the error is not io.EOF and is not nil
func IsNonEOF(err error) bool {
	return err != nil && !errors.Is(err, io.EOF)
}

// GetFileTime get the creation time, last access time, last modify time of the path
func GetFileTime(path string) (cTime time.Time, aTime time.Time, mTime time.Time, err error) {
	stat, err := os.Lstat(path)
	if err != nil {
		return
	}
	return GetFileTimeBySys(stat.Sys())
}

// IsSub whether it is a subdirectory of the parent
func IsSub(parent, child string) (bool, error) {
	pAbs, err := abs(parent)
	if err != nil {
		return false, err
	}
	cAbs, err := abs(child)
	if err != nil {
		return false, err
	}
	relPath, err := rel(pAbs, cAbs)
	if err != nil {
		return false, err
	}
	relPath = filepath.ToSlash(relPath)
	return relPath != ".." && !strings.HasPrefix(relPath, "../"), nil
}

// SafePath encode some special characters for the path like "?", "#" etc.
func SafePath(path string) string {
	if len(path) == 0 {
		return path
	}
	filterRule := []struct {
		old string
		new string
	}{
		{"%", "%25"}, // stay in the first position
		{"?", "%3F"},
		{"#", "%23"},
	}
	for _, r := range filterRule {
		path = strings.ReplaceAll(path, r.old, r.new)
	}
	return path
}

// IsSymlink the path is a symbolic link or not
func IsSymlink(path string) (bool, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return IsSymlinkMode(fi.Mode()), nil
}

// IsSymlinkMode check the mode is a symbolic link or not
func IsSymlinkMode(mode fs.FileMode) bool {
	return mode&fs.ModeSymlink != 0
}

// IsSymlinkSupported checks if the system supports symbolic links
func IsSymlinkSupported() bool {
	symlink := filepath.Join(os.TempDir(), fmt.Sprintf("%d.symlink_detect", time.Now().UnixNano()))
	defer os.RemoveAll(symlink)
	return Symlink(os.Args[0], symlink) == nil
}

// Symlink create a symbolic link
func Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

// Readlink returns the destination of the named symbolic link
func Readlink(name string) (string, error) {
	return os.Readlink(name)
}

// SymlinkText build custom symlink text content
func SymlinkText(realPath string) string {
	return fmt.Sprintf("# symlink\n%s", realPath)
}
