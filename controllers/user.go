package controllers

import (
    "database/sql"
    "go-interview/models"
    "go-interview/views"
    "github.com/google/uuid"
    "time"
)

type PostsController struct {
    DB *sql.DB
}

func (pc *PostsController) CreateUser(title, content, status, publishDate string) {
    post := models.Posts{
        Title:   title,
        Content: content,
        Status:  status,
    }
    if publishDate != "" {
        t, err := time.Parse("2006-01-02", publishDate)
        if err == nil {
            post.PublishDate = sql.NullTime{Time: t, Valid: true}
        }
    }
    err := post.Create(pc.DB)
    if err != nil {
        views.ShowError(err)
        return
    }
    views.ShowSuccess("Post created successfully")
    views.ShowPost(post)
}

func (pc *PostsController) GetUser(idStr string) {
    id, err := uuid.Parse(idStr)
    if err != nil {
        views.ShowError(err)
        return
    }
    post, err := models.GetUser(pc.DB, id)
    if err != nil {
        views.ShowError(err)
        return
    }
    views.ShowPost(post)
}

func (pc *PostsController) UpdateUser(idStr, title, content, status, publishDate string) {
    id, err := uuid.Parse(idStr)
    if err != nil {
        views.ShowError(err)
        return
    }
    post := models.Posts{ID: id, Title: title, Content: content, Status: status}
    if publishDate != "" {
        t, err := time.Parse("2006-01-02", publishDate)
        if err == nil {
            post.PublishDate = sql.NullTime{Time: t, Valid: true}
        }
    }
    err = post.Update(pc.DB)
    if err != nil {
        views.ShowError(err)
        return
    }
    views.ShowSuccess("Post updated successfully")
    views.ShowPost(post)
}

func (pc *PostsController) DeleteUser(idStr string) {
    id, err := uuid.Parse(idStr)
    if err != nil {
        views.ShowError(err)
        return
    }
    post := models.Posts{ID: id}
    err = post.Delete(pc.DB)
    if err != nil {
        views.ShowError(err)
        return
    }
    views.ShowSuccess("Post deleted successfully")
}