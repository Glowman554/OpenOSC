package oscmod

import "time"

func TickFPS(fps int, fn func()) {
	frameDuration := time.Second / time.Duration(fps)
	ticker := time.NewTicker(frameDuration)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fn()
			}
		}
	}()
}
