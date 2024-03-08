package logger

type RetroError struct {
	Err     string
	Details []error
}

func (e RetroError) Error() string {
	ers := e.Err
	for _, detail := range e.Details {
		ers += "\n" + detail.Error()
	}
	return ers
}

func GError(err string, details ...error) RetroError {
	return RetroError{
		Err:     err,
		Details: details,
	}
}
