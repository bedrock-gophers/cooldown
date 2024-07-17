package main

import (
	"fmt"
	"github.com/bedrock-gophers/cooldown/cooldown"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})
	c := server.DefaultConfig()
	c.Players.SaveData = false

	conf, err := c.Config(log)
	if err != nil {
		log.Fatalln(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()

	for srv.Accept(func(p *player.Player) {
		p.Handle(&handler{
			p:            p,
			coolDownChat: cooldown.NewCoolDown(),
		})
		p.SetGameMode(world.GameModeCreative)
	}) {

	}
}

type handler struct {
	player.NopHandler
	p            *player.Player
	coolDownChat *cooldown.CoolDown
}

func (h *handler) HandleChat(ctx *event.Context, message *string) {
	ctx.Cancel()

	if h.coolDownChat.Active() {
		h.p.Message(text.Colourf("<red>You are sending messages too quickly. Wait for %.2fs</red>", h.coolDownChat.Remaining().Seconds()))
		return
	}

	h.coolDownChat.Set(5 * time.Second)
	_, _ = chat.Global.WriteString(fmt.Sprintf("%s: %s", h.p.Name(), *message))
}
