package entities

type MessageInfo struct {
	Message   []byte
	By        string
	Timestamp string
	ID        int
}

type AllMessages struct {
	Messages []MessageInfo
}
