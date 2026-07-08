package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Task struct {
	ID          int       `orm:"column(id);auto;pk" json:"id"`
	Title       string    `orm:"column(title);size(255)" json:"title"`
	Description string    `orm:"column(description);type(text);null" json:"description"`
	Status      string    `orm:"column(status);size(20)" json:"status"`
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(timestamp)" json:"created_at"`
	UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(timestamp)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(Task))
}