package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	data := []string{"1336589085397450752", "1336589093232410624"}

	str, _ := json.Marshal(&data)

	fmt.Println(string(str))

	test := []Int64Str{}

	err := json.Unmarshal(str, &test)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(test)
}
