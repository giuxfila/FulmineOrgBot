package telegram

import (
	"context"
	"fmt"

        "github.com/giuxfila/FulmineOrgBot/internal"
	"github.com/giuxfila/FulmineOrgBot/internal/telegram/intercept"

	tb "gopkg.in/lightningtipbot/telebot.v3"
)

func (bot TipBot) makeHelpMessage(ctx context.Context, m *tb.Message) string {
	fromUser := LoadUser(ctx)
	dynamicHelpMessage := ""
	if len(m.Sender.Username) == 0 {
		dynamicHelpMessage = dynamicHelpMessage + "\n" + Translate(ctx, "helpNoUsernameMessage")
	}
	lnaddr, _ := bot.UserGetLightningAddress(fromUser)
	if len(lnaddr) > 0 {
		dynamicHelpMessage = dynamicHelpMessage + "\n" + fmt.Sprintf(Translate(ctx, "infoYourLightningAddress"), lnaddr)
	}
	if len(dynamicHelpMessage) > 0 {
		dynamicHelpMessage = Translate(ctx, "infoHelpMessage") + dynamicHelpMessage
	}
	helpMessage := Translate(ctx, "helpMessage")
	return fmt.Sprintf(helpMessage, dynamicHelpMessage)
}

func (bot TipBot) helpHandler(ctx intercept.Context) (intercept.Context, error) {
        if !internal.IsAuthorized(ctx.Sender().ID) {
		return ctx, nil
	}

	bot.anyTextHandler(ctx)
	if !ctx.Message().Private() {
		bot.tryDeleteMessage(ctx.Message())
	}
	bot.trySendMessage(ctx.Sender(), bot.makeHelpMessage(ctx, ctx.Message()), tb.NoPreview)
	return ctx, nil
}

func (bot TipBot) basicsHandler(ctx intercept.Context) (intercept.Context, error) {
	bot.anyTextHandler(ctx)
	if !ctx.Message().Private() {
		bot.tryDeleteMessage(ctx.Message())
	}
	bot.trySendMessage(ctx.Sender(), Translate(ctx, "basicsMessage"), tb.NoPreview)
	return ctx, nil
}

func (bot TipBot) makeAdvancedHelpMessage(ctx context.Context, m *tb.Message) string {
	fromUser := LoadUser(ctx)
	dynamicHelpMessage := "ℹ️ *Info*\n"
	if len(m.Sender.Username) == 0 {
		dynamicHelpMessage = dynamicHelpMessage + fmt.Sprintf("%s", Translate(ctx, "helpNoUsernameMessage")) + "\n"
	}
	lnaddr, err := bot.UserGetAnonLightningAddress(fromUser)
	if err == nil {
		dynamicHelpMessage = dynamicHelpMessage + fmt.Sprintf("Anonymous Lightning address: `%s`\n", lnaddr)
	}
	lnurl, err := UserGetAnonLNURL(fromUser)
	if err == nil {
		dynamicHelpMessage = dynamicHelpMessage + fmt.Sprintf("Anonymous LNURL: `%s`", lnurl)
	}

	return fmt.Sprintf(
		Translate(ctx, "advancedMessage"),
		dynamicHelpMessage,
	)
}

func (bot TipBot) advancedHelpHandler(ctx intercept.Context) (intercept.Context, error) {
	if !internal.IsAuthorized(ctx.Sender().ID) {
		return ctx, nil
	}

	bot.anyTextHandler(ctx)
	if !ctx.Message().Private() {
		bot.tryDeleteMessage(ctx.Message())
	}
	bot.trySendMessage(ctx.Sender(), bot.makeAdvancedHelpMessage(ctx, ctx.Message()), tb.NoPreview)
	return ctx, nil
}
