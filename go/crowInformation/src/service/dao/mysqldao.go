package dao

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
)

type MysqlDao struct {
	o orm.Ormer
}

var mysqlDao MysqlDao

func (d *MysqlDao) Connect() {
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8", 30)
	mysqlDao.o = orm.NewOrm()
	fmt.Println("connect mysql success")
}

func (d *MysqlDao) CreateN(n interface{}){
	mysqlDao.o.Insert(n)
}

func (d *MysqlDao) GetNById(id string) interface{} {
	return nil
}

func (d *MysqlDao) DeleteNById(id string) {

}

func (d *MysqlDao) UpdateN(n interface{}) {

}
