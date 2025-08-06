package telegram

import (
	"context"
	stderrors "errors"
	"fmt"
	"strconv"
	"time"

	"github.com/giuxfila/FulmineOrgBot/internal/telegram/intercept"
	"github.com/giuxfila/FulmineOrgBot/internal/errors"
	"github.com/giuxfila/FulmineOrgBot/internal"
	log "github.com/sirupsen/logrus"
	"github.com/giuxfila/FulmineOrgBot/internal/lnbits"
	"github.com/giuxfila/FulmineOrgBot/internal/str"
	tb "gopkg.in/lightningtipbot/telebot.v3"
	"gorm.io/gorm"
)

func (bot TipBot) startHandler(ctx intercept.Context) (intercept.Context, error) {
	if !ctx.Message().Private() {
		return ctx, errors.Create(errors.NoPrivateChatError)
	}
	if !internal.IsAuthorized(ctx.Sender().ID) {
		return ctx, nil
	}
	
	log.Printf("[💬 /start] New user: %s (%d)\n", GetUserStr(ctx.Sender()), ctx.Sender().ID)
	walletCreationMsg := bot.trySendMessageEditable(ctx.Sender(), Translate(ctx, "startSettingWalletMessage"))
	user, err := bot.initWallet(ctx.Sender())
	if err != nil {
		log.Errorln(fmt.Sprintf("[startHandler] Error with initWallet: %s", err.Error()))
		bot.tryEditMessage(walletCreationMsg, Translate(ctx, "startWalletErrorMessage"))
		return ctx, err
	}
	bot.tryDeleteMessage(walletCreationMsg)
	ctx.Context = context.WithValue(ctx, "user", user)
	bot.helpHandler(ctx)
	bot.trySendMessage(ctx.Sender(), Translate(ctx, "startWalletReadyMessage"))
	bot.balanceHandler(ctx)

	if len(ctx.Sender().Username) == 0 {
		bot.trySendMessage(ctx.Sender(), Translate(ctx, "startNoUsernameMessage"), tb.NoPreview)
	}
	return ctx, nil
}

func (bot TipBot) initWallet(tguser *tb.User) (*lnbits.User, error) {
	user, err := GetUser(tguser, bot)
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		user = &lnbits.User{Telegram: tguser}
		err = bot.createWallet(user)
		if err != nil {
			return user, err
		}
		user.Initialized = true
		err = UpdateUserRecord(user, bot)
		if err != nil {
			log.Errorln(fmt.Sprintf("[initWallet] error updating user: %s", err.Error()))
			return user, err
		}
	} else if !user.Initialized {
//		tipTooltipInitializedHandler(user.Telegram, bot)
		user.Initialized = true
		err = UpdateUserRecord(user, bot)
		if err != nil {
			log.Errorln(fmt.Sprintf("[initWallet] error updating user: %s", err.Error()))
			return user, err
		}
	} else if user.Initialized {
		return user, nil
	} else {
		err = fmt.Errorf("could not initialize wallet")
		return user, err
	}
	return user, nil
}

func (bot TipBot) createWallet(user *lnbits.User) error {
	UserStr := GetUserStr(user.Telegram)
	u, err := bot.Client.CreateUserWithInitialWallet(strconv.FormatInt(user.Telegram.ID, 10),
		fmt.Sprintf("%d (%s)", user.Telegram.ID, UserStr),
		internal.Configuration.Lnbits.AdminId,
		UserStr)
	if err != nil {
		errormsg := fmt.Sprintf("[createWallet] Create wallet error: %s", err.Error())
		log.Errorln(errormsg)
		return err
	}
	user.Wallet = &lnbits.Wallet{}
	user.ID = u.ID
	user.Name = u.Name
	wallet, err := bot.Client.Wallets(*user)
	if err != nil {
		errormsg := fmt.Sprintf("[createWallet] Get wallet error: %s", err.Error())
		log.Errorln(errormsg)
		return err
	}
	user.Wallet = &wallet[0]

	user.AnonID = fmt.Sprint(str.Int32Hash(user.ID))
	user.AnonIDSha256 = str.AnonIdSha256(user)
	user.UUID = str.UUIDSha256(user)

	user.Initialized = false
	user.CreatedAt = time.Now()
	err = UpdateUserRecord(user, bot)
	if err != nil {
		errormsg := fmt.Sprintf("[createWallet] Update user record error: %s", err.Error())
		log.Errorln(errormsg)
		return err
	}
	return nil
}
