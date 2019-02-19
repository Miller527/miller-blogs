/*
# __author__ = "Mr.chai"
# Date: 2019-02-18
*/
package sql_generater

import (
	"fmt"
	"strings"
)

type DriverType int

const (
	MYSQL DriverType = iota
    MARIADB
	ORACLE
	SQLITE
	POSTPRESQL

)

func SetDriver(d DriverType) {
	if d == MYSQL || d == MARIADB{
		driverChecker = MySQLChecker{}
	}else if d == ORACLE{
		driverChecker = OracleChecker{}
	}	else if d == SQLITE{
		driverChecker = SQLiteChecker{}
	}	else if d == POSTPRESQL{
		driverChecker = PostgreSQLChecker{}
	}else {
	}
}

var driverChecker IChecker

type Condition struct {
	Left     string
	Operator string
	Right    interface{}
}

func (c Condition) String() string {
	left, ok := c.leftCheck()
	if ! ok {
		return ""
	}
	right, ok := c.rightCheck()
	if ! ok {
		return ""
	}
	operator, ok := c.operatorCheck()
	if ! ok {
		return ""
	}
	return left + operator + right
}

// 检查字段名字,
func (c Condition) leftCheck() (string, bool) {
	if c.Left == "" {
		return "", false
	}
	return strings.Trim(c.Left, " "), true
}

func (c Condition) operatorCheck() (string, bool) {
	if c.Left == "" {
		return "", false
	}
	return strings.Trim(c.Operator, " "), true
}

func (c Condition) rightCheck() (string, bool) {

	switch c.Right.(type) {
	case float32, float64, int, int8, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", c.Right), true
	case string:
		r := c.Right.(string)
		if r != "" {
			return "'" + r + "'", true
		}
	}
	return "", false

}

type ConditionJoiner struct {
	Left     Condition
	Operator string
	Right    Condition
}

func init() {

}