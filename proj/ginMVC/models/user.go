package models

import (
	"github.com/davidddw2017/panzer/proj/ginMvc/drivers"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB = drivers.MysqlDb

type User struct {
	Id      int    `json:"id" form:"id" primaryKey:"true"`
	Name    string `json:"name" form:"name" binding:"required"`
	Address string `json:"address" form:"address" binding:"required"`
	Age     int    `json:"age" form:"age" binding:"required"`
}

// get one
func (model *User) UserGet(id int) (user User, err error) {
	// find one record
	if err = db.Table("user").Where("id = ?", id).First(&user).Error; err != nil {
		return
	}
	return
}

// get list
func (model *User) UserGetList(page int, pageSize int) (users []User, err error) {
	users = make([]User, 0)
	offset := pageSize * (page - 1)
	limit := pageSize

	if err = db.Table("user").Find(&users).Limit(limit).Offset(offset).Error; err != nil {
		return
	}
	return
}

// create
func (model *User) UserAdd() (id int64, err error) {
	//result, err := db.Exec("INSERT INTO `users`(`name`, `age`) VALUES (?, ?)", model.Name, model.Age)
	user := User{Name: model.Name, Age: model.Age, Address: model.Address}
	if err = db.Table("user").Create(&user).Error; err != nil {
		return
	}
	id = int64(user.Id)
	return
}

// update
func (model *User) UserUpdate(id int) (afr int64, err error) {
	user := User{Id: id, Name: model.Name, Age: model.Age, Address: model.Address}
	executeAction := db.Table("user").Save(&user)
	afr, err = executeAction.RowsAffected, executeAction.Error
	return
}

// delete
func (model *User) UserDelete(id int) (afr int64, err error) {
	executeAction := db.Table("user").Delete(&User{}, "id = ?", id)
	afr, err = executeAction.RowsAffected, executeAction.Error
	return
}
