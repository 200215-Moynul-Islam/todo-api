package repositories

import (
	"todo-api/models"

	"github.com/beego/beego/v2/client/orm"
)

type TaskRepository interface {
	Create(task *models.Task) error
	GetAll(userID int, status string, page, limit int) ([]models.Task, error)
	GetByID(id int) (*models.Task, error)
	Update(task *models.Task) error
	Delete(id int) (bool, error)
}

type postgresTaskRepository struct{}

func NewTaskRepository() TaskRepository {
	return &postgresTaskRepository{}
}

func (r *postgresTaskRepository) Create(task *models.Task) error {
	o := orm.NewOrm()
	_, err := o.Insert(task)
	return err
}

func (r *postgresTaskRepository) GetAll(userID int, status string, page, limit int) ([]models.Task, error) {
	o := orm.NewOrm()
	var tasks []models.Task

	qs := o.QueryTable(new(models.Task)).Filter("user_id", userID)

	// Apply filtering if provided
	if status != "" {
		qs = qs.Filter("status", status)
	}

	// Apply offset pagination calculation
	offset := (page - 1) * limit
	_, err := qs.Limit(limit, offset).All(&tasks)

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *postgresTaskRepository) GetByID(id int) (*models.Task, error) {
	o := orm.NewOrm()
	task := &models.Task{ID: id}

	err := o.Read(task)
	if err == orm.ErrNoRows {
		return nil, nil // Return nil task cleanly to indicate Not Found
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *postgresTaskRepository) Update(task *models.Task) error {
	o := orm.NewOrm()
	_, err := o.Update(task, "title", "description", "status", "updated_at")
	return err
}

func (r *postgresTaskRepository) Delete(id int) (bool, error) {
	o := orm.NewOrm()
	task := &models.Task{ID: id}

	num, err := o.Delete(task)
	if err != nil {
		return false, err
	}
	return num > 0, nil
}
