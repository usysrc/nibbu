package model

import (
	"fmt"
	"html/template"
	"log/slog"
)

type Post struct {
	ID      int           `json:"id"`
	Name    string        `json:"name"`
	Content template.HTML `json:"content"`
	URL     template.URL  `json:"url"`
	Author  string        `json:"author"`
	Date    string        `json:"string"`
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
	_, err := db.Exec("INSERT into post (name) VALUES ($1)", newPost.Name)
	if err != nil {
		return err
	}
	return nil
}

func GetPost(id int) (*Post, error) {
	rows, err := db.Query("SELECT id, name FROM items where id = ($1)", id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	var item Post
	if !rows.Next() {
		err := fmt.Errorf("Item not found.")
		slog.Error(err.Error())
		return nil, err
	}
	err = rows.Scan(&item.ID, &item.Name)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &item, nil
}
