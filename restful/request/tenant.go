package request

import (
	"net/mail"

	"github.com/dongri/phonenumber"
	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/model"
)

type CreateTenantRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	SendInvite  bool   `json:"send_invite"`
}

func (createTenantRequest CreateTenantRequest) Validate() error {
	if len(createTenantRequest.Email) > 0 {
		if !IsValidEmail(createTenantRequest.Email) {
			return exceptions.TenantCreateFaild.SetMessage("invalid email")
		}
	}
	if len(createTenantRequest.PhoneNumber) > 0 {
		// TODO: dynamically validate country code - KE is only supported for now
		number := phonenumber.Parse(createTenantRequest.PhoneNumber, "KE")
		if len(number) == 0 {
			return exceptions.TenantCreateFaild.SetMessage("invalid phone number")
		}
	}
	return nil
}

func (createTenantRequest CreateTenantRequest) ToTenant() (model.Tenant, error) {
	if err := createTenantRequest.Validate(); err != nil {
		return model.Tenant{}, err
	}
	return model.Tenant{
		FirstName:   createTenantRequest.FirstName,
		LastName:    createTenantRequest.LastName,
		Email:       createTenantRequest.Email,
		PhoneNumber: createTenantRequest.PhoneNumber,
	}, nil
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
