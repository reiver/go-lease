package lease

// An error that fits Timedout may be returned from lease.Type.Lease()
// if it times out trying to lock or unlock.
//
//	import "github.com/reiver/go-lease"
//	
//	// ...
//	
//	var tenant lease.Type
//	
//	// ...
//	
//	err := tenant.Lease(fn)
//	
//	switch casted := err.(type) {
//	case lease.Timedout:
//		//...
//	default:
//		//...
//	}
type Timedout interface {
	error
	Timedout()
}

type internalTimedout struct {
	msg string
}

var _ Timedout = internalTimedout{}

func (receiver internalTimedout) Error() string {
	switch receiver.msg {
	case "":
		return "timedout out"
	default:
		return receiver.msg
	}
}

func (internalTimedout) Timedout() {
	// Nothing here.
}
