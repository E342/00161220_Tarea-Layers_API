package csv

import (
	"encoding/csv"
	"errors"
	"layersapi/entities"
	"os"
	"time"
)

const filePath = "data/data.csv"

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u UserRepository) GetAll() ([]entities.User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []entities.User
	for i, record := range records {
		if i == 0 {
			continue
		}

		createdAt, _ := time.Parse(time.RFC3339, record[3])
		updatedAt, _ := time.Parse(time.RFC3339, record[4])
		meta := entities.Metadata{
			CreatedAt: createdAt.Format(time.RFC3339),
			UpdatedAt: updatedAt.Format(time.RFC3339),
			CreatedBy: record[5],
			UpdatedBy: record[6],
		}
		result = append(result, entities.NewUser(record[0], record[1], record[2], meta))
	}

	return result, nil
}

func (u UserRepository) GetById(id string) (entities.User, error) {
	file, err := os.Open(filePath)
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
		}
		if record[0] == id {
			createdAt, _ := time.Parse(time.RFC3339, record[3])
			updatedAt, _ := time.Parse(time.RFC3339, record[4])
			meta := entities.Metadata{
				CreatedAt: createdAt.Format(time.RFC3339),
				UpdatedAt: updatedAt.Format(time.RFC3339),
				CreatedBy: record[5],
				UpdatedBy: record[6],
			}
			return entities.NewUser(record[0], record[1], record[2], meta), nil
		}
	}

	return entities.User{}, errors.New("user not found")
}

func (u UserRepository) Create(user entities.User) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	newRecord := []string{
		user.Id,
		user.Name,
		user.Email,
		user.Metadata.CreatedAt,
		user.Metadata.UpdatedAt,
		user.Metadata.CreatedBy,
		user.Metadata.UpdatedBy,
	}

	return writer.Write(newRecord)
}

// Update actualiza un usuario existente en el archivo CSV.
// La función busca el ID del usuario en el archivo CSV,
// La función utiliza la biblioteca "encoding/csv" para leer y escribir en el archivo CSV.
func (u UserRepository) Update(id, name, email string) error {
	// Abriendo el archivo CSV en modo lectura
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Leyendo todos los registros del archivo CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	updated := false
	for i, record := range records {
		if i == 0 {
			continue
		}
		if record[0] == id {
			records[i][1] = name
			records[i][2] = email
			records[i][4] = time.Now().Format(time.RFC3339) // updatedAt (fecha de actualización)
			records[i][6] = "webapp"                        // updatedBy (Fuente de la actualización)
			updated = true
			break
		}
	}

	if !updated {
		return errors.New("user not found")
	}

	// Reabriendo el archivo en modo escritura para sobrescribirlo
	file, err = os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Creando un nuevo escritor CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribirndo todos los registros actualizados en el archivo CSV
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// Delete elimina un usuario del archivo CSV si su ID coincide.
// La función lee todo el archivo, encuenta el usuario correspondiente
// y sobrescribe sin él.
func (u UserRepository) Delete(id string) error {
	// Abriendo el archivo CSV en modo lectura
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Leyendo todos los registros del archivo CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var updatedRecords [][]string // Guardando los registros que se conservarán
	found := false

	for _, record := range records {
		// Si el ID no coincide, se mantiene el registro
		if len(record) > 0 && record[0] != id {
			updatedRecords = append(updatedRecords, record)
		} else if record[0] == id {
			// Marcando el usuario encontrado
			found = true
		}
	}

	// Si no se encontró el usuario, retorna error
	if !found {
		return errors.New("user not found")
	}

	// Reescribiendo el archivo sin el usuario eliminado
	file, err = os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribiendo los registros restantes
	for _, record := range updatedRecords {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
