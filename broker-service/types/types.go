package types

// JsonResponse is one instance of a json response to the frontend
type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// RequestPayload is the main upper struct service related structs are embed into this struct
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

// MailPayload is one instance of an email message in json
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// AuthPayload is one instance of an authentication json message
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LogPayload is one instance of a json log message
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
