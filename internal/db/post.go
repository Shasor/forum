package db

func CreatePost(sender int, categoryName, title, content, picture, date string) error {
	// Ouvrir la connexion à la base de données
	db := GetDB()
	defer db.Close()

	_ = CreateCategorie(categoryName)
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

func FetchPosts() []Post {
	db := GetDB()
	defer db.Close()

	query := `
		SELECT p.id, p.category, p.sender, p.title, p.content, p.picture, p.date, c.name, u.role, u.username, u.email, u.picture, u.password
		FROM posts p
		JOIN categories c ON p.category = c.id
		JOIN users u ON p.sender = u.id
		ORDER BY p.id DESC;`
	rows, err := db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Category.ID, &post.Sender.ID, &post.Title, &post.Content, &post.Picture, &post.Date, &post.Category.Name, &post.Sender.Role, &post.Sender.Username, &post.Sender.Email, &post.Sender.Picture, &post.Sender.Password)
		if err != nil {
			return nil
		}
		posts = append(posts, post)
	}
	return posts
}
