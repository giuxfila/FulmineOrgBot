package telegram

import (
	"github.com/giuxfila/FulmineOrgBot/internal/lnbits"
)

func (bot *TipBot) UserExistsByTelegramUsername(username string) (*lnbits.User, bool) {
	user, err := GetUserByTelegramUsername(username, *bot)
	if err != nil || user == nil || user.Telegram.ID == 0 {
		return nil, false
	}
	return user, true
}
