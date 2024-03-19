package signature

import "encoding/json"

type Signature struct {
	Device     string `json:"device"`
	Data       []byte `json:"data"`
	SignedData []byte `json:"signed_data"`
	Signature  string `json:"signature"`
}

func (s *Signature) MarshalJSON() ([]byte, error) {
	type Alias Signature
	return json.Marshal(&struct {
		Device     string `json:"device"`
		Data       string `json:"data,omitempty"`
		SignedData string `json:"signed_data,omitempty"`
		*Alias
	}{
		Device:     s.Device,
		Data:       string(s.Data),
		SignedData: string(s.SignedData),
		Alias:      (*Alias)(s),
	})
}
