package db

import (
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
	if categoryID != 0 && CategoryExist(GetCategoryNameByID(categoryID)) {
		query = `
        SELECT p.id, p.category, p.sender, p.title, p.content, p.picture, p.date, c.name, u.role, u.username, u.email, u.picture, u.password
        FROM posts p
        JOIN categories c ON p.category = c.id
        JOIN users u ON p.sender = u.id
        WHERE p.category = ?
        ORDER BY p.id DESC;`
		args = append(args, categoryID)
	} else if postID != 0 && PostExist(postID) {
		query = `
        SELECT p.id, p.category, p.sender, p.title, p.content, p.picture, p.date, c.name, u.role, u.username, u.email, u.picture, u.password
        FROM posts p
        JOIN categories c ON p.category = c.id
        JOIN users u ON p.sender = u.id
        WHERE p.id = ?;`
		args = append(args, postID)
	}
	if query == "" {
		query = `
        SELECT p.id, p.category, p.sender, p.title, p.content, p.picture, p.date, c.name, u.role, u.username, u.email, u.picture, u.password
        FROM posts p
        JOIN categories c ON p.category = c.id
        JOIN users u ON p.sender = u.id
        ORDER BY p.id DESC;`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		// Gérer l'erreur, par exemple en la journalisant
		log.Printf("Erreur lors de l'exécution de la requête : %v", err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Category.ID, &post.Sender.ID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Category.Name, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture, &post.Sender.Password)
		if err != nil {
			// Gérer l'erreur, par exemple en la journalisant
			log.Printf("Erreur lors du scan de la ligne : %v", err)
			continue
		}

		likes, dislikes, err := GetPostReactions(post.ID)
		if err != nil {
			// Gérer l'erreur, par exemple en la journalisant
			log.Printf("Erreur lors de la récupération des réactions : %v", err)
			continue
		}
		post.Likes = likes
		post.Dislikes = dislikes
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		// Gérer l'erreur finale, si elle existe
		log.Printf("Erreur lors de l'itération sur les lignes : %v", err)
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
