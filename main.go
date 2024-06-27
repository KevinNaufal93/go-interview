package main

import (
	"fmt"
	"log"
	"go-interview/controllers"
	"go-interview/configs"
)

func main() {
	database, err := configs.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	userController := controllers.PostsController{DB: database}

	for {
		fmt.Println("\n1. Create User")
		fmt.Println("2. Get User")
		fmt.Println("3. Update User")
		fmt.Println("4. Delete User")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var title, content string
			var status = "draft"
			var publishDate = "2024-06-04"
			fmt.Print("Enter title: ")
			fmt.Scanln(&title)
			fmt.Print("Enter content: ")
			fmt.Scanln(&content)
			userController.CreateUser(title, content, status, publishDate)
		case 2:
			var id string
			fmt.Print("Enter user ID: ")
			fmt.Scanln(&id)
			userController.GetUser(id)
		case 3:
			var id string
			var title, content string
			var status = "publish"
			var publishDate = "2024-06-04"
			fmt.Print("Enter user ID: ")
			fmt.Scanln(&id)
			fmt.Print("Enter new title: ")
			fmt.Scanln(&title)
			fmt.Print("Enter new content: ")
			fmt.Scanln(&content)
			userController.UpdateUser(id, title, content, status, publishDate)
		case 4:
			var id string
			fmt.Print("Enter user ID to delete: ")
			fmt.Scanln(&id)
			userController.DeleteUser(id)
		case 5:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}