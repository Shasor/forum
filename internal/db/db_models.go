package db

type User struct {
	ID       int
	Role     Role
	Username string
	Email    string
	Picture  string
	Password string
}

type Role struct {
	ID   int
	Name string
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
	Date     string
	Comments []Post
}

type Reaction struct {
	ID     int
	Sender User
	Post   Post
	Value  string
}
