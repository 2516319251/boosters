package log

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestA(t *testing.T)  {
	type Data struct {
		Field string `json:"field"`
	}

	encode := Data{Field: "xxxxxxxxx"}
	b, e := json.Marshal(encode)
	if e != nil {
		t.Error(e)
	}
	fmt.Println(string(b))

	decode := new(Data)
	if e = json.Unmarshal(b, decode); e != nil {
		t.Error(e)
	}
	fmt.Println(decode)
}

func TestLevel(t *testing.T) {
	encode := map[string]interface{}{
		"a": 12,
		"b": "xxx",
	}
	b, e := json.Marshal(encode)
	if e != nil {
		t.Error(e)
	}
	fmt.Println(string(b))

	decode := make(map[string]interface{})
	if e = json.Unmarshal(b, &decode); e != nil {
		t.Error(e)
	}
	fmt.Println(decode)
}
