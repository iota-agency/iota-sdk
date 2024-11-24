package upload

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"path/filepath"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/iota-agency/iota-sdk/pkg/configuration"
	"github.com/iota-agency/iota-sdk/pkg/constants"
)

type CreateDTO struct {
	File io.ReadSeeker `validate:"required"`
	Name string        `validate:"required"`
	Size int           `validate:"required"`
	Type string
}

type UpdateDTO struct {
}

func (d *CreateDTO) Ok(l ut.Translator) (map[string]string, bool) {
	errorMessages := map[string]string{}
	errs := constants.Validate.Struct(d)
	if errs == nil {
		return errorMessages, true
	}

	for _, err := range errs.(validator.ValidationErrors) {
		errorMessages[err.Field()] = err.Translate(l)
	}
	return errorMessages, len(errorMessages) == 0
}

func (d *UpdateDTO) Ok(l ut.Translator) (map[string]string, bool) {
	errorMessages := map[string]string{}
	errs := constants.Validate.Struct(d)
	if errs == nil {
		return errorMessages, true
	}

	for _, err := range errs.(validator.ValidationErrors) {
		errorMessages[err.Field()] = err.Translate(l)
	}
	return errorMessages, len(errorMessages) == 0
}

func (d *CreateDTO) ToEntity() (*Upload, []byte, error) {
	conf := configuration.Use()
	bytes, err := io.ReadAll(d.File)
	if err != nil {
		return nil, nil, nil
	}
	mdsum := md5.Sum(bytes)
	hash := hex.EncodeToString(mdsum[:])
	return &Upload{
		ID:        hash,
		Name:      d.Name,
		Size:      d.Size,
		Mimetype:  *mimetype.Detect(bytes),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		URL:       conf.Origin + filepath.Join("/", conf.UploadsPath, hash),
	}, bytes, nil
}

func (d *UpdateDTO) ToEntity(id string) (*Upload, error) {
	return &Upload{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}