package lndhub

import (
	"net/http"

	"github.com/giuxfila/FulmineOrgBot/internal"
	"github.com/giuxfila/FulmineOrgBot/internal/api"
	"github.com/giuxfila/FulmineOrgBot/internal/telegram"
	"gorm.io/gorm"
)

type LndHub struct {
	database *gorm.DB
}

func New(bot *telegram.TipBot) LndHub {
	return LndHub{database: bot.DB.Users}
}
func (w LndHub) Handle(writer http.ResponseWriter, request *http.Request) {
	api.Proxy(writer, request, internal.Configuration.Lnbits.Url)
}
