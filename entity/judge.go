package entity

import (
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/dtos"
	"github.com/arravoco/hackathon_backend/utils"
)

type Judge struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	passwordHash string
	Gender       string `json:"gender"`
	State        string `json:"state"`
	Role         string `json:"role"`
}

func (p *Judge) Register(input dtos.RegisterNewJudgeDTO) (*data.CreateJudgeAccountData, error) {
	passwordHash, err := utils.GenerateHashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	dataInput := &data.CreateJudgeAccountData{
		CreateAccountData: data.CreateAccountData{
			Email:        input.Email,
			PasswordHash: passwordHash,
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			Gender:       input.Gender,
			State:        input.State,
			Role:         "JUDGE"},
	}
	dataResponse, err := data.CreateJudgeAccount(dataInput)
	// emit created event

	return dataResponse, err
}
