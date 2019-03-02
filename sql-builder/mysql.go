/*
# __author__ = "Miller"
# Date: 2019-03-01
*/
package sqlbuilder

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type MySQLColumnStatic struct {
	Column string
}

func (cs *MySQLColumnStatic) verify() error {
	if cs.Column != "" {
		return nil
	}

	return errors.New("MySQLColumnStatic params 'Column' error")
}

func (cs *MySQLColumnStatic) Build() (string, error) {
	err := cs.verify()
	if err != nil {
		return "", err
	}
	return "`" + cs.Column + "`", nil
}

type MySQLColumnAlias struct {
	Column  interface{}
	Alias string
}

func (ca *MySQLColumnAlias) verify() error {
	if ca.Alias == "" {
		return errors.New("MySQLColumnAlias params 'Alias' is None")
	}

	switch ca.Column.(type) {
	case string:
		if ca.Column.(string) == "" {
			return errors.New("MySQLColumnAlias params 'Column' is None")
		}
		ca.Column = "`"+ca.Column.(string)+"`"
		return nil
	case *MySQLColumnFunc:
		str, err := ca.Column.(*MySQLColumnFunc).Build()
		if err != nil {
			return err
		}
		ca.Column = str
		return nil
	default:
		return errors.New("MySQLColumnAlias params 'Column' type error")
	}

}

func (ca *MySQLColumnAlias) Build() (string, error) {
	err := ca.verify()
	if err != nil {
		return "", err
	}
	return  fmt.Sprintf("%s AS `%s`", ca.Column.(string), ca.Alias),	nil
}

type MySQLColumnFunc struct {
	Column  string
	FuncName string
}


func (cf *MySQLColumnFunc) verify() error {
	if cf.Column == "" {
		return errors.New("MySQLColumnFunc params 'Column' is None")
	}
	if cf.FuncName == "" {
		return errors.New("MySQLColumnFunc params 'FuncName' is None")
	}
	return nil
}

func (cf *MySQLColumnFunc) Build() (string, error) {
	err := cf.verify()
	if err != nil {
		return "", err
	}

	return cf.Column + "(`" + cf.FuncName + "`)", nil
}

type MySQLTableStatic struct {
	Name string
}

func (ts *MySQLTableStatic) verify() error {

	if ts.Name == "" {
		return errors.New("MySQLTableStatic params 'Name' is None")
	}
	return nil
}

func (ts *MySQLTableStatic) Build() (string, error) {
	err := ts.verify()
	if err != nil {
		return "", err
	}

	return "`" + ts.Name + "`", nil
}

type MySQLTableAlias struct {
	Name  string
	Alias string
}

func (ta *MySQLTableAlias) verify() error {
	if ta.Name == "" {
		return errors.New("MySQLTableAlias params 'Name' is None")
	}
	if ta.Alias == "" {
		return errors.New("MySQLTableAlias params 'Alias' is None")
	}
	return nil
}

func (ta *MySQLTableAlias) Build() (string, error) {
	err := ta.verify()
	if err != nil {
		return "", err
	}
	return  fmt.Sprintf("`%s` AS `%s`", ta.Name, ta.Alias),	nil
}

func SelectColumn(args ...interface{}) {
	var columnList []IColumn

	for i := 0; i <= len(args); i++ {
		col := args[i]
		switch col.(type) {

		case string:
			columnList = append(columnList, &MySQLColumnStatic{col.(string)})
		}
	}
}

type MySQLSelect struct {
	column    []IColumn
	columnStr string
	tbName []ITableName
	tbNameStr string
	Err       error
}

// 将column生成字符串
func (sb *MySQLSelect) columnBuild() (string, error) {
	//var colList []string
	//for i := 0; i < len(sb.column); i++ {
	//	res, err := sb.column[i].Build()
	//	if err != nil {
	//		return "", err
	//	}
	//	colList = append(colList, res)
	//
	//}
	//result := strings.Join(colList, ", ")
	//return result, nil
	return sb.publicBuild("column")
}

func (sb *MySQLSelect) publicBuild(t string) (string, error) {
	length :=0
	if t == "table"{
		length = len(sb.tbName)
	}else {
		length = len(sb.column)
	}
	var resList []string
	var res string
	var err error

	for i := 0; i < length; i++ {
		if t =="table"{
			res, err = sb.tbName[i].Build()
		}else{
			res, err = sb.column[i].Build()

		}
		if err != nil {
			return "", err
		}
		resList = append(resList, res)

	}
	return strings.Join(resList, ", "), err

}


// 进行格式化和赋值
func (sb *MySQLSelect) Column(cols ...IColumn) ISelectBuilder {
	if sb.Err != nil {
		return sb
	}
	if len(cols) == 0 {
		sb.Err = errors.New("MySQLSelect params 'column' is None")
		return sb
	}
	sb.column = cols
	sb.columnStr, sb.Err = sb.columnBuild()

	return sb
}
func (sb *MySQLSelect) tableBuild() (string, error) {
	//var colList []string
	//for i := 0; i < len(sb.column); i++ {
	//	res, err := sb.column[i].Build()
	//	if err != nil {
	//		return "", err
	//	}
	//	colList = append(colList, res)
	//
	//}
	//result := strings.Join(colList, ", ")
	//return result, nil

	return sb.publicBuild("table")

}

// 表名
func (sb *MySQLSelect) Table(names ...ITableName) ISelectBuilder {
	if sb.Err != nil {
		return sb
	}
	if len(names) == 0 {
		sb.Err = errors.New("MySQLSelect tbName 'column' is None")
		return sb
	}
	sb.tbName = names
	sb.tbNameStr, sb.Err = sb.tableBuild()

	return sb
}


func (sb *MySQLSelect) Where() ISelectBuilder {
	if sb.Err != nil {
		return sb
	}
	return sb
}

// 校验各个流程的值是否合法
func (sb *MySQLSelect) checkBuild() error {
	if sb.Err != nil {
		return sb.Err
	}
	if sb.columnStr == ""{
		return errors.New("Column is None")
	}
	if sb.tbNameStr == ""{
		return errors.New("Table name is None")
	}
	return nil
}

// 最终的构建
func (sb *MySQLSelect) Build() (string, error) {
	err := sb.checkBuild()
	if err != nil{
		return "",err
	}
	cmd := "SELECT " + sb.columnStr +" FROM " + sb.tbNameStr
	return cmd, nil
}

func Select(column IColumn) ISelectBuilder {

	return &MySQLSelect{}
}

type MySqlBuilder struct {
	// todo 用户传入json配置自动生成
	Params interface{}
}

func (mb MySqlBuilder) Select() ISelectBuilder {
	return &MySQLSelect{}
}

func (mb MySqlBuilder) BuilderJson() (string, error) {

	return "", nil

}
