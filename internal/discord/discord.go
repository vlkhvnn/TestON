package discord

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vlkhvnn/TestON/internal/store"
)

var guildDefaultLang = make(map[string]string)

type Bot struct {
	session *discordgo.Session
	store   store.Storage
}

func NewBot(token string, storage store.Storage) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	bot := &Bot{
		session: dg,
		store:   storage,
	}

	dg.AddHandler(bot.messageHandler)
	return bot, nil
}

func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return err
	}
	log.Println("Discord bot started.")
	return nil
}

func (b *Bot) Stop() {
	b.session.Close()
}

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	guildID := m.GuildID
	if guildID == "" {
		guildID = m.Author.ID
	}

	content := m.Content
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	ctx := context.Background()

	switch parts[0] {
	case "!setLang":
		if len(parts) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: !setLang [language_code]")
			return
		}
		lang := parts[1]

		err := b.store.Lang.SetUserLang(ctx, guildID, lang)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to set language preference.")
			return
		}

		guildDefaultLang[guildID] = lang
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Default language set to '%s' for this session.", lang))

	case "!recent":
		var lang string
		if len(parts) >= 2 {
			lang = parts[1]
		} else {
			lang, _ = b.store.Lang.GetUserLang(ctx, guildID)
			if lang == "" {
				lang = "en"
			}
		}

		events, err := b.store.Event.GetRecent(ctx, lang, 10)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No recent changes for language: %s", lang))
				return
			}
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error retrieving recent changes: %v", err))
			return
		}

		var responseBuilder strings.Builder
		for i, event := range events {
			t := time.Unix(event.Timestamp, 0).Format(time.RFC822)
			encodedTitle := url.PathEscape(event.Title)
			urlStr := fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", lang, encodedTitle)
			entry := fmt.Sprintf("%d. [%s] [%s](%s) by **%s**\nAuthor Comment: %s\n\n",
				i+1, t, event.Title, urlStr, event.User, event.Comment)

			if responseBuilder.Len()+len(entry) > 2000 {
				s.ChannelMessageSend(m.ChannelID, responseBuilder.String())
				responseBuilder.Reset()
			}

			responseBuilder.WriteString(entry)
		}

		if responseBuilder.Len() > 0 {
			s.ChannelMessageSend(m.ChannelID, responseBuilder.String())
		}

	case "!stats":
		if len(parts) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: !stats [yyyy-mm-dd] [optional: language_code]")
			return
		}

		dateStr := parts[1]
		var lang string
		if len(parts) >= 3 {
			lang = parts[2]
		} else {
			lang, _ = b.store.Lang.GetUserLang(ctx, guildID)
			if lang == "" {
				lang = "en"
			}
		}

		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Invalid date format. Please use yyyy-mm-dd.")
			return
		}

		count, err := b.store.Stat.Get(ctx, lang, dateStr)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No stats found for %s on %s", lang, dateStr))
				return
			}
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error retrieving stats: %v", err))
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("On %s, there were %d changes for language '%s'.", dateStr, count, lang))
	}
}
