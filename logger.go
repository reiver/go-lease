package lease

type Logger interface{
	Alert(...interface{}) error
	Debug(...interface{})
	Error(...interface{}) error
	Warn(...interface{})
}
