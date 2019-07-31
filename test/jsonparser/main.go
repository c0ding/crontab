package main

import (
	"fmt"
	"github.com/buger/jsonparser"
)

var (
	data []byte
)

func main() {
	mock()
	result, err := jsonparser.GetString(data, "person", "name", "fullName")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	content, valueType, offset, err := jsonparser.Get(data, "person", "name", "fullName")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(content, valueType, offset)
	result1, err := jsonparser.ParseString(content)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result1)

	err = jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		fmt.Printf("key:%s\n value:%s\n Type:%s\n", string(key), string(value), dataType)
		return nil
	}, "person", "name")

}

func mock() {

	data = []byte(`{
  "person": {
    "name": {
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev"
    },
    "github": {
      "handle": "buger",
      "followers": 109
    },
    "avatars": [
      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
    ]
  },
  "company": {
    "name": "Acme"
  }
}`)

}
