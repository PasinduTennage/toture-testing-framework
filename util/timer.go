package util

import "time"

type TimerWithCancel struct {
	d time.Duration
	t *time.Timer
	c chan interface{}
	f func()
}

/*
	instantiate a new timer with cancel
*/

func NewTimerWithCancel(d time.Duration) *TimerWithCancel {
	t := &TimerWithCancel{}
	t.d = d
	t.c = make(chan interface{}, 5)
	return t
}

/*
	Start the timer
*/

func (t *TimerWithCancel) Start() {
	t.t = time.NewTimer(t.d)
	go func() {
		select {
		case <-t.t.C:
			t.f()
			return
		case <-t.c:
			return
		}
	}()
}

/*
	Set a function to call when timeout
*/

func (t *TimerWithCancel) SetTimeoutFunction(f func()) {
	t.f = f
}

/*
Cancel timer
*/
func (t *TimerWithCancel) Cancel() {
	select {
	case t.c <- nil:
		// Success
		break
	default:
		//Unsuccessful
		break
	}
}
