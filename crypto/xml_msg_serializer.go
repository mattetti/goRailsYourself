package crypto

import (
	"encoding/xml"
)

type XMLMsgSerializer struct {
}

func (s XMLMsgSerializer) Serialize(v interface{}) (string, error) {
	b, err := xml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s XMLMsgSerializer) Unserialize(data string, v interface{}) error {
	return xml.Unmarshal([]byte(data), v)
}
