package archive

/*
#cgo pkg-config: libarchive
#include <archive.h>
#include <archive_entry.h>
#include <stdlib.h>
*/
import "C"

import (
	"os"
	"time"
)

// Entry represents an individual Entry in an archive.
// An Entry is roughly analogous to an inode, and it satisfies the
// FileInfo interface.
// Reader.Next() iterates over these, and Reader.Read() reads
// the Entry that the Reader currently points at.
type Entry struct {
	entry *C.struct_archive_entry
}

// The full path for an entry.
func (e *Entry) PathName() string {
	name := C.archive_entry_pathname(e.entry)
	return C.GoString(name)
}

// Same as above, satisfies the FileInfo interface.
func (e *Entry) Name() string {
	return e.PathName()
}

// The size of an entry.  This may be undefined for
// directories.
func (e *Entry) Size() int64 {
	return int64(C.archive_entry_size(e.entry))
}

// The UNIX mode bits for an archive.  This is not always
// meaningful depending on where the archive was created.
func (e *Entry) Mode() os.FileMode {
	return os.FileMode(C.archive_entry_mode(e.entry))
}

// The last time an entry in the archive as modified.
func (e *Entry) ModTime() time.Time {
	if C.archive_entry_mtime_is_set(e.entry) != 0 {
		return time.Unix(int64(C.archive_entry_mtime(e.entry)),
			int64(C.archive_entry_mtime_nsec(e.entry)))
	} else {
		return time.Now()
	}
}

// Return the libarchive specific file type.
// Defined in constants.go.
func (e *Entry) FileType() int {
	return int(C.archive_entry_filetype(e.entry))
}

// Whether this entry is for a directory.
// Size and contents are not meaningful if this is true.
func (e *Entry) IsDir() bool {
	return  e.FileType() == FileTypeDir
}

// Whether this entry is for a regular file.
// Size and contents are meaningful if this is true.
func (e *Entry) IsRegular() bool {
	return e.FileType() == FileTypeRegFile
}

// Whether this entry describes a symlink.
// Size and contents are menaingful, but the contents
// contain a path.
func (e *Entry) IsSymLink() bool {
	return e.FileType() == FileTypeSymLink
}

// Present to allow an Entry to satisfy the FileInfo interface.
func (e *Entry) Sys() interface{} {
	return nil // fix
}
