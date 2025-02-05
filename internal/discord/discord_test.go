package discord

import (
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlkhvnn/TestON/internal/models"
	"github.com/vlkhvnn/TestON/internal/store"
)

type MockSession struct {
	messages []string
}

func (ms *MockSession) ChannelMessageSend(channelID, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	ms.messages = append(ms.messages, content)
	return &discordgo.Message{Content: content}, nil
}

func TestRecentCommandWithLimit(t *testing.T) {
	err := godotenv.Load("../../.env")
	require.NoError(t, err)

	mockEventStore := &store.MockEventStore{
		RecentEvents: []*models.RecentChangeEvent{
			{
				ID:         "1",
				Title:      "Test Page 1",
				User:       "User1",
				Comment:    "Comment 1",
				Timestamp:  time.Now().Unix(),
				Wiki:       "enwiki",
				ServerName: "en.wikipedia.org",
			},
			{
				ID:         "2",
				Title:      "Test Page 2",
				User:       "User2",
				Comment:    "Comment 2",
				Timestamp:  time.Now().Unix(),
				Wiki:       "enwiki",
				ServerName: "en.wikipedia.org",
			},
		},
	}
	mockLangStore := &store.MockLangStore{
		Langs: map[string]string{
			"guild1": "en",
		},
	}
	mockStatStore := &store.MockStatStore{}
	mockStorage := store.Storage{
		Event: mockEventStore,
		Lang:  mockLangStore,
		Stat:  mockStatStore,
	}

	b, err := NewBot("fake-token", mockStorage)
	require.NoError(t, err)

	msgContent := "!recent 5"
	m := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Content:   msgContent,
			ChannelID: "channel1",
			Author: &discordgo.User{
				ID: "user1",
			},
			GuildID: "guild1",
		},
	}

	ms := &MockSession{}

	b.HandleMessage(ms, m)

	assert.NotEmpty(t, ms.messages)

	found := false
	for _, response := range ms.messages {
		if strings.Contains(response, "Recent changes for") {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected response message containing 'Recent changes for'")
}

func TestStatsCommand(t *testing.T) {
	mockStatStore := &store.MockStatStore{
		Stats: map[string]int{
			"en_2025-02-04": 42,
		},
	}
	mockLangStore := &store.MockLangStore{
		Langs: map[string]string{
			"guild1": "en",
		},
	}
	mockEventStore := &store.MockEventStore{}
	mockStorage := store.Storage{
		Event: mockEventStore,
		Lang:  mockLangStore,
		Stat:  mockStatStore,
	}

	b, err := NewBot("fake-token", mockStorage)
	require.NoError(t, err)

	msgContent := "!stats 2025-02-04"
	m := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Content:   msgContent,
			ChannelID: "channel1",
			Author: &discordgo.User{
				ID: "user1",
			},
			GuildID: "guild1",
		},
	}

	ms := &MockSession{}

	b.HandleMessage(ms, m)

	found := false
	for _, response := range ms.messages {
		if strings.Contains(response, "42 changes") {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected response message containing '42 changes'")
}
