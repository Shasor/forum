package db

type User struct {
	ID       int
	Role     string
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
	ID       int
	Category Category
	Sender   User
	Title    string
	Content  string
	Picture  string
	Date     string
	Likes    int
	Dislikes int
}

type Reaction struct {
	ID     int
	Sender User
	Post   Post
	Value  string
}
