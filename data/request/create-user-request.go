package request

type CreateUserRequest struct {
	Username string `validate:"required,min=1,max=200" json:"username"`
	Password string `validate:"required,min=1,max=200" json:"password"`
	Email    string `validate:"required,min=1,max=200" json:"email"`
}
