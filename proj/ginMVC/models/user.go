package models

import (
	"fmt"

	"github.com/davidddw2017/panzer/proj/ginMvc/drivers"
	"github.com/go-xorm/xorm"
)

var db *xorm.Engine = drivers.MySQLDB

type User struct {
	Id      int    `json:"id" form:"id" primaryKey:"true"`
	Name    string `json:"name" form:"name" binding:"required"`
	Address string `json:"address" form:"address" binding:"required"`
	Age     int    `json:"age" form:"age" binding:"required"`
}

// get one
func (model *User) UserGet(id int) (user User, err error) {
	// find one record
	err = db.Table("user").Where("id = ?", id).Find(&user)
	return
}

// UserGetList get list
func (model *User) UserGetList(page int, pageSize int) (users []User, err error) {
	users = make([]User, 0)
	offset := pageSize * (page - 1)
	limit := pageSize
	fmt.Println("db= ", db)
	err = db.Table("user").Limit(limit, offset).Find(&users)
	fmt.Println(err)
	fmt.Println(users)
	return
}

// UserAdd create
func (model *User) UserAdd() (id int64, err error) {
	user := User{Name: model.Name, Age: model.Age, Address: model.Address}
	id, err = db.Table("user").Insert(&user)
	return
}

// UserUpdate update
func (model *User) UserUpdate(id int) (afr int64, err error) {
	user := User{Id: id, Name: model.Name, Age: model.Age, Address: model.Address}
	afr, err = db.Table("user").Update(&user)
	return
}

// UserDelete delete
func (model *User) UserDelete(id int) (afr int64, err error) {
	afr, err = db.Table("user").Delete(&User{Id: id})
	return
}
