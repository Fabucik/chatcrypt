package main

import (
	"github.com/Fabucik/chatcrypt/client/db"
	"github.com/Fabucik/chatcrypt/client/ui"
)

func main() {
	db.CreateDatabase()
	DB := db.OpenDatabase()
	defer DB.Close()

	ui.CreateApp(DB)

	DB.Close()
}
