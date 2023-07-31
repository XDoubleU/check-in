package dtos

type ErrorDto struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message any    `json:"message"`
} //	@name	ErrorDto
