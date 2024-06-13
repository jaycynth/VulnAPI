package request

type UpdateUserRequest struct {
	Id       int    `validate:"required"`
	Username string `validate:"required,max=200,min=1" json:"name"`
}
