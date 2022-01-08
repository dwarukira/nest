package request

import (
	"time"

	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/model"
)

type CreateLeaseRequest struct {
	StartDate         string `json:"start_date"`
	MonthlyRent       int    `json:"rent"`
	SecurityDeposit   int    `json:"security_deposit"`
	UnitID            string `json:"unit_id"`
	RentDueDayOfMonth int    `json:"rent_due_day_of_month"`
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

	return model.Lease{
		StartDate:         t,
		MonthlyRent:       createLeaseRequest.MonthlyRent,
		SecurityDeposit:   createLeaseRequest.SecurityDeposit,
		UnitID:            unitID,
		RentDueDayOfMonth: createLeaseRequest.RentDueDayOfMonth,
	}, nil
}
