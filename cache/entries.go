package cache

import (
	"fmt"
	"strings"
)

type Entry interface {
	fmt.Stringer
}

// PasswdEntry describes an entry of the /etc/passwd file
// https://sourceware.org/git/?p=glibc.git;a=blob;f=pwd/pwd.h;hb=HEAD#l49
// https://fossies.org/dox/glibc-2.25/structpasswd.html
type PasswdEntry struct {
	Name   string // Username
	Passwd string // Password
	UID    uint32 // User ID
	GID    uint32 // Group ID
	GECOS  string // Real name
	Dir    string // Home directory
	Shell  string // Shell program
}

func (e *PasswdEntry) String() string {
	if e.Passwd == "" {
		e.Passwd = "x"
	}

	return fmt.Sprintf("%s:%s:%d:%d:%s:%s:%s\n",
		e.Name,
		e.Passwd,
		e.UID,
		e.GID,
		e.GECOS,
		e.Dir,
		e.Shell,
	)
}

// ShadowEntry describes an entry of the /etc/shadow file
// https://sourceware.org/git/?p=glibc.git;a=blob;f=shadow/shadow.h;hb=HEAD#l39
// https://fossies.org/dox/glibc-2.25/structspwd.htmls
type ShadowEntry struct {
	Name   string     // Login name
	Passwd string     // Encrypted password
	Lstchg nullInt32  // Date of last change
	Min    nullInt32  // Minimum number of days between changes
	Max    nullInt32  // Maximum number of days between changes
	Warn   nullInt32  // Number of days to warn user to change the password
	Inact  nullInt32  // Number of days the account may be inactive
	Expire nullInt32  // Number of days since 1970-01-01 until account expires
	Flag   nullUInt32 // Reserved
}

func (e *ShadowEntry) String() string {
	if e.Passwd == "" {
		e.Passwd = "!!"
	}

	return fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s:%s\n",
		e.Name,
		e.Passwd,
		e.Lstchg.String(),
		e.Min.String(),
		e.Max.String(),
		e.Warn.String(),
		e.Inact.String(),
		e.Expire.String(),
		e.Flag.String(),
	)
}

// GroupEntry describes an entry of the /etc/group file
// https://sourceware.org/git/?p=glibc.git;a=blob;f=grp/grp.h;hb=HEAD#l41
// https://fossies.org/dox/glibc-2.25/structgroup.html
type GroupEntry struct {
	Name   string   // Group name
	Passwd string   // Passwords
	GID    uint32   // Group ID
	Mem    []string // Member list
}

func (e *GroupEntry) String() string {
	if e.Passwd == "" {
		e.Passwd = "x"
	}

	return fmt.Sprintf("%s:%s:%d:%s\n",
		e.Name,
		e.Passwd,
		e.GID,
		strings.Join(e.Mem, ","),
	)
}
