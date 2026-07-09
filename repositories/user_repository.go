package repositories

import (
	"todo-api/models"

	"github.com/beego/beego/v2/client/orm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
}

type postgresUserRepository struct{}

func NewUserRepository() UserRepository {
	return &postgresUserRepository{}
}

func (r *postgresUserRepository) Create(user *models.User) error {
	o := orm.NewOrm()
	_, err := o.Insert(user)
	return err
}

func (r *postgresUserRepository) GetByID(id int) (*models.User, error) {
	o := orm.NewOrm()
	user := &models.User{ID: id}
	err := o.Read(user)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) GetByEmail(email string) (*models.User, error) {
	o := orm.NewOrm()
	user := &models.User{}
	err := o.QueryTable(new(models.User)).Filter("email", email).One(user)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
