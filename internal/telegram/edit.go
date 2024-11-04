package telegram

import (
	"strings"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/lightningtipbot/telebot.v3"
)

var editStack cmap.ConcurrentMap

type edit struct {
	to       tb.Editable
	key      string
	what     interface{}
	options  []interface{}
	lastEdit time.Time
	edited   bool
}

func init() {
	editStack = cmap.New()
}

const resultTrueError = "telebot: result is True"
const editSameStringError = "specified new message content and reply markup are exactly the same as a current content and reply markup of the message"
const retryAfterError = "retry after"

// startEditWorker will loop through the editStack and run tryEditMessage on not edited messages.
// if editFromStack is older than 5 seconds, editFromStack will be removed.
func (bot TipBot) startEditWorker() {
	go func() {
		for {
			for _, k := range editStack.Keys() {
				if e, ok := editStack.Get(k); ok {
					editFromStack := e.(edit)
					if !editFromStack.edited {
						_, err := bot.tryEditMessage(editFromStack.to, editFromStack.what, editFromStack.options...)
						if err != nil && strings.Contains(err.Error(), retryAfterError) {
							// ignore any other error than retry after
							log.Errorf("[startEditWorker] Edit error: %s. len(editStack)=%d", err.Error(), len(editStack.Keys()))

						} else {
							if err != nil {
								log.Errorf("[startEditWorker] Ignoring edit error: %s. len(editStack)=%d", err.Error(), len(editStack.Keys()))
							}
							log.Tracef("[startEditWorker] message from stack edited %+v. len(editStack)=%d", editFromStack, len(editStack.Keys()))
							editFromStack.lastEdit = time.Now()
							editFromStack.edited = true
							editStack.Set(k, editFromStack)
						}
					} else {
						if editFromStack.lastEdit.Before(time.Now().Add(-(time.Duration(5) * time.Second))) {
							log.Tracef("[startEditWorker] removing message edit from stack %+v. len(editStack)=%d", editFromStack, len(editStack.Keys()))
							editStack.Remove(k)
						}
					}
				}
			}
			time.Sleep(time.Millisecond * 1000)
		}
	}()

}

// tryEditStack will add the editable to the edit stack, if what (message) changed.
func (bot TipBot) tryEditStack(to tb.Editable, key string, what interface{}, options ...interface{}) {
	sig, chat := to.MessageSig()
	log.Tracef("[tryEditStack] sig=%s, chat=%d, key=%s, what=%+v, options=%+v", sig, chat, key, what, options)
	// var sig = fmt.Sprintf("%s-%d", msgSig, chat)
	if e, ok := editStack.Get(key); ok {
		editFromStack := e.(edit)
		if editFromStack.what == what.(string) {
			log.Tracef("[tryEditStack] Message already in edit stack. Skipping")
			return
		}
	}
	e := edit{options: options, key: key, what: what, to: to}

	editStack.Set(key, e)
	log.Tracef("[tryEditStack] Added message %s to edit stack. len(editStack)=%d", key, len(editStack.Keys()))
}
