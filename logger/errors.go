package logger

type GoPlayError struct {
	Err     string
	Details []error
}

func (e GoPlayError) Error() string {
	ers := e.Err
	for _, detail := range e.Details {
		ers += "\n" + detail.Error()
	}
	return ers
}

func GError(err string, details ...error) GoPlayError {
	return GoPlayError{
		Err:     err,
		Details: details,
	}
}
