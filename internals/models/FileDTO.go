package models

import (
	"github.com/google/uuid"
	"time"
)

type FileDTO struct {
	Id               uuid.UUID `json:"id"`
	FileName         string    `json:"file_name"`
	OriginalFileName string    `json:"original_file_name"`
	FilePath         string    `json:"file_path"`
	CreatedOn        time.Time `json:"created_on"`
	UpdatedOn        time.Time `json:"updateed_on"`
}
