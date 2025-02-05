package discord

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vlkhvnn/TestON/internal/models"
	"github.com/vlkhvnn/TestON/internal/store"
)

var guildDefaultLang = make(map[string]string)

// for testing
type Sender interface {
	ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
}

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
	b.HandleMessage(s, m)
}

func (b *Bot) HandleMessage(s Sender, m *discordgo.MessageCreate) {
	// In production, s is a *discordgo.Session with State populated.
	if sess, ok := s.(*discordgo.Session); ok && sess.State != nil {
		if m.Author.ID == sess.State.User.ID {
			return
		}
	}

	guildID := m.GuildID
	if guildID == "" {
		guildID = m.Author.ID
	}

	parts := strings.Fields(m.Content)
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
		limit := 10
		if len(parts) >= 2 {
			if num, err := strconv.Atoi(parts[1]); err == nil {
				limit = num
			} else {
				lang = parts[1]
			}
		}
		if len(parts) >= 3 {
			if num, err := strconv.Atoi(parts[2]); err == nil {
				limit = num
			}
		}
		if limit < 1 {
			limit = 1
		} else if limit > 100 {
			limit = 100
		}
		if lang == "" {
			lang, _ = b.store.Lang.GetUserLang(ctx, guildID)
			if lang == "" {
				lang = "en"
			}
		}
		events, err := b.store.Event.GetRecent(ctx, lang, limit)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No recent changes for language: %s", lang))
				return
			}
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error retrieving recent changes: %v", err))
			return
		}
		b.sendRecentChanges(s, m, lang, events)

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

func (b *Bot) sendRecentChanges(s Sender, m *discordgo.MessageCreate, lang string, events []*models.RecentChangeEvent) {
	header := fmt.Sprintf("Recent changes for '%s':\n", lang)
	var responseBuilder strings.Builder
	responseBuilder.WriteString(header)
	for i, event := range events {
		t := time.Unix(event.Timestamp, 0).Format(time.RFC822)
		encodedTitle := url.PathEscape(event.Title)
		urlStr := fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", lang, encodedTitle)
		entry := fmt.Sprintf("%d. [%s] [%s](%s) by **%s**\nComment: %s\n\n",
			i+1, t, event.Title, urlStr, event.User, event.Comment)
		if responseBuilder.Len()+len(entry) > 2000 {
			s.ChannelMessageSend(m.ChannelID, responseBuilder.String())
			responseBuilder.Reset()
			responseBuilder.WriteString(header)
		}
		responseBuilder.WriteString(entry)
	}
	if responseBuilder.Len() > 0 {
		s.ChannelMessageSend(m.ChannelID, responseBuilder.String())
	}
}
