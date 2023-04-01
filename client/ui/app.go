package ui

import (
	"database/sql"
	"strings"
	"time"

	"github.com/Fabucik/chatcrypt/client/db"
	"github.com/Fabucik/chatcrypt/cryptography"
	"github.com/Fabucik/chatcrypt/entities"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateApp(DB *sql.DB) {
	var messageToSend string
	var currentContact string
	username := "fanda"

	app := tview.NewApplication()

	loginForm := tview.NewForm()
	menu := tview.NewTextView().SetText("(q) to quit\n(l) to login\n(c) to chat").SetTextColor(tcell.ColorGreen)

	exitButton := tview.NewButton("Quit")
	exitButton.SetBorder(true)
	exitButton.SetSelectedFunc(func() {
		app.Stop()
	})
	exitButton.SetActivatedStyle(tcell.Style{}.Background(loginForm.GetBackgroundColor()))
	exitButton.SetStyle(tcell.Style{}.Background(loginForm.GetBackgroundColor()))
	exitButton.SetLabelColorActivated(tcell.ColorDarkSalmon)

	messageInput := tview.NewTextArea()
	messageInput.SetBorder(true)
	messageInput.SetChangedFunc(func() {
		messageToSend = messageInput.GetText()
	})

	messageGrid := tview.NewGrid()
	messageGrid.SetBorders(false)
	messageGrid.SetTitle("chat")
	messageGrid.SetBorder(true)

	contactList := tview.NewList()
	contactList.SetTitle("contacts")
	contactList.SetBorder(true)
	contactList.SetSelectedFunc(func(index int, mainText string, secondaryText string, id rune) {
		contacts := db.GetContacts(DB)
		for _, contact := range contacts.Contacts {
			if int(id-'0') == contact.ID {
				messages := db.ReadMessages(DB, contact.Name)

				renderMessages(messageGrid, messages)
				currentContact = contact.Name
			}
		}
	})

	messageButton := tview.NewButton("Send")
	messageButton.SetBorder(true)
	messageButton.SetActivatedStyle(tcell.Style{}.Background(loginForm.GetBackgroundColor()))
	messageButton.SetStyle(tcell.Style{}.Background(loginForm.GetBackgroundColor()))
	messageButton.SetLabelColorActivated(tcell.ColorDarkSalmon)
	messageButton.SetSelectedFunc(func() {
		db.AddMessage(DB, []byte(messageToSend), currentContact, username, time.Now().Unix())
		messages := db.ReadMessages(DB, currentContact)
		renderMessages(messageGrid, messages)
	})

	chat := tview.NewFlex()
	chat.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).AddItem(contactList, 0, 8, false).
		AddItem(exitButton, 0, 1, false), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).AddItem(
			messageGrid, 0, 8, false).
			AddItem(tview.NewFlex().AddItem(messageInput, 0, 7, false).AddItem(messageButton, 0, 1, false), 0, 1, false), 0, 8, false)

	pages := tview.NewPages()
	pages.AddAndSwitchToPage("Menu", menu, true)
	pages.AddPage("Login", loginForm, true, false)
	pages.AddPage("Chat", chat, true, false)

	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case rune('q'):
			app.Stop()

		case rune('l'):
			addLoginForm(loginForm, pages)
			pages.SwitchToPage("Login")

		case rune('c'):
			contacts := db.GetContacts(DB)
			renderContacts(contactList, contacts)
			pages.SwitchToPage("Chat")
		}
		return event
	})

	err := app.SetRoot(pages, true).EnableMouse(true).Run()
	if err != nil {
		panic(err)
	}
}

func addLoginForm(form *tview.Form, pages *tview.Pages) {
	var login entities.LoginInfo

	form.AddInputField("Path to key dir", "", 50, nil, func(keysPath string) {
		login.KeysPath = keysPath
	})

	form.AddButton("Submit", func() {
		_, _, err := cryptography.ReadKeysFromFiles(login.KeysPath)
		if err != nil {
			form.Clear(false)

			form.AddInputField("Path to key dir", "", 50, nil, func(keysPath string) {
				login.KeysPath = keysPath
			})

			form.AddTextView("Error", "Failed reading keys", 25, 1, false, false)
		}
	})

	form.AddButton("Cancel", func() {
		pages.SwitchToPage("Menu")
	})
}

func renderMessages(grid *tview.Grid, messages entities.AllMessages) {
	grid.Clear()

	rowSizes := make([]int, len(messages.Messages))

	for index, message := range messages.Messages {
		messageArea := tview.NewTextArea().SetText(message.By+": "+string(message.Message), false)
		messageArea.SetDisabled(true)

		numOfNewLines := strings.Split(string(message.Message), "\n")
		rowSizes[index] = len(numOfNewLines)

		grid.AddItem(messageArea, index, 0, 1, 1, 0, 0, false)
	}

	grid.SetRows(rowSizes...)
}

func renderContacts(contactList *tview.List, contacts entities.AllContacts) {
	for _, contact := range contacts.Contacts {
		contactList.AddItem(contact.Name, "", rune(49+contact.ID-1), nil)
	}
}
