package database

// Structure représentant une catégorie
type Categorie struct {
	CategorieID int
	Name        string
}

// Structure représentant un commentaire
type Commentaire struct {
	CommentaireID int
	SenderID      int
	PostID        int
	Like          string
	Dislike       string
	Date          string
	Content       string
}

// Structure représentant un post
type Post struct {
	PostID         int
	CategorieID    int
	Title          string
	Content        string
	Date           string
	SenderID       int
	SenderUsername string
	Image          string
	Like           string
	Dislike        string
}

// Structure représentant un utilisateur
type User struct {
	UserID         int
	Email          string
	Pseudo         string
	Password       string
	Role           string
	ProfilePicture string
	FollowID       string
}
