package interfaces

import "gitverse-internship-zg/services/user-service/domain/entities"

type EntityRepository interface {
	GetByID(id string) (*entities.Entity, error)
	EditByID(newUser entities.EntityRequest) (*entities.Entity, error)
	CreateUser(userData entities.EntityRequest) (*entities.Entity, error)
}
