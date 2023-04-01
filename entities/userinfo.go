package entities

type Sender struct {
	To        []byte `json:"to"`
	From      []byte `json:"from"`
	Message   []byte `json:"message"`
	Signature []byte `json:"signature"`
}
