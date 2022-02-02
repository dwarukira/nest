package request

import (
	"time"

	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/logger"
	"github.com/solabsafrica/afrikanest/model"
)

// swagger:parameters createLease
type CreateLeaseRequestParam struct {
	// in: body
	// required: true
	Body CreateLeaseRequest
}

// swagger:parameters getLease
type LeaseRequest struct {
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:parameters getLeaseCharge
type LeaseChargeRequest struct {
	// in: path
	ID string `json:"id"`

	// in: path
	LeaseChargeID string `json:"leaseChargeId"`
}

// swagger:parameters createLeaseChargePayment
type CreateLeaseChargePaymentRequest struct {
	// in: body
	// required: true
	Body CreateLeaseChargePaymentRequestBody

	// in: path
	LeaseChargeID string `json:"leaseChargeId"`

	// in: path
	ID string `json:"id"`
}

//  swagger:parameters getLeaseBalance
type LeaseBalanceRequest struct {
	// in: path
	ID string `json:"id"`
}

type CreateLeaseChargePaymentRequestBody struct {
	// in: body
	// required: true
	Amount            int64  `json:"amount"`
	PaymentDate       string `json:"paymentDate"`
	SentEmailToTenant bool   `json:"sentEmailToTenant"`
}

func (r CreateLeaseChargePaymentRequestBody) Validate() error {
	logger.Info("validating lease charge payment request", r.Amount)
	if r.Amount == 0 {
		return exceptions.LeaseChargePaymentCreateFaild.SetMessage("amount must be greater than zero")
	}
	_, err := time.Parse("2006-01-02", r.PaymentDate)
	if err != nil {
		return exceptions.LeaseCreateFaild.SetMessage("invalid date")
	}

	return nil
}

func (r CreateLeaseChargePaymentRequestBody) ToModel(leaseId uuid.UUID) (model.LeaseChargePayment, error) {
	if err := r.Validate(); err != nil {
		return model.LeaseChargePayment{}, err
	}

	t, _ := time.Parse("2006-01-02", r.PaymentDate)

	return model.LeaseChargePayment{
		Amount:        r.Amount,
		LeaseChargeID: leaseId,
		PaymentDate:   t,
	}, nil
}

type CreateLeaseRequest struct {
	StartDate         string                `json:"start_date"`
	MonthlyRent       int                   `json:"rent"`
	SecurityDeposit   int                   `json:"security_deposit"`
	UnitID            string                `json:"unit_id"`
	RentDueDayOfMonth int                   `json:"rent_due_day_of_month"`
	Tenants           []CreateTenantRequest `json:"tenants"`
	InviteTenants     bool                  `json:"invite"`
}

func (createLeaseRequest CreateLeaseRequest) Validate() error {
	if len(createLeaseRequest.StartDate) == 0 {
		return exceptions.LeaseCreateFaild.SetMessage("start_date is empty")
	}
	if createLeaseRequest.RentDueDayOfMonth == 0 {
		return exceptions.LeaseCreateFaild.SetMessage("rent_due_day_of_month is empty")
	}

	if createLeaseRequest.MonthlyRent == 0 {
		return exceptions.LeaseCreateFaild.SetMessage("rent is empty")
	}

	if len(createLeaseRequest.Tenants) != 0 {
		for _, tenant := range createLeaseRequest.Tenants {
			if err := tenant.Validate(); err != nil {
				return err
			}
		}
	}

	_, err := time.Parse("2006-01-02", createLeaseRequest.StartDate)
	if err != nil {
		return exceptions.LeaseCreateFaild.SetMessage("invalid date")
	}

	return nil
}

func (createLeaseRequest CreateLeaseRequest) ToLease() (model.Lease, error) {
	if err := createLeaseRequest.Validate(); err != nil {
		return model.Lease{}, err
	}
	t, _ := time.Parse("2006-01-02", createLeaseRequest.StartDate)
	unitID, err := uuid.Parse(createLeaseRequest.UnitID)
	if err != nil {
		return model.Lease{}, exceptions.LeaseCreateFaild.SetMessage("invalid unit_id")
	}

	var tenants []model.Tenant

	for _, tenant := range createLeaseRequest.Tenants {
		if err := tenant.Validate(); err != nil {
			return model.Lease{}, err
		}
		t, err := tenant.ToTenant()
		if err != nil {
			return model.Lease{}, err
		}
		tenants = append(tenants, t)

	}

	return model.Lease{
		StartDate:         t,
		MonthlyRent:       createLeaseRequest.MonthlyRent,
		SecurityDeposit:   createLeaseRequest.SecurityDeposit,
		UnitID:            unitID,
		RentDueDayOfMonth: createLeaseRequest.RentDueDayOfMonth,
		Tenants:           &tenants,
	}, nil
}

// swagger:parameters createLeaseCharge
type CreateLeaseChargeRequest struct {

	// in: path
	ID string `json:"id"`

	// in: body
	// required: true
	Body CreateLeaseChargeRequestBody
}

type CreateLeaseChargeRequestBody struct {
	// in: body
	// required: true
	Amount int64 `json:"amount"`
	// in: body
	// required: true
	Description string `json:"description"`
	// in: body
	// required: true
	DueDate string `json:"due_date"`

	// in: body
	// required: true
	NotifyTenant bool `json:"notify_tenant"`

	// in: body
	// required: true
	ChargeType model.LeaseChargeTypeEmun `json:"charge_type"` //enum: "rent" "security_deposit" "other"

}

func (r CreateLeaseChargeRequestBody) Validate() error {
	if r.Amount == 0 {
		return exceptions.LeaseChargeCreateFaild.SetMessage("amount must be greater than zero")
	}
	_, err := time.Parse("2006-01-02", r.DueDate)
	if err != nil {
		return exceptions.LeaseCreateFaild.SetMessage("invalid date")
	}

	return nil
}

func (r CreateLeaseChargeRequestBody) ToModel(leaseId uuid.UUID) (model.LeaseCharge, error) {
	if err := r.Validate(); err != nil {
		return model.LeaseCharge{}, err
	}

	t, _ := time.Parse("2006-01-02", r.DueDate)

	return model.LeaseCharge{
		Amount:      r.Amount,
		Description: r.Description,
		DueDate:     t,
		ChargeType:  r.ChargeType,
		LeaseID:     leaseId,
	}, nil
}
