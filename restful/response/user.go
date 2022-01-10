package response

import (
	"strings"

	"github.com/solabsafrica/afrikanest/model"
)

type CreateUserResponse struct {
	ID string `json:"id"`
}

// swagger:model UpdateUserResponse
type UpdateUserResponse struct {
	ID        string `json:"id"`
	UpdatedAt string `json:"updated_at"`
}

// swagger:model GetUserResponse
type GetUserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type GetUsersResponse struct {
	Pagination Pagination   `json:"pagination"`
	Users      []model.User `json:"users"`
}

func NewCreateUserResponse(user model.User) CreateUserResponse {
	return CreateUserResponse{
		ID: user.ID.String(),
	}
}

func NewGetUserResponse(user model.User) GetUserResponse {
	return GetUserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      strings.Title(strings.ToLower(user.FirstName)) + " " + strings.Title(strings.ToLower(user.LastName)),
	}
}

func NewGetUsersResponse(users []model.User, pagination Pagination) GetUsersResponse {
	return GetUsersResponse{
		Users:      users,
		Pagination: pagination,
	}
}

func NewUpdateUserResponse(user model.User) UpdateUserResponse {
	return UpdateUserResponse{
		ID:        user.ID.String(),
		UpdatedAt: user.UpdatedAt.GoString(),
	}
}
