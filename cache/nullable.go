package cache

import "fmt"

type nullInt32 struct {
	val   int32
	valid bool
}

func (n *nullInt32) String() string {
	if !n.valid {
		return ""
	}
	return fmt.Sprintf("%d", n.val)
}

func Int32(v int32) nullInt32 {
	return nullInt32{
		val:   v,
		valid: true,
	}
}

type nullUInt32 struct {
	val   uint32
	valid bool
}

func (n *nullUInt32) String() string {
	if !n.valid {
		return ""
	}
	return fmt.Sprintf("%d", n.val)
}

func UInt32(v uint32) nullUInt32 {
	return nullUInt32{
		val:   v,
		valid: true,
	}
}
