package discord

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vlkhvnn/TestON/internal/store"
)

var guildDefaultLang = make(map[string]string)

// Bot represents the Discord bot.
type Bot struct {
	session    *discordgo.Session
	eventStore *store.Store
}

// NewBot creates a new Discord bot instance.
func NewBot(token string, eventStore *store.Store) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	bot := &Bot{
		session:    dg,
		eventStore: eventStore,
	}
	// Register message handler.
	dg.AddHandler(bot.messageHandler)
	return bot, nil
}

// Start opens the Discord session.
func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return err
	}
	log.Println("Discord bot started.")
	return nil
}

// Stop closes the Discord session.
func (b *Bot) Stop() {
	b.session.Close()
}

// messageHandler processes incoming Discord messages.
func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Use guild ID as key for default language; if DM, use the author ID.
	guildID := m.GuildID
	if guildID == "" {
		guildID = m.Author.ID
	}

	content := m.Content
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "!setLang":
		if len(parts) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: !setLang [language_code]")
			return
		}
		lang := parts[1]
		guildDefaultLang[guildID] = lang
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Default language set to '%s' for this session.", lang))
	case "!recent":
		// Usage: !recent [optional: language_code]
		var lang string
		if len(parts) >= 2 {
			lang = parts[1]
		} else {
			lang = guildDefaultLang[guildID]
			if lang == "" {
				lang = "en"
			}
		}

		events := b.eventStore.GetEvents(lang)
		if len(events) == 0 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No recent changes for language: %s", lang))
			return
		}

		var response strings.Builder
		response.WriteString(fmt.Sprintf("Recent changes for '%s':\n", lang))
		for i, event := range events {
			t := time.Unix(event.Timestamp, 0).Format(time.RFC822)
			// Generate a URL for the change.
			encodedTitle := url.PathEscape(event.Title)
			urlStr := fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", lang, encodedTitle)
			response.WriteString(fmt.Sprintf("%d. [%s] %s (%s) by %s - %s\n", i+1, t, event.Title, urlStr, event.User, event.Comment))
		}
		s.ChannelMessageSend(m.ChannelID, response.String())
	case "!stats":
		// Usage: !stats [yyyy-mm-dd] [optional: language_code]
		if len(parts) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: !stats [yyyy-mm-dd] [optional: language_code]")
			return
		}
		dateStr := parts[1]
		var lang string
		if len(parts) >= 3 {
			lang = parts[2]
		} else {
			lang = guildDefaultLang[guildID]
			if lang == "" {
				lang = "en"
			}
		}
		// Validate date format.
		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Invalid date format. Please use yyyy-mm-dd.")
			return
		}
		count := b.eventStore.GetStats(lang, dateStr)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("On %s, there were %d changes for language '%s'.", dateStr, count, lang))
	}
}
