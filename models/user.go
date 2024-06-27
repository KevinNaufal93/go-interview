package models

import (
    "database/sql"
    "github.com/google/uuid"
)

type Posts struct {
    ID          uuid.UUID
    Title       string
    Content     string
    Status      string
    PublishDate sql.NullTime // Use sql.NullTime for nullable date
}

func (u *Posts) Create(db *sql.DB) error {
    query := `INSERT INTO posts (title, content, status, publish_date) VALUES ($1, $2, $3, $4) RETURNING id`
    return db.QueryRow(query, u.Title, u.Content, u.Status, u.PublishDate).Scan(&u.ID)
}

func GetUser(db *sql.DB, id uuid.UUID) (Posts, error) {
    u := Posts{}
    query := `SELECT id, title, content, status, publish_date FROM posts WHERE id = $1`
    err := db.QueryRow(query, id).Scan(&u.ID, &u.Title, &u.Content, &u.Status, &u.PublishDate)
    return u, err
}

func (u *Posts) Update(db *sql.DB) error {
    query := `UPDATE posts SET title = $2, content = $3, status = $4, publish_date = $5 WHERE id = $1`
    _, err := db.Exec(query, u.ID, u.Title, u.Content, u.Status, u.PublishDate)
    return err
}

func (u *Posts) Delete(db *sql.DB) error {
    query := `DELETE FROM posts WHERE id = $1`
    _, err := db.Exec(query, u.ID)
    return err
}