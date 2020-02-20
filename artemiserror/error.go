package artemiserror

type ExceptionType string

const (
	ArgsNotInitialized ExceptionType = "Arguments not initialized"
	PathNotSet         ExceptionType = "Path not set"
	InvalidName        ExceptionType = "Invalid name"
	InvalidParameter   ExceptionType = "Invalid parameter"
)

type ArtemisError struct {
	message string
}

func (e ArtemisError) Error() string {
	return e.message
}

func New(et ExceptionType) ArtemisError {
	return ArtemisError{string(et)}
}
