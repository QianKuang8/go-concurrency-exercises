//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mu        sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	stopCh := make(chan struct{}, 1)
	interval := time.Tick(time.Second)
	go func() {
		u.mu.Lock()
		canUse := u.IsPremium || u.TimeUsed < 10
		u.mu.Unlock()
		if !canUse {
			return
		}
		process()
		stopCh <- struct{}{}
		close(stopCh)
	}()
	for {
		select {
		case <-stopCh:
			return true
		case <-interval:
			u.mu.Lock()
			u.TimeUsed++
			canUse := u.IsPremium || u.TimeUsed < 10
			u.mu.Unlock()
			if !canUse {
				return false
			}
		}
	}
}

func main() {
	RunMockServer()
}
