package db

import (
	"fmt"
	"log"
)

func CreatePost(sender int, categoryName, title, content, picture, date string) error {
	// Ouvrir la connexion à la base de données
	db := GetDB()
	defer db.Close()

	_ = CreateCategory(categoryName)
	category, _ := SelectCategoryByName(categoryName)

	// Démarrer une transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Préparer la requête d'insertion
	stmt, err := tx.Prepare("INSERT INTO posts (category, sender, title, content, picture, date) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Exécuter l'insertion
	_, err = stmt.Exec(category.ID, sender, title, content, picture, date)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit de la transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func FetchPosts(categoryID, postID int) []Post {
	db := GetDB()
	defer db.Close()

	var query string
	var args []interface{}

	// Use LEFT JOIN to allow posts from deleted users
	if categoryID != 0 && CategoryExist(GetCategoryNameByID(categoryID)) {
		query = `
        SELECT p.id, p.category, p.sender, p.title, p.content, p.picture, p.date, c.name,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
        FROM posts p
        JOIN categories c ON p.category = c.id
        LEFT JOIN users u ON p.sender = u.id
        WHERE p.category = ?
        ORDER BY p.id DESC;`
		args = append(args, categoryID)
	} else if postID != 0 && PostExist(postID) {
		query = `
        SELECT p.id, p.category, p.sender, p.title, p.content, p.picture, p.date, c.name,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
        FROM posts p
        JOIN categories c ON p.category = c.id
        LEFT JOIN users u ON p.sender = u.id
        WHERE p.id = ?;`
		args = append(args, postID)
	} else {
		query = `
        SELECT p.id, p.category, p.sender, p.title, p.content, p.picture, p.date, c.name,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
        FROM posts p
        JOIN categories c ON p.category = c.id
        LEFT JOIN users u ON p.sender = u.id
        ORDER BY p.id DESC;`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Category.ID, &post.Sender.ID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Category.Name, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Get post reactions
		likes, dislikes, err := GetPostReactions(post.ID)
		if err != nil {
			log.Printf("Error fetching reactions: %v", err)
			continue
		}
		post.Likes = likes
		post.Dislikes = dislikes

		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
	}
	return posts
}

// FetchPostsLiked retrieves all posts that the user with the specified senderID liked.
func FetchPostsLiked(senderID int) []Post {
	db := GetDB()
	// Query to fetch posts that the user liked, joining Post and PostReaction tables
	query := `
		SELECT p.id, p.category, p.sender, p.title, p.content, p.picture, p.date, c.name, u.role, u.username, u.email, u.picture, u.password
		FROM posts p 
		JOIN categories c ON p.category = c.id
		JOIN reactions r ON p.id = r.post
		JOIN users u ON p.sender = u.id
		WHERE r.sender = ? AND r.value = 'LIKE';
    `

	// Execute the query
	rows, err := db.Query(query, senderID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	// Slice to hold the posts
	var likedPosts []Post

	// Loop through the result set and scan each row into a Post struct
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Category.ID, &post.Sender.ID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Category.Name, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture, &post.Sender.Password); err != nil {
			fmt.Println(err)
			return nil
		}
		likedPosts = append(likedPosts, post)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	// Return the list of liked posts
	return likedPosts
}

func PostExist(id int) bool {
	db := GetDB()
	defer db.Close()
	var exist bool
	err := db.QueryRow("SELECT EXISTS( SELECT 1 FROM posts WHERE id = ?)", id).Scan(&exist)
	if err != nil {
		return false
	}
	return exist
}
