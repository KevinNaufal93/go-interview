// server.go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "go-interview/controllers"
    "go-interview/configs"
)

type Server struct {
    pc *controllers.PostsController
}

func NewServer() (*Server, error) {
    db, err := configs.InitDB()
    if err != nil {
        return nil, fmt.Errorf("init db failed >>>>>>>>> %v", err)
    }
    pc := &controllers.PostsController{DB: db}
    return &Server{
        pc: pc,
    }, nil
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    switch {
	case req.Method == "POST" && req.URL.Path == "/api/posts":
        s.handleCreatePost(res, req)
    case req.Method == "PUT" && strings.HasPrefix(req.URL.Path, "/api/posts/"):
        s.handleUpdatePost(res, req)
    // case req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/api/posts/"):
    //     s.handleGetOnePost(res, req)
    // case req.Method == "GET" && req.URL.Path == "/api/posts":
    //     s.handleGetAllPost(res, req)
    // case req.Method == "DELETE" && strings.HasPrefix(req.URL.Path, "/api/posts/"):
    //     s.handleDeletePost(res, req)
    default:
        http.NotFound(res, req)
    }
}

func (s *Server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
    var post struct {
        Title   string   `json:"title"`
        Content string   `json:"content"`
        Tags    []string `json:"tags"`
    }

    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := s.pc.CreatePost(post.Title, post.Content, post.Tags)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func (s *Server) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/posts/")
    var post struct {
        Title       string   `json:"title"`
        Content     string   `json:"content"`
        Status      string   `json:"status"`
        PublishDate string   `json:"publish_date"`
        Tags        []string `json:"tags"`
    }

    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := s.pc.UpdatePost(id, post.Title, post.Content, post.Status, post.PublishDate, post.Tags)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}


