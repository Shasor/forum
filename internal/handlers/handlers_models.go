package handlers

type Response struct {
	Msg         []string
	Action      string
	Form        Form
	Broadcasted bool
}

type Form struct {
	Username string
	Email    string
}

var (
	Resp Response
)
