package csv

import (
	"encoding/csv"
	"errors"
	"layersapi/entities"
	"os"
	"time"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u UserRepository) GetAll() ([]entities.User, error) {
	file, err := os.Open(`C:\Users\david\OneDrive\Documents\catedras\ncapas-01-25\repo\layers-backend\data\data.csv`)
	if err != nil {
		return []entities.User{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return []entities.User{}, err
	}

	var result []entities.User

	for i, record := range records {
		if i == 0 {
			continue
		}

		createdAt, _ := time.Parse(time.RFC3339, record[3])
		updatedAt, _ := time.Parse(time.RFC3339, record[4])
		meta := entities.Metadata{
			CreatedAt: createdAt.String(),
			UpdatedAt: updatedAt.String(),
			CreatedBy: record[5],
			UpdatedBy: record[6],
		}
		result = append(result, entities.NewUser(record[0], record[1], record[2], meta))
	}

	return result, nil
}

func (u UserRepository) GetById(id string) (entities.User, error) {
	file, err := os.Open("../../../data/data.csv")
	if err != nil {
		return entities.User{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return entities.User{}, err
	}

	for i, record := range records {
		if i == 0 {
			continue
		} else if record[0] == id {

			createdAt, _ := time.Parse(time.RFC3339, record[3])
			updatedAt, _ := time.Parse(time.RFC3339, record[4])
			meta := entities.Metadata{
				CreatedAt: createdAt.String(),
				UpdatedAt: updatedAt.String(),
				CreatedBy: record[5],
				UpdatedBy: record[6],
			}
			return entities.NewUser(record[0], record[1], record[2], meta), nil
		}

	}

	return entities.User{}, errors.New("user not found")
}

func (u UserRepository) Create(user entities.User) error {
	file, err := os.Open("data.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	newUser := []string{
		user.Id,
		user.Name,
		user.Email,
		user.Metadata.CreatedAt,
		user.Metadata.UpdatedAt,
		"webapp",
		"webapp",
	}

	if err := writer.Write(newUser); err != nil {
		return err
	}

	return nil
}

func (u UserRepository) Update(id, name, email string) error {
	//TODO
	return nil
}
