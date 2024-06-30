// models/posts.go
package models

import (
	"database/sql"
    "fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Posts struct {
    ID          uuid.UUID    `json:"id"`
    Title       string       `json:"title"`
    Content     string       `json:"content"`
    Status      string       `json:"status"`
    PublishDate sql.NullTime `json:"publish_date"`
    Tags        []string     `json:"tags"`
}

func (p *Posts) Create(db *sql.DB) (map[string]interface{}, error) {
    tx, err := db.Begin()
    if err != nil {
        return nil,err
    }
    defer tx.Rollback()

    query := `INSERT INTO posts (id, title, content, status, publish_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    p.ID = uuid.New()
    err = tx.QueryRow(query, p.ID, p.Title, p.Content, p.Status, p.PublishDate).Scan(&p.ID)
    if err != nil {
        return nil,err
    }

    var tagIDs []uuid.UUID
    var tagLabels []string

    for _, tag := range p.Tags {
        var tagID uuid.UUID
        err = tx.QueryRow("INSERT INTO tags (label) VALUES ($1) ON CONFLICT (label) DO UPDATE SET label = EXCLUDED.label RETURNING id", tag).Scan(&tagID)
        if err != nil {
            return nil,err
        }

        _, err = tx.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)", p.ID, tagID)
        if err != nil {
            fmt.Printf("Error post_tag IDs: %v\n", err)
            return nil,err
        }
        tagLabels = append(tagLabels, tag)
        tagIDs = append(tagIDs, tagID)
    }

    fmt.Printf("Tag IDs: %v\n", tagIDs)
    _, err = tx.Exec("UPDATE posts SET tags = $1 WHERE id = $2", pq.Array(tagIDs), p.ID)
    if err != nil {
        return nil,err
    }

    err = tx.Commit()

    if err != nil {
        return nil, err
    }

    result := map[string]interface{}{
        "id":      p.ID,
        "title":   p.Title,
        "content": p.Content,
        "tags":    tagLabels,
    }
    
    return result, nil
}

func (p *Posts) Update(db *sql.DB) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Update post
    query := `UPDATE posts SET title = $2, content = $3, status = $4, publish_date = $5 WHERE id = $1`
    _, err = tx.Exec(query, p.ID, p.Title, p.Content, p.Status, p.PublishDate)
    if err != nil {
        return err
    }

    // Delete existing post_tags
    _, err = tx.Exec("DELETE FROM post_tags WHERE post_id = $1", p.ID)
    if err != nil {
        return err
    }

    // Handle tags
    for _, tag := range p.Tags {
        var tagID uuid.UUID
        // Try to insert the tag, if it already exists, get its ID
        err = tx.QueryRow("INSERT INTO tags (label) VALUES ($1) ON CONFLICT (label) DO UPDATE SET label = EXCLUDED.label RETURNING id", tag).Scan(&tagID)
        if err != nil {
            return err
        }

        // Insert into post_tags junction table
        _, err = tx.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)", p.ID, tagID)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func GetPost(db *sql.DB, id uuid.UUID) (Posts, error) {
    p := Posts{}
    query := `SELECT id, title, content, status, publish_date FROM posts WHERE id = $1`
    err := db.QueryRow(query, id).Scan(&p.ID, &p.Title, &p.Content, &p.Status, &p.PublishDate)
    if err != nil {
        return p, err
    }

    // Fetch tags
    rows, err := db.Query("SELECT t.label FROM tags t JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = $1", id)
    if err != nil {
        return p, err
    }
    defer rows.Close()

    for rows.Next() {
        var tag string
        if err := rows.Scan(&tag); err != nil {
            return p, err
        }
        p.Tags = append(p.Tags, tag)
    }

    return p, nil
}

