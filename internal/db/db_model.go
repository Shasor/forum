package db

type User struct {
	ID       int
	Role     string
	Follows  []Category
	Username string
	Email    string
	Picture  string
	Password string
}

type Category struct {
	ID   int
	Name string
}

type Post struct {
	ID         int
	Categories []Category
	Sender     User
	ParentID   int
	Title      string
	Content    string
	Picture    string
	Date       string
	Likes      int
	Dislikes   int
	NbComments int
}

type Reaction struct {
	ID     int
	Sender User
	Post   Post
	Value  string
}

type Activity struct {
	ID     int
	User   User
	Post   Post
	Action string
}

type Notification struct {
	ID       int
	Type     string
	Sender   User
	Receiver User
	Post     Post
	Content  string
	Date     string
	Readed	 bool
}
