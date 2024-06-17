package bindme

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJson = "application/json"
)

var (
	ErrEOF           = errors.New("body must not be empty")
	ErrInvalidJson   = errors.New("body contains badly-formed JSON")
	ErrDuplicateJson = errors.New("body contains only one JSON object")
)

func ReadJson(r *http.Request, dst interface{}) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return ErrInvalidJson
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return ErrEOF
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	if err = d.Decode(&struct{}{}); err != io.EOF {
		return ErrDuplicateJson
	}
	v := validator.New()
	return v.Struct(dst)
}

func ReadForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	f := form.NewDecoder()
	if err := f.Decode(dst, r.Form); err != nil {
		return err
	}
	v := validator.New()
	return v.Struct(dst)
}

func ReadFile(r *http.Request, fileName string, maxFileSize int64) (multipart.File, *multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(maxFileSize << 20); err != nil {
		return nil, nil, err
	}
	file, handler, err := r.FormFile(fileName)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, nil, fmt.Errorf("missing required field: %s", fileName)
		}
		return nil, nil, err
	}
	defer file.Close()
	return file, handler, nil
}

func WriteJson(w http.ResponseWriter, status int, v interface{}, headers http.Header) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set(ContentType, ContentTypeJson)
	w.WriteHeader(status)
	w.Write(data)
	return nil
}
