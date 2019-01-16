//
// __author__ = "Miller"
// Date: 2018/11/25
//

package sugar

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"miller-blogs/sugar/utils"
	"strconv"
)

// connection database config
type DBConfig struct {
	Driver       string
	Username     string
	Password     string
	Protocol     string
	Address      string
	DBName       string
	Params       map[string]string
	MaxOpenConns int
	MaxIdleConns int
}

func (dc *DBConfig) joinParams() string {
	params := dc.Params
	if params == nil {
		return ""
	}
	s := "?"
	lendth := len(params)
	count := 0
	for k, v := range params {
		count++
		s = s + k + "=" + v
		if count != lendth {
			s = s + "&"
		}
	}
	return s
}

// join connection source name
func (dc *DBConfig) source() string {
	paramStr := dc.joinParams()
	return fmt.Sprintf("%s:%s@%s(%s)/%s%s",
		dc.Username, dc.Password, dc.Protocol, dc.Address, dc.DBName, paramStr)
}

// change connection source name
func (dc *DBConfig) changeSourceName(name string) string {
	paramStr := dc.joinParams()
	return fmt.Sprintf("%s:%s@%s(%s)/%s%s",
		dc.Username, dc.Password, dc.Protocol, dc.Address, name, paramStr)
}

// database manager
var Dbm DBManager

type DBManager struct {
	Conf      DBConfig
	DBPool    map[string]*sql.DB
	DefaultDB *sql.DB
}

// 获取一个连接池对象
func (dbm *DBManager) GetDB(name string) *sql.DB {
	db, ok := dbm.DBPool[name]
	if ok {
		return db
	}
	return dbm.DefaultDB
}

// 创建一个新连接池对象
func (dbm *DBManager) NewDB(source string) (*sql.DB, error) {
	db, err := sql.Open(dbm.Conf.Driver, source)
	if err != nil {
		return nil, err
	}
	if dbm.Conf.MaxOpenConns == 0 {
		dbm.Conf.MaxOpenConns = 100
	}
	db.SetMaxOpenConns(dbm.Conf.MaxOpenConns)
	if dbm.Conf.MaxIdleConns != 0 {
		dbm.Conf.MaxOpenConns = 10
	}
	db.SetMaxIdleConns(dbm.Conf.MaxOpenConns)
	return db, err

}

// update database conn pool
func (dbm *DBManager) UpdateDBPool(name string) {
	_, ok := dbm.DBPool[name]
	if ! ok {
		db, err := dbm.NewDB(dbm.Conf.changeSourceName(name))
		utils.PanicCheck(err)
		err = db.Ping()
		utils.PanicCheck(err)
		dbm.DBPool[name] = db
	}
}

// close and delete database conn pool
func (dbm *DBManager) CloseDB(name string) {
	defer delete(dbm.DBPool, name)
	if db, ok := dbm.DBPool[name]; ok {
		db.Close()
	}
}

// database manager init
func (dbm *DBManager) init() {
	db, err := dbm.NewDB(dbm.Conf.source())
	utils.PanicCheck(err)
	err = db.Ping()
	utils.PanicCheck(err)
	dbm.DefaultDB = db
}

// 查询所有database列表
func (dbm *DBManager) showDatabase() []string {
	stmt, _ := dbm.DefaultDB.Prepare("SHOW DATABASES")
	result, err := dbm.SelectValues(stmt)
	utils.PanicCheck(err)
	return result
}

// 查询某database的table列表
func (dbm *DBManager) showTables(name string) []string {
	db, ok := dbm.DBPool[name]
	if !ok {
		errStr := fmt.Sprintf("DBMNotFoundPool: DBManager not found databases '%s' conn pool", name)
		utils.PanicCheck(errors.New(errStr))
	}
	stmt, _ := db.Prepare("SHOW TABLES")

	result, err := dbm.SelectValues(stmt)
	utils.PanicCheck(err)

	return result
}

// get raw line
func (dbm *DBManager) list(vals []sql.RawBytes) []string {
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
func (dbm *DBManager) dict(vals []sql.RawBytes, columns []string) map[string]interface{} {
	var val interface{}
	var resLine = make(map[string]interface{})
	for i, col := range vals {
		// Here we can check if the value is nil (NULL value)
		//switch col {
		//case nil:
		//	val = "NULL"
		//case []byte("false"):
		//	val = false
		//case []byte("true"):
		//	val = true
		//default:
		//	val = string(col)
		//	k, err := strconv.Atoi(val.(string))
		//	if err == nil{
		//		val = k
		//	}
		//}
		// todo 这里优化取值
		if col == nil {
			val = "NULL"
		} else {
			val = string(col)
			k, err := strconv.Atoi(val.(string))
			if err == nil {
				val = k
			}
		}
		resLine[columns[i]] = val

	}
	return resLine
}

// select query return type  []string
func (dbm *DBManager) SelectValues(stmt *sql.Stmt, args ...interface{}) ([]string, error) {
	rows, err := dbm.selectQuery(stmt, args...)
	if err != nil {
		return nil, err
	}
	values, scanArgs, err := dbm.valuesScan(rows)
	if err != nil {
		return nil, err
	}
	var result []string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		line := dbm.list(values)
		result = append(result, line...)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil

}

func (dbm *DBManager) selectQuery(stmt *sql.Stmt, args ...interface{}) (*sql.Rows, error) {
	fmt.Println(args, stmt == nil)
	return stmt.Query(args...)
}

func (dbm *DBManager) valuesScan(rows *sql.Rows) ([]sql.RawBytes, []interface{}, error) {

	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	return values, scanArgs, nil

}

// select query return type  [][]string
func (dbm *DBManager) SelectSlice(stmt *sql.Stmt, args ...interface{}) ([][]string, error) {
	fmt.Println(args)
	rows, err := dbm.selectQuery(stmt, args...)
	if err != nil {
		return nil, err
	}
	values, scanArgs, err := dbm.valuesScan(rows)
	if err != nil {
		return nil, err
	}
	var result [][]string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		line := dbm.list(values)

		result = append(result, line)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (dbm *DBManager) SelectDict(stmt *sql.Stmt, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := dbm.selectQuery(stmt, args...)
	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}
	values, scanArgs, err := dbm.valuesScan(rows)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		line := dbm.dict(values, columns)

		result = append(result, line)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// Create New DBModel
func NewDbm(dbc DBConfig) DBManager {
	return DBManager{Conf: dbc, DBPool: map[string]*sql.DB{}}
}

func DBMInit(dbConfig DBConfig) {
	//dbc := DBConfig{}
	//v, err := json.Marshal(dbConfig)
	//utils.PanicCheck(err)

	//err = json.Unmarshal(v, &dbc)
	//utils.PanicCheck(err)
	Dbm = NewDbm(dbConfig)
	Dbm.init()
}
