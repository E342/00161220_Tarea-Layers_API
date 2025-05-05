package services

import (
	"errors"
	"layersapi/entities"
	"layersapi/entities/dto"
	"layersapi/repositories"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// Definiendo expresiones para validar nombre y correo electrónico
var (
	nameRegex  = regexp.MustCompile(`^[A-Za-z\s]+$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// Definiendo estructura del servicio que contiene una instancia del repositorio de usuarios
type UserService struct {
	userRepository repositories.UserRepository
}

// Creando una nueva instancia del servicio de usuario
func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

// Validando que el nombre no esté vacío y contenga solo letras y espacios
func (u UserService) validateName(name string) error {
	if len(name) == 0 {
		return errors.New("name cannot be empty")
	}
	if !nameRegex.MatchString(name) {
		return errors.New("name must only contain alphabetic characters and spaces")
	}
	return nil
}

// Validando que el correo electrónico tenga el formato correcto
func (u UserService) validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email address")
	}
	return nil
}

// Creando un nuevo usuario si el nombre y el correo electrónico son válidos
func (u UserService) Create(user dto.CreateUser) error {
	if err := u.validateName(user.Name); err != nil {
		return err
	}
	if err := u.validateEmail(user.Email); err != nil {
		return err
	}

	// Obteniendo todos los usuarios
	// y verificando duplicados segun su correo
	users, _ := u.userRepository.GetAll()
	for _, existing := range users {
		if existing.Email == user.Email {
			return errors.New("email already exists")
		}
	}

	// Generando UUID y metadatos de auditoría
	id, _ := uuid.NewUUID()
	now := time.Now().Format(time.RFC3339)
	meta := entities.Metadata{
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "webapp",
		UpdatedBy: "webapp",
	}

	// Creando la entidad de usuario y guardándola en el repositorio
	newUser := entities.NewUser(id.String(), user.Name, user.Email, meta)
	return u.userRepository.Create(newUser)
}

// Actualizando los datos de un usuario existente
func (u UserService) Update(id string, user dto.UpdateUser) error {
	// Validando nombre y correo
	if err := u.validateName(user.Name); err != nil {
		return err
	}
	if err := u.validateEmail(user.Email); err != nil {
		return err
	}

	// Verificando si el usuario existe
	_, err := u.userRepository.GetById(id)
	if err != nil {
		return errors.New("user not found")
	}

	// Realizando la actualización
	return u.userRepository.Update(id, user.Name, user.Email)
}

// Obteniendo la lista completa de usuarios desde el repositorio
func (u UserService) GetAll() ([]entities.User, error) {
	return u.userRepository.GetAll()
}

// Obteniendo un usuario por su ID, validando que no esté vacío
func (u UserService) GetById(id string) (entities.User, error) {
	if id == "" {
		return entities.User{}, errors.New("id cannot be empty")
	}
	return u.userRepository.GetById(id)
}

// Eliminando un usuario por su ID, validando que no esté vacío
func (u UserService) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return u.userRepository.Delete(id)
}
