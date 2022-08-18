package api

import (
	"strings"
	"time"

	"github.com/bytebase/bytebase/api"
	"github.com/pkg/errors"
)

// validPlans is a string array of valid plan types.
var validPlans = []api.PlanType{
	api.TEAM,
	api.ENTERPRISE,
}

// License is the API message for enterprise license.
type License struct {
	Subject       string
	InstanceCount int
	ExpiresTs     int64
	IssuedTs      int64
	Plan          api.PlanType
	Trialing      bool
	OrgName       string
}

// Valid will check if license expired or has correct plan type.
func (l *License) Valid() error {
	if expireTime := time.Unix(l.ExpiresTs, 0); expireTime.Before(time.Now()) {
		return errors.Errorf("license has expired at %v", expireTime)
	}

	return l.validPlanType()
}

func (l *License) validPlanType() error {
	for _, plan := range validPlans {
		if plan == l.Plan {
			return nil
		}
	}

	return errors.Errorf("plan %q is not valid, expect %s or %s",
		l.Plan.String(),
		api.TEAM.String(),
		api.ENTERPRISE.String(),
	)
}

// OrgID extract the organization id from license subject.
func (l *License) OrgID() string {
	return strings.Split(l.Subject, ".")[0]
}

// LicenseService is the service for enterprise license.
type LicenseService interface {
	// StoreLicense will store license into file.
	StoreLicense(patch *SubscriptionPatch) error
	// LoadLicense will load license from file and validate it.
	LoadLicense() (*License, error)
}
