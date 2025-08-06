package telegram

import (
	"github.com/giuxfila/FulmineOrgBot/internal/lnbits"
	"github.com/giuxfila/FulmineOrgBot/internal/telegram/intercept"
)

type StateCallbackMessage map[lnbits.UserStateKey]func(ctx intercept.Context) (intercept.Context, error)

var stateCallbackMessage StateCallbackMessage

func initializeStateCallbackMessage(bot *TipBot) {
	stateCallbackMessage = StateCallbackMessage{
		lnbits.UserStateLNURLEnterAmount:     bot.enterAmountHandler,
		lnbits.UserEnterAmount:               bot.enterAmountHandler,
		lnbits.UserEnterUser:                 bot.enterUserHandler,
	}
}
