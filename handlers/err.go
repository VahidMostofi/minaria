package handlers

type GenericError struct {
	Message        string      `json:"message"`
	AdditionalInfo interface{} `json:"more"`
	Err            error       `json:"-"`
	HTTPStatusCode int         `json:"-"`
}
