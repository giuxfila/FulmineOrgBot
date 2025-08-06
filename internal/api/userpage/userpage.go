package userpage

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/giuxfila/FulmineOrgBot/internal"
	"github.com/giuxfila/FulmineOrgBot/internal/telegram"
	"github.com/PuerkitoBio/goquery"
	"github.com/fiatjaf/go-lnurl"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	bot *telegram.TipBot
}

func New(b *telegram.TipBot) Service {
	return Service{
		bot: b,
	}
}

const botImage = "https://raw.githubusercontent.com/giuxfila/FulmineOrgBot/1248e4fe19df37e0f811a67c706ff9279b25d8c4/resources/fulmineorgbot_logo.png"

//go:embed static/userpage.html static/webapp.html
var templates embed.FS
var userpage_tmpl = template.Must(template.ParseFS(templates, "static/userpage.html"))
var qr_tmpl = template.Must(template.ParseFS(templates, "static/webapp.html"))

var Client = &http.Client{
	Timeout: 10 * time.Second,
}

// thank you fiatjaf for this code
func (s Service) getTelegramUserPictureURL(username string) (string, error) {
	// with proxy:
	// client, err := s.network.GetHttpClient()
	// if err != nil {
	// 	return "", err
	// }
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://t.me/" + username)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	url, ok := doc.Find(`meta[property="og:image"]`).First().Attr("content")
	if !ok {
		return "", errors.New("no image available for this user")
	}

	return url, nil
}

func (s Service) UserPageHandler(w http.ResponseWriter, r *http.Request) {
	// https://fulmine.org/@<username>
	username := strings.ToLower(mux.Vars(r)["username"])

	_, exists := s.bot.UserExistsByTelegramUsername(username)
	if !exists {
		http.Error(w, "⚠️ L'utente non è autorizzato ad utilizzare questo wallet", http.StatusNotFound)
		return
	}

	callback := fmt.Sprintf("%s/.well-known/lnurlp/%s", internal.Configuration.Bot.LNURLHostName, username)
	botName := internal.Configuration.Bot.Name
	botUsername := internal.Configuration.Bot.Username
	log.Infof("[UserPage] rendering page of %s", username)

	lnurlEncode, err := lnurl.LNURLEncode(callback)
	if err != nil {
		log.Errorln("[UserPage]", err)
		return
	}

	image, err := s.getTelegramUserPictureURL(username)
	if err != nil || image == "https://telegram.org/img/t_logo.png" {
		image = botImage
	}

	if err := userpage_tmpl.ExecuteTemplate(w, "userpage", struct {
		Username    string
		Image       string
		LNURLPay    string
		BotUsername string
		BotName     string
	}{username, image, lnurlEncode, botUsername, botName}); err != nil {
		log.Errorf("failed to render template")
	}
}


func (s Service) UserWebAppHandler(w http.ResponseWriter, r *http.Request) {
	// https://fulmine.org/app/<username>
	username := strings.ToLower(mux.Vars(r)["username"])
	callback := fmt.Sprintf("%s/.well-known/lnurlp/%s", internal.Configuration.Bot.LNURLHostName, username)
	botName := internal.Configuration.Bot.Name
	botUsername := internal.Configuration.Bot.Username
	log.Infof("[UserPage] rendering webapp of %s", username)
	lnurlEncode, err := lnurl.LNURLEncode(callback)
	if err != nil {
		log.Errorln("[UserPage]", err)
		return
	}
	if err := qr_tmpl.ExecuteTemplate(w, "webapp", struct {
		Username    string
		LNURLPay    string
		Callback    string
		BotUsername string
		BotName     string
	}{username, lnurlEncode, callback, botUsername, botName}); err != nil {
		log.Errorf("failed to render template")
	}
}
