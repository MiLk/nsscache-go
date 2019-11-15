package cache

import (
	"fmt"
	"io"
	"strings"
)

// Entry specifies a generic entry in an unspecified cache.  Specific
// implementations are provided for passwd, group, and shadow caches.
type Entry interface {
	fmt.Stringer
	io.WriterTo

	Column(int) string
}

// PasswdEntry describes an entry of the /etc/passwd file
// https://sourceware.org/git/?p=glibc.git;a=blob;f=pwd/pwd.h;hb=HEAD#l49
// https://fossies.org/dox/glibc-2.30/structpasswd.html
type PasswdEntry struct {
	Name   string `json:"name"`   // Username
	Passwd string `json:"passwd"` // Password
	UID    uint32 `json:"uid"`    // User ID
	GID    uint32 `json:"gid"`    // Group ID
	GECOS  string `json:"gecos"`  // Real name
	Dir    string `json:"dir"`    // Home directory
	Shell  string `json:"shell"`  // Shell program
}

func (e *PasswdEntry) format() string {
	return "%s:%s:%d:%d:%s:%s:%s\n"
}

func (e *PasswdEntry) args() []interface{} {
	if e.Passwd == "" {
		e.Passwd = "x"
	}

	return []interface{}{
		e.Name,
		e.Passwd,
		e.UID,
		e.GID,
		e.GECOS,
		e.Dir,
		e.Shell,
	}
}

// Column returns the information from the requested columns or an
// empty string if no column is known.
func (e *PasswdEntry) Column(col int) string {
	switch col {
	case 0:
		return e.Name
	case 2:
		return fmt.Sprintf("%d", e.UID)
	default:
		return ""
	}
}

func (e *PasswdEntry) String() string {
	return fmt.Sprintf(e.format(), e.args()...)
}

// WriteTo writes the specified entry to the provided writer.
func (e *PasswdEntry) WriteTo(w io.Writer) (int64, error) {
	return toInt64(fmt.Fprintf(w, e.format(), e.args()...))
}

// ShadowEntry describes an entry of the /etc/shadow file
// https://sourceware.org/git/?p=glibc.git;a=blob;f=shadow/shadow.h;hb=HEAD#l39
// https://fossies.org/dox/glibc-2.30/structspwd.htmls
type ShadowEntry struct {
	Name   string     `json:"name"`             // Login name
	Passwd string     `json:"passwd"`           // Encrypted password
	Lstchg nullInt32  `json:"lstchg,omitempty"` // Date of last change
	Min    nullInt32  `json:"min,omitempty"`    // Minimum number of days between changes
	Max    nullInt32  `json:"max,omitempty"`    // Maximum number of days between changes
	Warn   nullInt32  `json:"warn,omitempty"`   // Number of days to warn user to change the password
	Inact  nullInt32  `json:"inact,omitempty"`  // Number of days the account may be inactive
	Expire nullInt32  `json:"expire,omitempty"` // Number of days since 1970-01-01 until account expires
	Flag   nullUInt32 `json:"flag,omitempty"`   // Reserved
}

func (e *ShadowEntry) format() string {
	return "%s:%s:%s:%s:%s:%s:%s:%s:%s\n"
}

func (e *ShadowEntry) args() []interface{} {
	if e.Passwd == "" {
		e.Passwd = "!!"
	}

	return []interface{}{
		e.Name,
		e.Passwd,
		e.Lstchg.String(),
		e.Min.String(),
		e.Max.String(),
		e.Warn.String(),
		e.Inact.String(),
		e.Expire.String(),
		e.Flag.String(),
	}
}

// Column returns the information from the requested columns or an
// empty string if no column is known.
func (e *ShadowEntry) Column(col int) string {
	switch col {
	case 0:
		return e.Name
	default:
		return ""
	}
}

func (e *ShadowEntry) String() string {
	return fmt.Sprintf(e.format(), e.args()...)
}

// WriteTo writes the specified entry to the provided writer.
func (e *ShadowEntry) WriteTo(w io.Writer) (int64, error) {
	return toInt64(fmt.Fprintf(w, e.format(), e.args()...))
}

// GroupEntry describes an entry of the /etc/group file
// https://sourceware.org/git/?p=glibc.git;a=blob;f=grp/grp.h;hb=HEAD#l41
// https://fossies.org/dox/glibc-2.30/structgroup.html
type GroupEntry struct {
	Name   string   `json:"name"`   // Group name
	Passwd string   `json:"passwd"` // Password
	GID    uint32   `json:"gid"`    // Group ID
	Mem    []string `json:"mem"`    // Member list
}

func (e *GroupEntry) format() string {
	return "%s:%s:%d:%s\n"
}

func (e *GroupEntry) args() []interface{} {
	if e.Passwd == "" {
		e.Passwd = "x"
	}

	return []interface{}{
		e.Name,
		e.Passwd,
		e.GID,
		strings.Join(e.Mem, ","),
	}
}

// Column returns the information from the requested columns or an
// empty string if no column is known.
func (e *GroupEntry) Column(col int) string {
	switch col {
	case 0:
		return e.Name
	case 2:
		return fmt.Sprintf("%d", e.GID)
	default:
		return ""
	}
}

func (e *GroupEntry) String() string {
	return fmt.Sprintf(e.format(), e.args()...)
}

// WriteTo writes the specified entry to the provided writer.
func (e *GroupEntry) WriteTo(w io.Writer) (int64, error) {
	return toInt64(fmt.Fprintf(w, e.format(), e.args()...))
}

func toInt64(i int, e error) (int64, error) {
	return int64(i), e
}
