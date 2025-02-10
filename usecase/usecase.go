package usecase

import (
	"gitverse-internship-zg/services/user-service/domain/entities"
	"gitverse-internship-zg/services/user-service/domain/interfaces"

	"github.com/sirupsen/logrus"
)

type Usecase struct {
	Repo interfaces.EntityRepository
	MesQ interfaces.EntityMessageQ
}

const (
	userCreated = "user_created"
	userEdited  = "user_edited"
	userTopic   = "user_updates"
	// userCreated = "user_created"
)

func (u *Usecase) GetEntityByID(id string) (*entities.Entity, error) {
	return u.Repo.GetByID(id)
}

func (u *Usecase) EditUserByID(newUser entities.EntityRequest) (*entities.Entity, error) {
	res, err := u.Repo.EditByID(newUser)
	if err == nil {
		if err = u.MesQ.Send(userEdited, userTopic); err != nil {
			logrus.Error(err)
		}
	}
	return res, err
}

func (u *Usecase) CreateUser(userData entities.EntityRequest) (*entities.Entity, error) {
	return u.Repo.CreateUser(userData)
}

func NewUsecase(repo interfaces.EntityRepository, mesq interfaces.EntityMessageQ) *Usecase {
	return &Usecase{Repo: repo, MesQ: mesq}
}
