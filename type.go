package lease

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"
)

var defaultTimeout = time.Second * time.Duration(7)

type Type struct {
	state int32
	Logger Logger
}

func (receiver Type) logger() Logger {
	log := receiver.Logger

	if nil == log {
		return internalDiscardLogger{}
	}

	return log
}

func (receiver *Type) lock(traceid int64) error {
	started := time.Now()

	log := receiver.logger()
	log.Debug("lock: BEGIN ⧼",traceid,"⧽")

	if nil == receiver {
		err := errNilReceiver
		log.Error("lock: ", err, " ⧼",traceid,"⧽")
		log.Debug("lock: END NIL RECEIVER ⧼",traceid,"⧽")
		return err
	}

	for !atomic.CompareAndSwapInt32(&receiver.state, unlocked, locked) {
		runtime.Gosched()
	}

	log.Debug("lock: END ", time.Now().Sub(started), " ⧼",traceid,"⧽")
	return nil
}

func (receiver *Type) locktry(traceid int64, timeout time.Duration) error {
	started := time.Now()

	log := receiver.logger()
	log.Debug("locktry: BEGIN ⧼",traceid,"⧽")

	if nil == receiver {
		err := errNilReceiver
		log.Error("locktry: ", err, " ⧼",traceid,"⧽")
		log.Debug("locktry: END NIL RECEIVER ⧼",traceid,"⧽")
		return err
	}

	{
		start := time.Now()

		for !atomic.CompareAndSwapInt32(&receiver.state, unlocked, locked) {
			runtime.Gosched()

			diff := time.Now().Sub(start)

			if timeout < diff  {
				var err error = internalTimedout{fmt.Sprintf("locking timed out after %s", diff)}
				log.Error("locktry: ", err, " ⧼",traceid,"⧽")
				log.Debug("locktry: END ERROR ", diff, " ⧼",traceid,"⧽")
				return err
			}
		}
	}

	log.Debug("locktry: END ", time.Now().Sub(started), " ⧼",traceid,"⧽")
	return nil
}

func (receiver *Type) unlock(traceid int64) error {
	started := time.Now()

	log := receiver.logger()
	log.Debug("unlock: BEGIN ⧼",traceid,"⧽")

	if nil == receiver {
		err := errNilReceiver
		log.Error("unlock: ", err, " ⧼",traceid,"⧽")
		log.Debug("unlock: END NIL RECEIVER ⧼",traceid,"⧽")
		return err
	}

	{
		start := time.Now()

		for !atomic.CompareAndSwapInt32(&receiver.state, locked, unlocked) {
			runtime.Gosched()

			diff := time.Now().Sub(start)

			if defaultTimeout < diff  {
				var err error = internalTimedout{fmt.Sprintf("unlocking timed out after %s", diff)}
				log.Error("unlock: ", err, " ⧼",traceid,"⧽")
				log.Debug("unlock: END ERROR ", diff, " ⧼",traceid,"⧽")
				return err
			}
		}
	}

	log.Debug("unlock: END ", time.Now().Sub(started), " ⧼",traceid,"⧽")
	return nil
}

func (receiver *Type) Lease(fn func()) (err error) {
	started := time.Now()
	traceid := rand.Int63()

	log := receiver.logger()
	log.Debug("lease: BEGIN ⧼",traceid,"⧽")

	if nil == receiver {
		err := errNilReceiver
		log.Error("lease: ", err, " ⧼",traceid,"⧽")
		log.Debug("lease: END NIL RECEIVER ⧼",traceid,"⧽")
		return err
	}

	wait := defaultTimeout

	err = receiver.locktry(traceid, wait)
	if nil != err {
		log.Error("lease: ", err, " ⧼",traceid,"⧽")
		log.Debug("lease: END ERROR ", time.Now().Sub(started), " ⧼",traceid,"⧽")
		return err
	}
	defer func(){
		err = receiver.unlock(traceid)
		if nil != err {
			log.Error("lease: ", err, " ⧼",traceid,"⧽")
			log.Debug("lease: END ERROR ", time.Now().Sub(started), " ⧼",traceid,"⧽")
		} else {
			log.Debug("lease: END ", time.Now().Sub(started), " ⧼",traceid,"⧽")
		}
	}()

	fn()

	return
}
