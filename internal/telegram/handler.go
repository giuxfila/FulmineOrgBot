package telegram

import (
	"fmt"
	"strings"

	"github.com/giuxfila/FulmineOrgBot/internal/telegram/intercept"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/lightningtipbot/telebot.v3"
)

type InterceptionWrapper struct {
	Endpoints   []interface{}
	Handler     intercept.Func
	Interceptor *Interceptor
}

func (bot TipBot) registerTelegramHandlers() {
	telegramHandlerRegistration.Do(func() {
		// Set up handlers
		for _, h := range bot.getHandler() {
			log.Traceln("registering", h.Endpoints)
			bot.register(h)
		}

	})
}

func getDefaultBeforeInterceptor(bot TipBot) []intercept.Func {
	return []intercept.Func{bot.idInterceptor}
}
func getDefaultDeferInterceptor(bot TipBot) []intercept.Func {
	return []intercept.Func{bot.unlockInterceptor}
}
func getDefaultAfterInterceptor(bot TipBot) []intercept.Func {
	return []intercept.Func{}
}

// registerHandlerWithInterceptor will register a ctx with all the predefined interceptors, based on the interceptor type
func (bot TipBot) registerHandlerWithInterceptor(h InterceptionWrapper) {
	h.Interceptor.Before = append(getDefaultBeforeInterceptor(bot), h.Interceptor.Before...)
	//h.Interceptor.After = append(h.Interceptor.After, getDefaultAfterInterceptor(bot)...)
	//h.Interceptor.OnDefer = append(h.Interceptor.OnDefer, getDefaultDeferInterceptor(bot)...)
	for _, endpoint := range h.Endpoints {
		bot.handle(endpoint, intercept.WithHandler(h.Handler,
			intercept.WithBefore(h.Interceptor.Before...),
			intercept.WithAfter(h.Interceptor.After...),
			intercept.WithDefer(h.Interceptor.OnDefer...)))
	}
}

// handle accepts an endpoint and handler for Telegram handler registration.
// function will automatically register string handlers as uppercase and first letter uppercase.
func (bot TipBot) handle(endpoint interface{}, handler tb.HandlerFunc) {
	// register the endpoint
	bot.Telegram.Handle(endpoint, handler)
	switch endpoint.(type) {
	case string:
		// check if this is a string endpoint
		sEndpoint := endpoint.(string)
		if strings.HasPrefix(sEndpoint, "/") {
			// Uppercase endpoint registration, because starting with slash
			bot.Telegram.Handle(strings.ToUpper(sEndpoint), handler)
			if len(sEndpoint) > 2 {
				// Also register endpoint with first letter uppercase
				bot.Telegram.Handle(fmt.Sprintf("/%s%s", strings.ToUpper(string(sEndpoint[1])), sEndpoint[2:]), handler)
			}
		}
	}
}

// register registers a handler, so that Telegram can handle the endpoint correctly.
func (bot TipBot) register(h InterceptionWrapper) {
	if h.Interceptor != nil {
		bot.registerHandlerWithInterceptor(h)
	} else {
		for _, endpoint := range h.Endpoints {
			bot.handle(endpoint, intercept.WithHandler(h.Handler))
		}
	}
}

// getHandler returns a list of all handlers, that need to be registered with Telegram
func (bot TipBot) getHandler() []InterceptionWrapper {
	return []InterceptionWrapper{
		{
			Endpoints: []interface{}{"/start"},
			Handler:   bot.startHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.loadUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/pay"},
			Handler:   bot.payHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/invoice", &btnInvoiceMainMenu},
			Handler:   bot.invoiceHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/set"},
			Handler:   bot.settingHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/nostr"},
			Handler:   bot.nostrHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/node"},
			Handler:   bot.nodeHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{&btnSatdressCheckInvoice},
			Handler:   bot.satdressCheckInvoiceHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.requireUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/balance", &btnBalanceMainMenu},
			Handler:   bot.balanceHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/send", &btnSendMenuEnter},
			Handler:   bot.sendHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.loadReplyToInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{&btnSendMainMenu},
			Handler:   bot.keyboardSendHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.loadReplyToInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/transactions"},
			Handler:   bot.transactionsHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
				}},
		},
		{
			Endpoints: []interface{}{&btnLeftTransactionsButton},
			Handler:   bot.transactionsScrollLeftHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.loadUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{&btnRightTransactionsButton},
			Handler:   bot.transactionsScrollRightHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.loadUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/help", &btnHelpMainMenu},
			Handler:   bot.helpHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.loadUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/basics"},
			Handler:   bot.basicsHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.loadUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/donate"},
			Handler:   bot.donationHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/advanced"},
			Handler:   bot.advancedHelpHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/link"},
			Handler:   bot.lndhubHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/api"},
			Handler:   bot.apiHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{"/lnurl"},
			Handler:   bot.lnurlHandler,
			Interceptor: &Interceptor{
				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{tb.OnPhoto},
			Handler:   bot.photoHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.requireUserInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{tb.OnDocument, tb.OnVideo, tb.OnAnimation, tb.OnVoice, tb.OnAudio, tb.OnSticker, tb.OnVideoNote},
			Handler:   bot.fileHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor,
					bot.logMessageInterceptor,
					bot.loadUserInterceptor}},
		},
		{
			Endpoints: []interface{}{tb.OnText},
			Handler:   bot.anyTextHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.requirePrivateChatInterceptor, // Respond to any text only in private chat
					bot.localizerInterceptor,
					bot.logMessageInterceptor,
					bot.loadUserInterceptor, // need to use loadUserInterceptor instead of requireUserInterceptor, because user might not be registered yet
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
//		{
//			Endpoints: []interface{}{tb.OnQuery},
//			Handler:   bot.anyQueryHandler,
//			Interceptor: &Interceptor{
//				Before: []intercept.Func{
//					bot.localizerInterceptor,
//					bot.requireUserInterceptor,
//					bot.lockInterceptor,
//				},
//				OnDefer: []intercept.Func{
//					bot.unlockInterceptor,
//				},
//			},
//		},
//		{
//			Endpoints: []interface{}{tb.OnInlineResult},
//			Handler:   bot.anyChosenInlineHandler,
//		},
		{
			Endpoints: []interface{}{&btnPay},
			Handler:   bot.confirmPayHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.requireUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{&btnCancelPay},
			Handler:   bot.cancelPaymentHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.requireUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{&btnSend},
			Handler:   bot.confirmSendHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.requireUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{&btnCancelSend},
			Handler:   bot.cancelSendHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.requireUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},,
		{
			Endpoints: []interface{}{&btnWithdraw},
			Handler:   bot.confirmWithdrawHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.requireUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{&btnCancelWithdraw},
			Handler:   bot.cancelWithdrawHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.requireUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{&btnAuth},
			Handler:   bot.confirmLnurlAuthHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.requireUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
		{
			Endpoints: []interface{}{&btnCancelAuth},
			Handler:   bot.cancelLnurlAuthHandler,
			Interceptor: &Interceptor{

				Before: []intercept.Func{
					bot.localizerInterceptor,
					bot.requireUserInterceptor,
					bot.answerCallbackInterceptor,
					bot.lockInterceptor,
				},
				OnDefer: []intercept.Func{
					bot.unlockInterceptor,
				},
			},
		},
	}
}
