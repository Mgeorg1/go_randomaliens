package event_handler

import (
	"context"
	"log"
	"math"
	"sync"
	"time"

	"github.com/Mgeorg1/go_randomaliens/gen/pb"
	"github.com/Mgeorg1/go_randomaliens/internal/client/repo"
	"github.com/Mgeorg1/go_randomaliens/internal/client/welford"
)

type EventHandler interface {
	Handle(*pb.FrequencyEvent) bool //returns true if the event is an anomaly, false otherwise
}

type LogEventHandler struct{}

func NewLogEventHandler() *LogEventHandler {
	return &LogEventHandler{}
}

func (h *LogEventHandler) Handle(event *pb.FrequencyEvent) bool {
	log.Printf("Received event: SessionID=%s, Frequency=%f, Timestamp=%s",
		event.SessionId, event.Frequency, event.Timestamp.AsTime().String())
	return false
}

type WelfordEventHandler struct {
	welford.Welford
	k float64 //k is the coofficient for the difference between the new value and the current mean,
	//it is used to detect anomalies, if the diff > k * stddev, then the value is considered an anomaly
	ticker *time.Ticker
	repo   *repo.Repo
}

func NewWelfordEventHandler(diffK float64, repo *repo.Repo) *WelfordEventHandler {
	return &WelfordEventHandler{k: diffK, ticker: time.NewTicker(10 * time.Second), repo: repo}
}

func (h *WelfordEventHandler) Stop() {
	h.ticker.Stop()
}

func (h *WelfordEventHandler) Handle(event *pb.FrequencyEvent) bool {
	if h.Count() > 30 {
		diff := math.Abs(event.Frequency - h.Mean())
		if diff > h.k*h.StdDeviation() {
			modelEvent := &repo.FrequencyEvent{}
			modelEvent.FromPbEvent(event)
			if err := h.repo.SaveFrequencyEvent(modelEvent); err != nil {
				log.Printf("Error saving anomaly event to database: %v", err)
			}
			return true
		}
	}
	h.Add(event.Frequency)

	select {
	case <-h.ticker.C:
		log.Printf("Processed %d events: SessionID=%s, Frequency=%f, Timestamp=%s, Mean=%f, StdDeviation=%f",
			h.Count(), event.SessionId,
			event.Frequency, event.Timestamp.AsTime().String(), h.Mean(), h.StdDeviation())
	default:
	}
	return false
}

func HandlerWorker(wg *sync.WaitGroup, ctx context.Context, events <-chan *pb.FrequencyEvent, handle func(*pb.FrequencyEvent) bool) {
	defer wg.Done()
	for {
		select {
		case event, ok := <-events:
			if !ok {
				return
			}
			handle(event)
		case <-ctx.Done():
			for event := range events {
				handle(event) //handle remaining events before exiting
			}
			return
		}
	}
}
