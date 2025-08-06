package telegram

import (
	"fmt"
	"strings"

	"github.com/giuxfila/FulmineOrgBot/internal/telegram/intercept"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/lightningtipbot/telebot.v3"
)

func (bot *TipBot) settingHandler(ctx intercept.Context) (intercept.Context, error) {
	m := ctx.Message()
	splits := strings.Split(m.Text, " ")

	if len(splits) == 1 {
		bot.trySendMessage(m.Sender, Translate(ctx, "settingsHelp"), tb.NoPreview)
	} else if len(splits) > 1 {
		switch strings.ToLower(splits[1]) {
		case "unit":
			return bot.addFiatCurrency(ctx)
		case "help":
			return bot.nostrHelpHandler(ctx)
		}
	}
	return ctx, nil
}

func (bot *TipBot) addFiatCurrency(ctx intercept.Context) (intercept.Context, error) {
	m := ctx.Message()
	user, err := GetLnbitsUserWithSettings(m.Sender, *bot)
	if err != nil {
		return ctx, err
	}

	splits := strings.Split(m.Text, " ")
	if len(splits) < 3 {
		currentCurrency := strings.ToUpper(user.Settings.Display.DisplayCurrency)
		if currentCurrency == "" {
			currentCurrency = "BTC"
		}
		msg := fmt.Sprintf(Translate(ctx, "settingsCurrentCurrency"), currentCurrency)
		bot.trySendMessage(m.Sender, msg)
		return ctx, nil
	}

	currencyInput := strings.ToLower(splits[2])
	if currencyInput != "usd" && currencyInput != "eur" && currencyInput != "gbp" && currencyInput != "btc" && currencyInput != "sat" {
		bot.trySendMessage(m.Sender, Translate(ctx, "settingsInvalidCurrency"))
		return ctx, fmt.Errorf("invalid currency")
	}

	if currencyInput == "sat" {
		currencyInput = "BTC"
	}

	user.Settings.Display.DisplayCurrency = currencyInput
	err = UpdateUserRecord(user, *bot)
	if err != nil {
		log.Errorf("[addFiatCurrency] could not update record of user %s: %v", GetUserStr(user.Telegram), err)
		return ctx, err
	}

	bot.trySendMessage(m.Sender, Translate(ctx, "settingsUpdatedCurrency"))
	return ctx, nil
}
