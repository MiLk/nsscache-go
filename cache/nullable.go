package cache

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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

func (n *nullInt32) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	n.val = 0
	n.valid = false
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	n.val = int32(v)
	n.valid = true
	return nil
}

func (n *nullInt32) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
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

func (n *nullUInt32) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	n.val = 0
	n.valid = false
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	n.val = uint32(v)
	n.valid = true
	return nil
}

func (n *nullUInt32) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}
