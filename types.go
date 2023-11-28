package main

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
type Tenant struct {
	ID          uint   `json:"-"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	TrialPeriod bool   `json:"periodo_teste"`
}
