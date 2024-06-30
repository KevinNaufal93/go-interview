// controllers/posts_controller.go
package controllers

import (
	"database/sql"
	"go-interview/models"
	"time"
	"github.com/google/uuid"
)

type PostsController struct {
    DB *sql.DB
}

func (pc *PostsController) CreatePost(title, content string, tags []string) (map[string]interface{}, error) {
    post := models.Posts{
        Title:   title,
        Content: content,
        Tags:    tags,
        Status:  "draft",
        PublishDate: sql.NullTime{Time: time.Now(), Valid: true},
    }

    result, err := post.Create(pc.DB)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func (pc *PostsController) UpdatePost(idStr, title, content, status, publishDate string, tags []string) error {
    id, err := uuid.Parse(idStr)
    if err != nil {
        return err
    }
    post := models.Posts{ID: id, Title: title, Content: content, Status: status, Tags: tags}
    if publishDate != "" {
        t, err := time.Parse("2006-01-02", publishDate)
        if err == nil {
            post.PublishDate = sql.NullTime{Time: t, Valid: true}
        }
    }
    return post.Update(pc.DB)
}