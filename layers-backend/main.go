package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"layersapi/controllers"
	"layersapi/repositories/files/csv"
	"layersapi/services"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	userRepository := csv.NewUserRepository()
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(*userService)

	r.HandleFunc("/users/{id}", userController.GetUserByIDHandler).Methods(http.MethodGet)
	r.HandleFunc("/users", userController.CreateUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", userController.UpdateUserHandler).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", userController.DeleteUserHandler).Methods(http.MethodDelete)
	r.HandleFunc("/users", userController.GetAllUsersHandler).Methods(http.MethodGet)

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server error: ", err)
	}
}
