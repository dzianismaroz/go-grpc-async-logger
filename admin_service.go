package main

import (
	"time"
)

func cleanStat() *Stat {
	return &Stat{
		ByMethod:   map[string]uint64{},
		ByConsumer: map[string]uint64{},
	}
}

// --- implementation of Admin Service ---------
func (s *MyMicroService) Logging(none *Nothing, stream Admin_LoggingServer) error {
	logCh := make(chan *Event)
	defer close(logCh)
	s.addLogListenerCh <- logCh
	for {
		select {
		case event := <-logCh:
			if err := stream.Send(event); err != nil {
				return err
			}
		case <-s.ctx.Done():
			return nil
		}
	}
}

func (s *MyMicroService) Statistics(interval *StatInterval, stream Admin_StatisticsServer) error {
	statCh := make(chan *Stat)
	defer close(statCh)
	s.addStatListenerCh <- statCh
	ticker := time.NewTicker(time.Duration(interval.IntervalSeconds) * time.Second)
	response := cleanStat()

	for {
		select {
		case stat := <-statCh:
			appendStat(response, stat)

		case <-ticker.C:
			if err := stream.Send(response); err != nil {
				return err
			}
			response = cleanStat()

		case <-s.ctx.Done():
			return nil
		}
	}
}

func appendStat(dst, src *Stat) {
	for k, v := range src.ByConsumer {
		dst.ByConsumer[k] += v
	}
	for k, v := range src.ByMethod {
		dst.ByMethod[k] += v
	}
}
