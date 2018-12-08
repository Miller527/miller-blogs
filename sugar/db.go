//
// __author__ = "Miller"
// Date: 2018/11/25
//

package sugar

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"miller-blogs/settings"
)

// database table list
var tables []string

// connection database config
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

// join connection source name
func (dc *DBConfig) dataSourceName() string {
	return fmt.Sprintf("%s:%s@%s(%s)/%s", dc.Username, dc.Password, dc.Protocol, dc.Address, dc.DBName)
}

// database manager
var Dbm DBManager

type DBManager struct {
	Conf DBConfig
	Db   *sql.DB
}

// database manager init
func (dbm *DBManager) init() {
	db, err := sql.Open(dbm.Conf.Driver, dbm.Conf.dataSourceName())
	if err != nil {
		panic(err)
	}

	if dbm.Conf.MaxOpenConns == 0 {
		dbm.Conf.MaxOpenConns = 100
	}
	db.SetMaxOpenConns(dbm.Conf.MaxOpenConns)
	if dbm.Conf.MaxIdleConns != 0 {
		dbm.Conf.MaxOpenConns = 10
	}
	db.SetMaxIdleConns(dbm.Conf.MaxOpenConns)

	if err := db.Ping(); err != nil {
		panic(err)
	}
	dbm.Db = db
}

func (dbm *DBManager) showTables() {
	stmt, _ := dbm.Db.Prepare("SHOW TABLES")
	type showTables struct {
		name string
	}
	tc := &TableConf{
		Field: []string{"name"},
		Title: []string{"数据表"},
		Desc:  &showTables{},
	}
	result, err := dbm.SelectSlice(stmt, tc)
	if err != nil {
		fmt.Println("tables error")
	}
	for _, line := range result {
		tables = append(tables, line[0])
	}
	fmt.Println(tables)
}

//// todo 优化其他字段类型
//func (dbm *DBManager) dest(tc *TableConf) ([]interface{}, error) {
//	// todo 其他类型时候反射取值
//	value := reflect.ValueOf(tc.Desc)
//	if value.Kind() != reflect.Ptr {
//		return nil, errors.New("KindError: TableConf.Desc error")
//	}
//	elem := value.Elem()
//
//	var dest []interface{}
//	for _, field := range tc.Field {
//		elemField := elem.FieldByName(field)
//		var err error
//		switch elemField.Kind() {
//		case reflect.String:
//			// todo 通过判断取值
//			dest = append(dest, &*(*string)(unsafe.Pointer(elemField.Addr().Pointer())))
//		case reflect.Int:
//			dest = append(dest, &*(*int)(unsafe.Pointer(elemField.Addr().Pointer())))
//		case reflect.Bool:
//			dest = append(dest, &*(*bool)(unsafe.Pointer(elemField.Addr().Pointer())))
//		case reflect.Uint:
//			dest = append(dest, &*(*uint)(unsafe.Pointer(elemField.Addr().Pointer())))
//		case reflect.Float64:
//			dest = append(dest, &*(*float64)(unsafe.Pointer(elemField.Addr().Pointer())))
//		case reflect.Float32:
//			dest = append(dest, &*(*float32)(unsafe.Pointer(elemField.Addr().Pointer())))
//			//todo 这里取值
//			//case reflect.Slice:
//			//dest = append(dest, &*(*)(unsafe.Pointer(elemField.Addr().Pointer())))
//		default:
//			err = errors.New("ValueError: filed kind error")
//			fmt.Println(elemField.Kind())
//		}
//		fmt.Println("xxxxxxxxxxxxxxxx")
//		if err != nil {
//			return nil, err
//		}
//	}
//	return dest, nil
//}

// todo 优化其他字段类型
//func (dbm *DBManager) line(tc *TableConf) ([]interface{}, error) {
//	value := reflect.ValueOf(tc.Desc)
//	if value.Kind() != reflect.Ptr {
//		return nil, errors.New("KindError: TableConf.Desc error")
//	}
//	elem := value.Elem()
//
//	var line []interface{}
//	for _, field := range tc.Field {
//		elemField := elem.FieldByName(field)
//		var err error
//		switch elemField.Kind() {
//		case reflect.String:
//			// todo 通过判断取值
//			line = append(line, elemField.String())
//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//			line = append(line, elemField.Int())
//		case reflect.Bool:
//			line = append(line, elemField.Bool())
//		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//			line = append(line, elemField.Uint())
//		case reflect.Float32, reflect.Float64:
//			line = append(line, elemField.Float())
//		default:
//			err = errors.New("ValueError: xxxxxxxxxxxxxxxxxx")
//		}
//		if err != nil {
//			return nil, err
//		}
//	}
//	return line, nil
//}

func (dbm *DBManager) line(vals []sql.RawBytes) []string{
	var val string
	var resLine []string
	for _, col := range vals {
		// Here we can check if the value is nil (NULL value)
		if col == nil {
			val = "NULL"
		} else {
			val = string(col)
		}
		resLine = append(resLine, val)
	}
	return resLine
}

func (dbm *DBManager) SelectValues(stmt *sql.Stmt, tc *TableConf, args ...interface{})([]string, error){
	rows, err := dbm.selectQuery(stmt, args...)
	if err != nil {
		return nil, err
	}
	values, scanArgs, err := dbm.valuesScan(rows)
	if err != nil{
		return nil,err
	}
	var result []string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		line := dbm.line(values)
		result = append(result, line...)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil

}

func (dbm *DBManager) selectQuery(stmt *sql.Stmt,  args ...interface{})(*sql.Rows, error){
	return  stmt.Query(args...)
}


func (dbm *DBManager) valuesScan(rows *sql.Rows)([]sql.RawBytes, []interface{}, error){

	columns, err := rows.Columns()
	if err != nil {
		return nil,nil, err
	}
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	return values, scanArgs, nil

}

func (dbm *DBManager) SelectSlice(stmt *sql.Stmt, tc *TableConf, args ...interface{}) ([][]string, error) {
	rows, err := dbm.selectQuery(stmt, args...)
	if err != nil {
		return nil, err
	}
	values, scanArgs, err := dbm.valuesScan(rows)
	if err != nil{
		return nil,err
	}
	var result [][]string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil,err
		}
		line := dbm.line(values)

		result = append(result, line)
	}
	if err = rows.Err(); err != nil {
		return nil,err
	}


	return result, nil
}

// Create New DBModel
func NewDbm(dbc DBConfig) DBManager {
	return DBManager{Conf: dbc}
}

func DbmInit() {
	dbc := DBConfig{
	}
	v,e :=json.Marshal(settings.Settings.DBConfig)
	if e != nil{
		fmt.Println("config dbconfig error", e)
	}
	err := json.Unmarshal(v, &dbc)
	if err != nil{
		fmt.Println("json dbconfig",err)
	}
	fmt.Println(dbc)
	Dbm = NewDbm(dbc)
	Dbm.init()
	Dbm.showTables()
}
