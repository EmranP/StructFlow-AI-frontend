package dto

type ProjectRequest struct {
	Title string `json:"title" validate:"required"`

	ProjectType string `json:"projectType"`

	Stack string `json:"stack"`

	Architecture string `json:"architecture"`

	Features string `json:"features"`

	AdditionalInfo string `json:"additionalInfo"`

	Prompt string `json:"prompt" validate:"required,min=10"`
}
