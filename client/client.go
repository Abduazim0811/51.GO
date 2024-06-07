package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"User/server/models"
)

func readUsersFromFile(filename string) ([]models.User, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []models.User
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}
		users = append(users, models.User{Name: parts[0], Email: parts[1]})
	}
	return users, scanner.Err()
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err)
	}

	users, err := readUsersFromFile("users.txt")
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		var reply models.User
		err = client.Call("UserService.CreateUser", user, &reply)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created user: %+v\n", reply)
	}

	var user models.User
	err = client.Call("UserService.GetUser", 1, &user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User 1: %+v\n", user)

	user.Name = "Updated Name"
	err = client.Call("UserService.UpdateUser", user, &user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated user: %+v\n", user)

	var success bool
	err = client.Call("UserService.DeleteUser", 1, &success)
	if err != nil {
		log.Fatal(err)
	}
	if success {
		fmt.Println("Deleted user 1")
	}

	var usersList []models.User
	err = client.Call("UserService.ListUsers", struct{}{}, &usersList)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Users: %+v\n", usersList)
}
