package models

import (
	"encoding/json"
	"sort"
	"strconv"
)

type Int64Str uint64

func (i Int64Str) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(i), 10))
}

func (i *Int64Str) UnmarshalJSON(b []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		value, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*i = Int64Str(value)
		return nil
	}

	// Fallback to number
	return json.Unmarshal(b, (*uint64)(i))
}

type IDLIST []Int64Str

func (l IDLIST) Len() int {
	return len(l)
}

func (l IDLIST) Less(i, j int) bool {
	return l[i] < l[j]
}

func (l IDLIST) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

//Diff 返回在b中不在a中的数
func Diff(a, b IDLIST) IDLIST {
	c := IDLIST{}
	sort.Sort(&a)
	sort.Sort(&b)
	posb := 0
	pos := 0

	for posb < len(b) {
		element := b[posb]
		for pos < len(a) {
			elementA := a[pos]
			if elementA == element {
				posb++
				pos++
				break
			}
			if element < elementA {
				posb++
				c = append(c, element)
				break
			}
			pos++
		}

		if pos == len(a) {
			break
		}

	}

	if posb < len(b) {
		c = append(c, b[posb:]...)
	}

	return c

}
