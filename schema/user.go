package schema

const (
	StatusActive      = "ACTIVE"
	StatusDeactivated = "DEACTIVATED"
)

type UserRequest struct {
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	UserEmail      *string `json:"user_email,omitempty"`
	EmployeeNumber *string `json:"employee_number,omitempty"`
	OnboardTime    *int    `json:"onboard_time,omitempty"`
	Status         *string `json:"status,omitempty"`
}

type UserResponse struct {
	Id              string         `json:"Id"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	FullName        string         `json:"full_name"`
	Alias           string         `json:"alias"`
	UserEmail       string         `json:"user_email"`
	EmailStatus     string         `json:"email_status"`
	Phone           string         `json:"phone"`
	EmployeeNumber  string         `json:"employee_number"`
	OnboardTime     int            `json:"onboard_time"`
	NfcCards        []NfcCard      `json:"nfc_cards"`
	PinCode         PinCode        `json:"pin_code"`
	AccessPolicyIds []string       `json:"access_policy_ids"`
	AccessPolicies  []AccessPolicy `json:"access_policies"`

	// TODO make this an enum
	Status string `json:"status"`
}
