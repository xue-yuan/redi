package models

type UpdateOpenGraphRequest struct {
	OpenGraph
	NewImage string `json:"new_image" validate:"exist"`
}
