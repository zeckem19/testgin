package service

import "github.com/zeckem19/testgin/entity"

type PersonService interface {
	Save(entity.Person) entity.Person
	FindAll() []entity.Person
}
