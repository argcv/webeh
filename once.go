package webeh

import (
	"runtime"
	"sync/atomic"
)

// This is just another implementation of Once
// According to some benchmark, their performance are nearly the same
// in this case, it just added a new function "reset"

// STATUS:
// 0: not executed yet
// 1: in progressing...
// 2: done
type Once struct {
	status uint32
}

// call once
func (o *Once) Do(f func()) {
	if o.status != 2 {
		if atomic.CompareAndSwapUint32(&(o.status), 0, 1) {
			f()
			atomic.StoreUint32(&(o.status), 2)
		} else {
			for o.status != 2 {
				runtime.Gosched()
				c := atomic.LoadUint32(&(o.status))
				switch c {
				case 2:
					// successful
					return
				case 0:
					// this means it was resetted
					// current status is expired
					// re-initialize it
					o.Do(f)
					return
				default:
					// 1 continue waiting
				}
			}
		}
	}
}

// call once again in the future
func (o *Once) Reset() {
	atomic.StoreUint32(&(o.status), 0)
}
