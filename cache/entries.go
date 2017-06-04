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
	Name   string `json:"name"`   // Username
	Passwd string `json:"passwd"` // Password
	UID    uint32 `json:"uid"`    // User ID
	GID    uint32 `json:"gid"`    // Group ID
	GECOS  string `json:"gecos"`  // Real name
	Dir    string `json:"dir"`    // Home directory
	Shell  string `json:"shell"`  // Shell program
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
	Name   string   `json:"name"`   // Group name
	Passwd string   `json:"passwd"` // Password
	GID    uint32   `json:"gid"`    // Group ID
	Mem    []string `json:"mem"`    // Member list
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
