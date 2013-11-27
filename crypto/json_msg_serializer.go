package crypto

import (
	"encoding/json"
)

type JsonMsgSerializer struct {
}

func (s JsonMsgSerializer) Serialize(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s JsonMsgSerializer) Unserialize(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}
