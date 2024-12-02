package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

// CreatePost handles both posts and comments, depending on parentID
func CreatePost(sender int, categories []string, title, content, picture, date string, parentID *int) error {
	// Open the database connection
	db := GetDB()
	defer db.Close()

	category, _ := SelectCategoryByName(categories[0])

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

	postID, _ := GetLastPostIDByUserID(sender)
	post, _ := SelectPostByID(postID)

	var categoriesID []int
	for i, category := range categories {
		category = strings.TrimSpace(category)
		id, _ := GetCategoryIDByName(category)
		categoriesID = append(categoriesID, id)
		err = LinkPostToCategory(post.ID, categoriesID[i])
		if err != nil {
			panic(err)
		}
	}
	if parentID == nil {
		addActivity(sender, post.ID, "post")
		for _, categoryID := range categoriesID {
			receivers, _ := GetUsersByFollowedCategory(categoryID)
			for _, receiver := range receivers {
				AddNotification("category", date, sender, receiver.ID, post.ID, 0)
			}
		}
	} else {
		addActivity(sender, post.ID, "comment")
		parentPost, _ := SelectPostByID(post.ParentID)
		AddNotification("post", date, sender, parentPost.Sender.ID, post.ID, parentPost.ID)
	}

	return nil
}

func FetchPosts() []Post {
	db := GetDB()
	defer db.Close()

	query := `
        SELECT p.id, p.sender, p.parent_id, p.title, p.content, p.picture, p.date,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
        FROM posts p
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
		err := rows.Scan(&post.ID, &post.Sender.ID, &post.ParentID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture)
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

		post.Categories, _ = GetPostCategories(post.ID)

		post.NbComments, err = NbCommentsFromPost(post.ID)
		if err != nil {
			post.NbComments = 0
			fmt.Println("Error at fetching nb comments: ", err)
		}

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
        SELECT p.id, p.sender, p.parent_id, p.title, p.content, p.picture, p.date,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
        FROM posts p 
        JOIN categories c ON p.category = c.id
        JOIN reactions r ON p.id = r.post
        LEFT JOIN users u ON p.sender = u.id
        WHERE r.user = ? AND r.value = 'LIKE';
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
			&post.ID, &post.Sender.ID, &post.ParentID, &post.Title, &post.Content,
			&post.Picture, &post.Date,
			&post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil
		}

		likes, dislikes, err := GetPostReactions(post.ID)
		if err != nil {
			log.Printf("Error fetching reactions: %v", err)
			continue
		}

		post.Likes = likes
		post.Dislikes = dislikes

		post.Categories, _ = GetPostCategories(post.ID)

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
        SELECT p.id, p.sender, p.parent_id, p.title, p.content, p.picture, p.date,
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
		err := rows.Scan(&post.ID, &post.Sender.ID, &post.ParentID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture)
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

		post.Categories, _ = GetPostCategories(post.ID)

		post.NbComments, err = NbCommentsFromPost(post.ID)
		if err != nil {
			post.NbComments = 0
			fmt.Println("Error at fetching nb comments: ", err)
		}

		if post.ParentID != 0 {
			posts = append(posts, post)
		}
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
	}
	return posts
}

func FetchFollowPosts(senderID int) []Post {
	db := GetDB()
	defer db.Close()

	query := `
        SELECT p.id, p.sender, p.parent_id, p.title, p.content, p.picture, p.date,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
        FROM posts p
        JOIN post_category pc ON p.id = pc.post_id
		JOIN follows f ON pc.category_id = f.category
        LEFT JOIN users u ON p.sender = u.id
        WHERE f.user = ?
		GROUP BY p.id
        ORDER BY p.id DESC;`

	rows, err := db.Query(query, senderID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Sender.ID, &post.ParentID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture)
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
		post.Categories, _ = GetPostCategories(post.ID)
		post.Likes = likes
		post.Dislikes = dislikes

		post.NbComments, err = NbCommentsFromPost(post.ID)
		if err != nil {
			post.NbComments = 0
			fmt.Println("Error at fetching nb comments: ", err)
		}

		posts = append(posts, post)
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

func GetLastPostIDByUserID(id int) (int, error) {
	db := GetDB()
	defer db.Close()

	query := `
        SELECT id
        FROM posts
		WHERE sender = ?
		ORDER BY id DESC
		LIMIT 1;`
	var postID int

	err := db.QueryRow(query, id).Scan(&postID)

	if err != nil {
		// Gérer l'erreur
		log.Printf("Erreur lors de l'exécution de la requête : %v", err)
		return 0, err
	}

	return postID, nil
}

func GetPostFromUserById(id int) []Post {
	db := GetDB()
	defer db.Close()

	query := `
        SELECT p.id, p.sender, p.parent_id, p.title, p.content, p.picture, p.date,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
		FROM posts p
        LEFT JOIN users u ON p.sender = u.id
		WHERE p.sender = ?
        ORDER BY p.id DESC;`

	rows, err := db.Query(query, id)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Sender.ID, &post.ParentID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture)
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

		post.Categories, _ = GetPostCategories(post.ID)

		post.NbComments, err = NbCommentsFromPost(post.ID)
		if err != nil {
			post.NbComments = 0
			fmt.Println("Error at fetching nb comments: ", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
	}
	return posts

}

func SortPostsByDateDesc(posts []Post) {
	// Définir le layout pour analyser les dates
	const layout = "15:04 | 02/01/2006"

	sort.Slice(posts, func(i, j int) bool {
		// Convertir les dates en objets time.Time
		date1, err1 := time.Parse(layout, posts[i].Date)
		date2, err2 := time.Parse(layout, posts[j].Date)

		// Si la conversion échoue, on considère la date comme moins récente
		if err1 != nil || err2 != nil {
			return err1 != nil
		}

		// Trier par ordre antichronologique (plus récente en premier)
		return date1.After(date2)
	})
}

func FetchPostsReactions(userID int) []Post {
	db := GetDB()
	defer db.Close()

	// Query to fetch posts that the user liked, with LEFT JOIN for users and placeholders for deleted users
	query := `
        SELECT p.id, p.sender, p.parent_id, p.title, p.content, p.picture, p.date,
               IFNULL(u.role, 'Deleted') AS role,
               IFNULL(u.username, 'Deleted User') AS username,
               IFNULL(u.email, '') AS email,
               IFNULL(u.picture, 'default-profile.png') AS picture
        FROM posts p 
        JOIN categories c ON p.category = c.id
        JOIN reactions r ON p.id = r.post
        LEFT JOIN users u ON p.sender = u.id
        WHERE r.user = ?;
    `

	// Execute the query
	rows, err := db.Query(query, userID)
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
			&post.ID, &post.Sender.ID, &post.ParentID, &post.Title, &post.Content,
			&post.Picture, &post.Date,
			&post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil
		}

		likes, dislikes, err := GetPostReactions(post.ID)
		if err != nil {
			log.Printf("Error fetching reactions: %v", err)
			continue
		}

		post.Likes = likes
		post.Dislikes = dislikes

		post.Categories, _ = GetPostCategories(post.ID)

		// post.Reactions = GetReactionsByPostID(post.ID)

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

func SelectPostByID(postID int) (Post, error) {
	db := GetDB()
	defer db.Close()

	query := `
			SELECT p.id, p.sender, p.parent_id, p.title, p.content, p.picture, p.date,
				   IFNULL(u.role, 'Deleted') AS role,
				   IFNULL(u.username, 'Deleted User') AS username,
				   IFNULL(u.email, '') AS email,
				   IFNULL(u.picture, 'default-profile.png') AS picture
			FROM posts p
			LEFT JOIN users u ON p.sender = u.id
			WHERE p.id = ?;`

	var post Post
	err := db.QueryRow(query, postID).Scan(&post.ID, &post.Sender.ID, &post.ParentID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture)
	if err != nil {
		if err == sql.ErrNoRows {
			return Post{}, errors.New("post not found")
		}
		return Post{}, err
	}
	// Get post reactions
	likes, dislikes, err := GetPostReactions(post.ID)
	if err != nil {
		log.Printf("Error fetching reactions: %v", err)
		return Post{}, err
	}
	post.Likes = likes
	post.Dislikes = dislikes

	post.Categories, _ = GetPostCategories(post.ID)

	post.NbComments, err = NbCommentsFromPost(post.ID)
	if err != nil {
		post.NbComments = 0
		fmt.Println("Error at fetching nb comments: ", err)
	}
	return post, nil
}

func DeletePostByID(postID int) error {
	db := GetDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	recursive := `
		WITH RECURSIVE comments_tree AS (
    		SELECT id FROM posts WHERE id = ?
    		UNION ALL
    		SELECT p.id FROM posts p
    		JOIN comments_tree ct ON p.parent_id = ct.id
		)`

	queries := []string{
		"DELETE FROM post_category WHERE post_id IN (SELECT id FROM comments_tree);",
		"DELETE FROM reactions WHERE post IN (SELECT id FROM comments_tree);",
		"DELETE FROM activity WHERE post IN (SELECT id FROM comments_tree);",
		"DELETE FROM posts WHERE id IN (SELECT id FROM comments_tree);",
	}

	rowsAffected := int64(0)
	for _, query := range queries {
		query = recursive + query
		result, err := tx.Exec(query, postID)
		if err != nil {
			tx.Rollback()
			return err
		}
		affected, _ := result.RowsAffected()
		rowsAffected += affected
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("the requested post doesn't exist")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func ModifyContentPostByID(postID int, content string) error {
	db := GetDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	update := `UPDATE posts SET content = ? WHERE id = ?;`

	if _, err := tx.Exec(update, content, postID); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
