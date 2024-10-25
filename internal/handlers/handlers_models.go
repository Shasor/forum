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

type ReactionRequest struct {
	PostID   int    `json:"postId"`
	Reaction string `json:"reaction"`
}

var (
	Resp Response
)
