package entities

type Account struct {
	ID          string `validate:"required"`
	Username    string `validate:"required,lowercase,gte=3,lte=12"`
	Name        string `validate:"required,gte=3,lte=16"`
	Description string `validate:"required,lte=140"`
	Email       string `validate:"required,email"`
	Password    string `validate:"required,lowercase,gte=6,lte=15"`
	CreatedAt   string
	UpdatedAt   string
	Deleted     bool
}
