package event_handler

import (
	"log"
	"math"

	"github.com/Mgeorg1/go_randomaliens/gen/pb"
	"github.com/Mgeorg1/go_randomaliens/internal/client/welford"
)

type EventHandler interface {
	Handle(*pb.FrequencyEvent)
}

type LogEventHandler struct{}

func NewLogEventHandler() *LogEventHandler {
	return &LogEventHandler{}
}

func (h *LogEventHandler) Handle(event *pb.FrequencyEvent) {
	log.Printf("Received event: SessionID=%s, Frequency=%f, Timestamp=%s",
		event.SessionId, event.Frequency, event.Timestamp.AsTime().String())
}

type WelfordEventHandler struct {
	welford.Welford
	k float64 //k is the coofficient for the difference between the new value and the current mean,
	//it is used to detect anomalies, if the diff > k * stddev, then the value is considered an anomaly
}

func NewWelfordEventHandler(diffK float64) *WelfordEventHandler {
	return &WelfordEventHandler{k: diffK}
}

func (h *WelfordEventHandler) Handle(event *pb.FrequencyEvent) {
	if h.Count() > 30 {
		diff := math.Abs(event.Frequency - h.Mean())
		if diff > h.k*h.StdDeviation() {
			log.Printf("Anomaly detected: SessionID=%s, Frequency=%f, Timestamp=%s, Mean=%f, Stddev=%f",
				event.SessionId, event.Frequency, event.Timestamp.AsTime().String(), h.Mean(), h.StdDeviation())
			return
		}
	}
	h.Add(event.Frequency)
	if h.Count()%1000 == 0 {
		log.Printf("Processed %d events: SessionID=%s, Frequency=%f, Timestamp=%s, Mean=%f, Stddev=%f",
			h.Count, event.SessionId,
			event.Frequency, event.Timestamp.AsTime().String(), h.Mean(), h.StdDeviation())
	}
}
