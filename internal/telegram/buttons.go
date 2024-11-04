package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/giuxfila/FulmineOrgBot/internal"
	"github.com/giuxfila/FulmineOrgBot/internal/lnbits"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/lightningtipbot/telebot.v3"
)

// Costanti per i comandi dei pulsanti
const (
	MainMenuCommandWebApp  = "🗳️ App"
	MainMenuCommandBalance  = "₿alance"
	MainMenuCommandInvoice  = "⚡️ Invoice"
	MainMenuCommandHelp     = "📖 Help"
	MainMenuCommandSend     = "⤴️"
	SendMenuCommandEnter    = "👤 Enter"
)

var (
	mainMenu          = &tb.ReplyMarkup{ResizeKeyboard: true}
	btnHelpMainMenu   = mainMenu.Text(MainMenuCommandHelp)
	btnWebAppMainMenu = mainMenu.Text(MainMenuCommandWebApp)
	btnSendMainMenu   = mainMenu.Text(MainMenuCommandSend)
	btnBalanceMainMenu = mainMenu.Text(MainMenuCommandBalance)
	btnInvoiceMainMenu = mainMenu.Text(MainMenuCommandInvoice)

	sendToMenu       = &tb.ReplyMarkup{ResizeKeyboard: true}
	sendToButtons    = []tb.Btn{}
	btnSendMenuEnter = mainMenu.Text(SendMenuCommandEnter)
)

// Funzione per controllare se l'ID dell'utente è autorizzato
func isUserAuthorized(userID int64) bool {
	return strconv.FormatInt(userID, 10) == internal.Configuration.Bot.AuthorizedID
}

func init() {
	// Inizializza i pulsanti solo se l'utente è autorizzato
	// Non possiamo verificare direttamente qui l'ID autorizzato senza avere un ID utente
}

// Funzione per impacchettare i pulsanti in righe di lunghezza specificata
func buttonWrapper(buttons []tb.Btn, markup *tb.ReplyMarkup, length int) []tb.Row {
	buttonLength := len(buttons)
	rows := make([]tb.Row, 0)

	if buttonLength > length {
		for i := 0; i < buttonLength; i = i + length {
			buttonRow := make([]tb.Btn, length)
			if i+length < buttonLength {
				buttonRow = buttons[i : i+length]
			} else {
				buttonRow = buttons[i:]
			}
			rows = append(rows, markup.Row(buttonRow...))
		}
		return rows
	}
	rows = append(rows, markup.Row(buttons...))
	return rows
}

// Aggiunge un link WebApp a un pulsante
func (bot *TipBot) appendWebAppLinkToButton(btn *tb.Btn, user *lnbits.User) {
	var url string
	if len(user.Telegram.Username) > 0 {
		url = fmt.Sprintf("%s/app/@%s", internal.Configuration.Bot.LNURLHostName, user.Telegram.Username)
	} else {
		url = fmt.Sprintf("%s/app/@%s", internal.Configuration.Bot.LNURLHostName, user.AnonIDSha256)
	}
	if strings.HasPrefix(url, "https://") {
		btn.WebApp = &tb.WebAppInfo{Url: url}
	}
}

// Aggiorna il pulsante del saldo nel menu principale
func (bot *TipBot) mainMenuBalanceButtonUpdate(to int64) {
	var user *lnbits.User
	var err error
	if user, err = getCachedUser(&tb.User{ID: to}, *bot); err != nil {
		user, err = GetLnbitsUser(&tb.User{ID: to}, *bot)
		if err != nil {
			return
		}
		updateCachedUser(user, *bot)
	}
	if user.Wallet != nil {
		amount, err := bot.GetUserBalanceCached(user)
		if err == nil {
			log.Tracef("[appendMainMenu] user %s balance %d sat", GetUserStr(user.Telegram), amount)
			MainMenuCommandBalance := fmt.Sprintf("%s %d sat", MainMenuCommandBalance, amount)
			btnBalanceMainMenu = mainMenu.Text(MainMenuCommandBalance)
		}

		bot.appendWebAppLinkToButton(&btnWebAppMainMenu, user)
		mainMenu.Reply(
			mainMenu.Row(btnBalanceMainMenu),
			mainMenu.Row(btnInvoiceMainMenu, btnWebAppMainMenu, btnHelpMainMenu),
		)
	}
}

// Crea un array di pulsanti per il menu di invio
func (bot *TipBot) makeContactsButtons(ctx context.Context) []tb.Btn {
	var records []Transaction

	sendToButtons = []tb.Btn{}
	user := LoadUser(ctx)
	bot.DB.Transactions.Where("from_id = ? AND to_user LIKE ? AND to_user <> ?", user.Telegram.ID, "@%", GetUserStr(user.Telegram)).Distinct("to_user").Order("id desc").Limit(5).Find(&records)
	log.Debugf("[makeContactsButtons] found %d records", len(records))

	for i, r := range records {
		log.Tracef("[makeContactsButtons] toNames[%d] = %s (id=%d)", i, r.ToUser, r.ID)
		sendToButtons = append(sendToButtons, tb.Btn{Text: r.ToUser})
	}

	sendToButtons = append(sendToButtons, tb.Btn{Text: SendMenuCommandEnter})
	sendToMenu.Reply(buttonWrapper(sendToButtons, sendToMenu, 3)...)
	return sendToButtons
}

// Aggiunge il menu principale se l'utente è in chat privata
func (bot *TipBot) appendMainMenu(to int64, recipient interface{}, options []interface{}) []interface{} {
	if to > 0 {
		bot.mainMenuBalanceButtonUpdate(to)
	}

	appendKeyboard := true
	for _, option := range options {
		if option == tb.ForceReply {
			appendKeyboard = false
		}
		switch option.(type) {
		case *tb.ReplyMarkup:
			appendKeyboard = false
		}
	}

	if to > 0 && appendKeyboard {
		if isUserAuthorized(to) { // Controlla se l'utente è autorizzato qui
			options = append(options, mainMenu)
		}
	}
	return options
}
