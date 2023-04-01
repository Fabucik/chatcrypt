package db

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/Fabucik/chatcrypt/entities"
	_ "github.com/mattn/go-sqlite3"
)

func CreateDatabase() {
	os.Create("./database/chatcrypt.db")
}

func OpenDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/chatcrypt.db")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func CreateContactTable(db *sql.DB) {
	contactTable := `CREATE TABLE contacts (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Name" TEXT,
		"PublicKey" BLOB);`

	query, err := db.Prepare(contactTable)
	if err != nil {
		log.Fatalln(err)
	}

	query.Exec()
}

func AddContact(db *sql.DB, name string, publicKey []byte) {
	records := `INSERT INTO contacts(Name, Publickey) VALUES(?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = query.Exec(name, publicKey)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetContacts(db *sql.DB) entities.AllContacts {
	records := `SELECT * FROM contacts`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := query.Query()
	if err != nil {
		log.Fatalln(err)
	}

	var contact entities.Contact
	var allContacts entities.AllContacts

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&contact.ID, &contact.Name, &contact.PublicKey)
		allContacts.Contacts = append(allContacts.Contacts, contact)
	}

	return allContacts
}

func CreateMessageTable(db *sql.DB) {
	messageTable := `CREATE TABLE messages (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Message" BLOB,
		"Contact" TEXT,
		"By" TEXT,
		"Timestamp" TEXT);`

	query, err := db.Prepare(messageTable)
	if err != nil {
		log.Fatalln(err)
	}

	query.Exec()
}

func AddMessage(db *sql.DB, message []byte, contact string, by string, timestamp int64) {
	records := `INSERT INTO messages(Message, Contact, By, Timestamp) VALUES(?, ?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = query.Exec(message, contact, by, strconv.FormatInt(timestamp, 10))
	if err != nil {
		log.Fatalln(err)
	}
}

func ReadMessages(db *sql.DB, contact string) entities.AllMessages {
	records := `SELECT id, Message, By, Timestamp FROM messages WHERE Contact = ?`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatalln((err))
	}

	rows, err := query.Query(contact)
	if err != nil {
		log.Fatalln(err)
	}

	var message entities.MessageInfo
	var allMessages entities.AllMessages

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&message.ID, &message.Message, &message.By, &message.Timestamp)
		allMessages.Messages = append(allMessages.Messages, message)
	}

	return allMessages
}
