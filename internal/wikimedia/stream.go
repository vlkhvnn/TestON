package wikimedia

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/r3labs/sse/v2"
	"github.com/vlkhvnn/TestON/internal/models"
	"github.com/vlkhvnn/TestON/internal/store"
	"go.uber.org/zap"
)

var (
	wikiURL = "https://stream.wikimedia.org/v2/stream/recentchange"
)

func StartStream(ctx context.Context, eventStore *store.Storage, logger *zap.SugaredLogger) error {
	client := sse.NewClient(wikiURL)
	errCh := make(chan error, 1)

	go func() {
		err := client.SubscribeRawWithContext(ctx, func(msg *sse.Event) {
			var event models.RecentChangeEvent
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				logger.Errorw("Error unmarshalling event", "error", err)
				return
			}

			if event.Bot {
				return
			}

			parts := strings.Split(event.ServerName, ".")
			if len(parts) < 1 {
				logger.Warnw("Unexpected ServerName format", "serverName", event.ServerName)
				return
			}
			lang := parts[0]

			storageCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := eventStore.Event.Add(storageCtx, lang, &event); err != nil {
				logger.Errorw("Error storing event", "error", err)
				return
			}

			t := time.Unix(event.Timestamp, 0).UTC()
			dateStr := t.Format("2006-01-02")

			if err := eventStore.Stat.IncrementByLang(storageCtx, lang, dateStr); err != nil {
				logger.Errorw("Error updating stats", "error", err)
				return
			}
		})
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}
