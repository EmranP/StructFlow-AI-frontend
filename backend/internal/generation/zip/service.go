package zip

import (
	"archive/zip"
	"bytes"
	"encoding/json"
)

type Service struct {
}

type TemplateContent struct {
	Name string `json:"name"`

	Directories []string `json:"directories"`

	Files []string `json:"files"`
}

func New() *Service {
	return &Service{}
}

func (s *Service) Generate(
	content []byte,
) ([]byte, error) {
	var tpl TemplateContent

	err := json.Unmarshal(
		content,
		&tpl,
	)

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	writer := zip.NewWriter(buf)

	for _, dir := range tpl.Directories {

		_, err := writer.Create(
			dir + "/",
		)

		if err != nil {
			return nil, err
		}
	}

	for _, file := range tpl.Files {

		f, err := writer.Create(
			file,
		)

		if err != nil {
			return nil, err
		}

		_, err = f.Write(
			[]byte(""),
		)

		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
