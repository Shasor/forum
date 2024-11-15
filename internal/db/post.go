package db

import (
	"log"
)

// CreatePost handles both posts and comments, depending on parentID
func CreatePost(sender int, categoryName, title, content, picture, date string, parentID *int) error {
	// Open the database connection
	db := GetDB()
	defer db.Close()

	// Create the category if it doesn't exist
	_ = CreateCategory(categoryName)
	category, _ := SelectCategoryByName(categoryName)

	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the SQL statement for inserting a post (including parent_id)
	stmt, err := tx.Prepare("INSERT INTO posts (category, sender, parent_id, title, content, picture, date) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// If parentID is nil, it means it's an original post, so we set parent_id to 0
	if parentID == nil {
		_, err = stmt.Exec(category.ID, sender, 0, title, content, picture, date)
	} else {
		// Otherwise, use the provided parent_id (for comments)
		_, err = stmt.Exec(category.ID, sender, *parentID, title, content, picture, date)
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func FetchPosts() []Post {
	db := GetDB()
	defer db.Close()

	query := `
        SELECT p.id, p.category, p.sender, p.parent_id, p.title, p.content, p.picture, p.date, c.name,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
        FROM posts p
        JOIN categories c ON p.category = c.id
        LEFT JOIN users u ON p.sender = u.id
        ORDER BY p.id DESC;`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Category.ID, &post.Sender.ID, &post.ParentID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Category.Name, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture)
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
	defer db.Close()

	// Query to fetch posts that the user liked, with LEFT JOIN for users and placeholders for deleted users
	query := `
        SELECT p.id, p.category, p.sender, p.parent_id, p.title, p.content, p.picture, p.date, c.name,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
        FROM posts p 
        JOIN categories c ON p.category = c.id
        JOIN reactions r ON p.id = r.post
        LEFT JOIN users u ON p.sender = u.id
        WHERE r.sender = ? AND r.value = 'LIKE';
    `

	// Execute the query
	rows, err := db.Query(query, senderID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil
	}
	defer rows.Close()

	// Slice to hold the posts
	var likedPosts []Post

	// Loop through the result set and scan each row into a Post struct
	for rows.Next() {
		var post Post
		// Scan each row's values, including placeholders for missing user info
		if err := rows.Scan(
			&post.ID, &post.Category.ID, &post.Sender.ID, &post.ParentID, &post.Title, &post.Content,
			&post.Picture, &post.Date, &post.Category.Name,
			&post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil
		}
		likedPosts = append(likedPosts, post)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		return nil
	}

	// Return the list of liked posts
	return likedPosts
}

func FetchComments() []Post {
	db := GetDB()
	defer db.Close()

	query := `
        SELECT p.id, p.category, p.sender, p.parent_id, p.title, p.content, p.picture, p.date, c.name,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
        FROM posts p
        JOIN categories c ON p.category = c.id
        LEFT JOIN users u ON p.sender = u.id
        ORDER BY p.id DESC;`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Category.ID, &post.Sender.ID, &post.ParentID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Category.Name, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture)
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

		if post.ParentID != 0 {
			posts = append(posts, post)
		}
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
	}
	return posts
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

func NbCommentsFromPost(id int) (int, error) {
	db := GetDB()
	defer db.Close()
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM posts WHERE parent_id =?", id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
