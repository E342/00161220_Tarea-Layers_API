package dto

type MetadataDTO struct {
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}

type CreateUser struct {
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Metadata MetadataDTO `json:"metadata"`
}

type UpdateUser struct {
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Metadata MetadataDTO `json:"metadata"`
}
