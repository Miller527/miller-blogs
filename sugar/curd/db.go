//
// __author__ = "Miller"
// Date: 2018/11/25
//

package curd

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"unsafe"
)

var tables []string

var Dbm DBModels

type DBModels struct {
	Conf DBConfig
	Db *sql.DB
}

func (dbm *DBModels) init() {
	db, err := sql.Open(dbm.Conf.Driver, dbm.Conf.dataSourceName())
	if err != nil {
		panic(err)
	}

	if dbm.Conf.MaxOpenConns == 0 {
		dbm.Conf.MaxOpenConns=100
	}
	db.SetMaxOpenConns(dbm.Conf.MaxOpenConns)
	if dbm.Conf.MaxIdleConns != 0 {
		dbm.Conf.MaxOpenConns=10
	}
	db.SetMaxIdleConns(dbm.Conf.MaxOpenConns)

	if err := db.Ping(); err != nil {
		panic(err)
	}
	dbm.Db = db
}

func (dbm *DBModels) SelectSlice(stmt *sql.Stmt, tc *TableConf, args ...interface{}) ([]string,error){
	rows, err := stmt.Query(args...);
	if err != nil{
		fmt.Println("---", err)
		return nil, err
	}

	//value:=reflect.ValueOf(tc.Desc)
	//if value.Kind()==reflect.Ptr {
	//	elem := value.Elem()
	//	name := elem.FieldByName("name")
	//	if name.Kind() == reflect.String {
	//		*(*string)(unsafe.Pointer(name.Addr().Pointer())) = "fangwendong"
	//	}
	//}

	var result []string
	for rows.Next() {
		// todo 其他类型时候反射取值
		value:=reflect.ValueOf(tc.Desc)
		if value.Kind()==reflect.Ptr {
			elem := value.Elem()
			name := elem.FieldByName("name")
			//fmt.Println("name",name)
			if name.Kind() == reflect.String {
				*(*string)(unsafe.Pointer(name.Addr().Pointer())) = "fangwendong"
				// todo 通过判断取值
				err := rows.Scan(&*(*string)(unsafe.Pointer(name.Addr().Pointer())))
				if err == nil{
					fmt.Println("yyyyyyyyyyyyyyyyyyy",name)
					//result = append(result, name)

				}
			}
		}
	}
	return result,nil
}


func (dbm *DBModels) showTables(){
	stmt, _ := dbm.Db.Prepare("show tables")
	type tb struct {
		name string
	}
	x := &TableConf{
		Slice:[]interface{}{"id","rid","name"},
		Title:[]string{"ID","角色ID","角色名称"} ,
		Desc: &tb{},
	}


	tables, _ = dbm.SelectSlice(stmt, x)
	fmt.Println(tables)
}

func NewDbm(dc DBConfig) DBModels {


	return DBModels{Conf:dc}
}

func DbmInit(dc DBConfig) {
	Dbm = NewDbm(dc)
	Dbm.init()
	Dbm.showTables()
}

type DBConfig struct {
	Driver       string
	Username     string
	Password     string
	Protocol     string
	Address      string
	DBName       string
	Parameters   map[string]interface{}
	MaxOpenConns int
	MaxIdleConns int
}

func (dc *DBConfig) dataSourceName() string {
	return fmt.Sprintf("%s:%s@%s(%s)/%s", dc.Username, dc.Password, dc.Protocol, dc.Address, dc.DBName)
}
