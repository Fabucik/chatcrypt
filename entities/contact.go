package entities

type Contact struct {
	Name      string
	PublicKey []byte
	ID        int
}

type AllContacts struct {
	Contacts []Contact
}
