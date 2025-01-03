package model

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
)

type Post struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Content   template.HTML `json:"content"`
	URL       template.URL  `json:"url"`
	Author    string        `json:"author"`
	Date      string        `json:"string"`
	Published string        `json:"published"`
}

func GetAllPostsFromUser(user string) ([]Post, error) {
	rows, err := db.Query("SELECT id, name, content, url, author, date, published FROM post WHERE author = ($1)", user)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID, &post.Name, &post.Content, &post.URL, &post.Author, &post.Date, &post.Published)
		if err != nil {
			slog.Error(err.Error())
		}
		if err != nil {
			slog.Error(err.Error())
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetAllPublishedPostsFromUser(user string) ([]Post, error) {
	rows, err := db.Query("SELECT id, name, content, url, author, date, published FROM post WHERE author = ($1) AND published = 'public'", user)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID, &post.Name, &post.Content, &post.URL, &post.Author, &post.Date, &post.Published)
		if err != nil {
			slog.Error(err.Error())
		}
		if err != nil {
			slog.Error(err.Error())
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetAllPosts() ([]Post, error) {
	rows, err := db.Query("SELECT id, name, content, url, author, date FROM post")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID, &post.Name, &post.Content, &post.URL, &post.Author, &post.Date)
		if err != nil {
			slog.Error(err.Error())
		}
		if err != nil {
			slog.Error(err.Error())
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func NewPost(newPost Post) error {
	_, err := db.Exec("INSERT INTO post (name, content, url, author, date, published) VALUES ($1, $2, $3, $4, $5, $6)",
		newPost.Name, newPost.Content, newPost.URL, newPost.Author, newPost.Date, newPost.Published)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePost(post Post) error {
	_, err := db.Exec(`
		UPDATE post
		SET name = ?, content = ?, url = ?, author = ?, date = ?
		WHERE id = ?`,
		post.Name, post.Content, post.URL, post.Author, post.Date, post.ID) // Assuming 'id' is the unique identifier
	if err != nil {
		return err
	}
	return nil
}

func GetPost(id int) (*Post, error) {
	rows, err := db.Query("SELECT id, name, content, url, author, date FROM post where id = ($1)", id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	var post Post
	if !rows.Next() {
		err := fmt.Errorf("Post not found.")
		slog.Error(err.Error())
		return nil, err
	}
	err = rows.Scan(&post.ID, &post.Name, &post.Content, &post.URL, &post.Author, &post.Date)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &post, nil
}

func GetPostByUrl(url string) (*Post, error) {
	rows, err := db.Query("SELECT id, name, content, url, author, date, published FROM post WHERE url = ($1)", url)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	var post Post
	if !rows.Next() {
		err := fmt.Errorf("Post not found.")
		slog.Error(err.Error())
		return nil, err
	}
	err = rows.Scan(&post.ID, &post.Name, &post.Content, &post.URL, &post.Author, &post.Date, &post.Published)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &post, nil
}

func DeletePost(id int) error {
	query := "DELETE FROM post WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not execute delete query: %w", err)
	}

	// Check the number of affected rows (optional)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not fetch affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found with ID %d", id)
	}

	return nil
}

func PublishPost(id int) error {
	query := "UPDATE post SET published = 'public' WHERE ID = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not execute delete query: %w", err)
	}

	// Check the number of affected rows (optional)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not fetch affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found with ID %d", id)
	}
	return nil
}

func UnpublishPost(id int) error {
	query := "UPDATE post SET published = 'draft' WHERE ID = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not execute delete query: %w", err)
	}

	// Check the number of affected rows (optional)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not fetch affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found with ID %d", id)
	}
	return nil
}

func UpvotePost(id int, ip string) error {
	// Insert upvote into the database
	query := `INSERT INTO upvotes (ip, post_id) VALUES (?, ?)`
	_, err := db.Exec(query, ip, id)
	if err != nil {
		if sqliteErr, ok := err.(interface{ ErrorCode() int }); ok && sqliteErr.ErrorCode() == 2067 {
			// UNIQUE constraint failed
			return fmt.Errorf("You have already upvoted this post.")
		}
		log.Println(err)
		return fmt.Errorf("Could not register upvote.")
	}
	return nil
}

func GetUpvotesByPostID(id int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM upvotes WHERE post_id = ?`
	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("Could not fetch upvote count.")
	}

	return count, nil
}
