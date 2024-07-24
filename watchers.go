package main

// Manage subscribers, cancelation and broadcasting of logs and statistics.
func (s *MyMicroService) startBroadcasting() {
	go func() {
		for {
			select {
			case ch := <-s.addLogListenerCh: // Submit new broadcast log-listener
				s.logListeners[ch] = struct{}{}
			case ch := <-s.removeLogListenerCh: // Unsubscribe existing log-listener
				close(ch)
				delete(s.logListeners, ch)
			case event := <-s.broadcastLogCh:
				for ch := range s.logListeners {
					ch <- event
				}
			case <-s.ctx.Done():
				return
			default:
			}
		}
	}()

	go func() {
		for {
			select {
			case ch := <-s.addStatListenerCh:
				s.statListeners[ch] = struct{}{}
			case ch := <-s.removeStatListenerCh:
				close(ch)
				delete(s.statListeners, ch)
			case stat := <-s.broadcastStatCh:
				for ch := range s.statListeners {
					ch <- stat
				}
			case <-s.ctx.Done():
				return
			default: // Linter is wrong. Default is very important.
			}
		}
	}()
}
