package main

import (
	"time"
)

// Prepare clean actual statistics after it was consumed already by subscriber.
func cleanStat() *Stat {
	return &Stat{
		ByMethod:   map[string]uint64{},
		ByConsumer: map[string]uint64{},
	}
}

// --- implementation of Admin Service ---------

func (s *MyMicroService) Logging(none *Nothing, stream Admin_LoggingServer) error {
	logCh := make(chan *Event)
	s.addLogListenerCh <- logCh // Subscribe new logging consumer.
	defer func() {
		s.removeLogListenerCh <- logCh //Unsubsribe
	}()
	for {
		select {
		case event := <-logCh: // Stream all loggging events.
			if err := stream.Send(event); err != nil {
				return err
			}
		case <-s.ctx.Done(): //Canceled.
			return nil
		default:
		}
	}
}

func (s *MyMicroService) Statistics(interval *StatInterval, stream Admin_StatisticsServer) error {
	statCh := make(chan *Stat)
	s.addStatListenerCh <- statCh //Subscribe new statistics consumer.
	defer func() {
		s.removeStatListenerCh <- statCh //Unsubscribe.
	}()
	ticker := time.NewTicker(time.Duration(interval.IntervalSeconds) * time.Second)
	response := cleanStat()

	for {
		select {
		case stat := <-statCh: // Collect actual statistics eventually.
			appendStat(response, stat)

		case <-ticker.C: //Send statistics to subscriber on specified interval.
			if err := stream.Send(response); err != nil {
				return err
			}
			response = cleanStat()

		case <-s.ctx.Done(): //Cancelation.
			return nil
		default:
		}
	}
}

// Collect statistics from broadcaster.
func appendStat(dst, src *Stat) {
	for k, v := range src.ByConsumer {
		dst.ByConsumer[k] += v
	}
	for k, v := range src.ByMethod {
		dst.ByMethod[k] += v
	}
}
