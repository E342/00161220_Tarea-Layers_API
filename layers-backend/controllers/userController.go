package controllers

import (
	"encoding/json"
	"layersapi/entities"
	"layersapi/entities/dto"
	"layersapi/services"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (u UserController) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	resData, err := u.userService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resData)
}

func (u UserController) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Extrayendo los parámetros de la URL
	id := vars["id"]

	user, err := u.userService.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error adapting to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// CreateUserHandler maneja la creación de un nuevo usuario
// y utiliza un DTO para recibir los datos de entrada.
func (u UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user entities.User

	// Se espera que el cuerpo de la solicitud contenga un JSON con los campos
	// a utilizar para crear un nuevo usuario.
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Creando metadata
	metadata := services.CreateMetadata()

	// Mapeando a DTO
	createData := dto.CreateUser{
		Name:  user.Name,
		Email: user.Email,
	}
	createData.Metadata.CreatedAt = metadata.CreatedAt
	createData.Metadata.CreatedBy = metadata.CreatedBy
	createData.Metadata.UpdatedAt = metadata.UpdatedAt
	createData.Metadata.UpdatedBy = metadata.UpdatedBy

	// Llamando a Service
	err = u.userService.Create(createData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created successfully"))
}

// UpdateUserHandler maneja la actualización de un usuario existente
// y utiliza un DTO para recibir los datos de entrada.
func (u UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// El ID del usuario se extrae de la URL.
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Se espera que el cuerpo de la solicitud contenga un JSON con los campos
	// necesarios para actualizar el usuario.
	var input dto.UpdateUser
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Actualizando campos de metadata automáticamente
	input.Metadata.UpdatedAt = time.Now().Format(time.RFC3339)
	input.Metadata.UpdatedBy = "webapp"

	err = u.userService.Update(id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User updated successfully"))
}

// DeleteUserHandler maneja la eliminación de un usuario
// y utiliza un DTO para recibir los datos de entrada.

func (u UserController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// En este caso, solo se necesita el ID del usuario.
	id := vars["id"]

	err := u.userService.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Write([]byte("user delete successfully"))
}
