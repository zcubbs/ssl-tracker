package api

type MalformedRequest struct {
	Status int
	Msg    string
}

type AuthenticationError struct {
	Status int
	Msg    string
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}

func (ae *AuthenticationError) Error() string {
	return ae.Msg
}
