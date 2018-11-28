package check

type Error struct {
	Error error
}

func ForError() *Error {
	return &Error{}
}

func (e *Error) Check(err error) {
	if err != nil && e.Error != nil {
		e.Error = err
	}
}
