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
	Column interface{}
	Alias  string
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
		ca.Column = "`" + ca.Column.(string) + "`"
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
	return fmt.Sprintf("%s AS `%s`", ca.Column.(string), ca.Alias), nil
}

type MySQLColumnFunc struct {
	Column   string
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
	//todo 支持不同库
	//Db string
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
	return fmt.Sprintf("`%s` AS `%s`", ta.Name, ta.Alias), nil
}

type MySQLLogicCondition struct {
	Name  string
	Alias string
}

func (mlc *MySQLLogicCondition) verify() error {
	if mlc.Name == "" {
		return errors.New("MySQLTableAlias params 'Name' is None")
	}
	if mlc.Alias == "" {
		return errors.New("MySQLTableAlias params 'Alias' is None")
	}
	return nil
}

func (mlc *MySQLLogicCondition) Build() (string, error) {
	err := mlc.verify()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("`%s` AS `%s`", mlc.Name, mlc.Alias), nil
}

type mySQLCondition struct {
	Column string
	Op     string
	Value  interface{}
	DBName string
	Format string
}


func (mc *mySQLCondition) verify() error {
	if mc.Column == "" {
		return errors.New("mySQLCondition params 'Column' is None")
	}
	if mc.Op == "" {
		return errors.New("mySQLCondition params 'Op' is None")
	}
	fmt.Println("xxxx",mc.Value, mc.Value=="", mc.Value==nil)

	switch mc.Value.(type) {
	case string:
		fmt.Println("xxxxxxxx")

		if mc.Value.(string) == ""{
			return errors.New("mySQLCondition params 'Value' is None")
		}
	default:
		if mc.Value == nil{
			return errors.New("mySQLCondition params 'Value' is None")

		}
	}
	if mc.Format == "" {
		return errors.New("mySQLCondition params 'Format' is None")
	}
	return nil
}

func (mc *mySQLCondition) Build() (string, error) {
	err := mc.verify()
	if err != nil {
		return "", err
	}
	if mc.DBName != "" {
		mc.Column = mc.DBName + "." + mc.Column
	}

	return fmt.Sprintf(mc.Format, mc.Column, mc.Op, mc.Value), nil
}

func SelectColumn(args ...interface{}) {
	var columnList []IBlock

	for i := 0; i <= len(args); i++ {
		col := args[i]
		switch col.(type) {

		case string:
			columnList = append(columnList, &MySQLColumnStatic{col.(string)})
		}
	}
}

type MySQLSelect struct {
	column    []IBlock
	columnStr string
	tbName    []IBlock
	tbNameStr string
	where     []IBlock
	whereStr  string
	Err       error
}

// 包含
func (sb *MySQLSelect) IN(column, tb string, values ...interface{}) IBlock {
	return nil

}
// 不包含
func (sb *MySQLSelect) NotIN(column, tb string, values ...interface{}) IBlock {
	return nil

}

// 是否为空
func (sb *MySQLSelect) IsNull(column, tb string) IBlock {
	return nil
}
func (sb *MySQLSelect) IsNotNull(column, tb string) IBlock {
	return nil

}

func (sb *MySQLSelect) operator(col, op string, val interface{}, tb, format, strSep string ) IBlock {
	switch val.(type) {
	case string:
		if val.(string) != ""{
			val = strSep + val.(string) + strSep

		}
	}
	return &mySQLCondition{col, op, val, tb, format}
}

// less than 小于
func (sb *MySQLSelect) LT(column string, value interface{}, tb string) IBlock {
	return sb.operator(column, "<", value, tb, "`%s`%s%v","'")

}

// less than or equal to 小于等于
func (sb *MySQLSelect) LE(column string, value interface{}, tb string) IBlock {
	return sb.operator(column, "<=", value, tb, "`%s`%s%v","'")

}

// equal to 等于
func (sb *MySQLSelect) EQ(column string, value interface{}, tb string) IBlock {
	return sb.operator(column, "=", value, tb, "`%s`%s%v","'")

}

// not equal to 不等于
func (sb *MySQLSelect) NE(column string, value interface{}, tb string) IBlock {
	return sb.operator(column, "!=", value, tb, "`%s`%s%v","'")

}

// greater than or equal to 大于等于
func (sb *MySQLSelect) GE(column string, value interface{}, tb string) IBlock {
	return sb.operator(column, ">=", value, tb, "`%s`%s%v","'")

}

// greater than 大于
func (sb *MySQLSelect) GT(column string, value interface{}, tb string) IBlock {
	return sb.operator(column, ">", value, tb, "`%s`%s%v","'")

}
// like
func (sb *MySQLSelect) LIKE(column, value, tb string) IBlock {
	return sb.operator(column, " LIKE ", value, tb, "`%s`%s%v","'")

}
// 与
func (sb *MySQLSelect) AND(cols ...IBlock) IBlock {

	if len(cols) == 0{
		return nil
	}else	if len(cols) == 1{
		return cols[0]
	}
	s1,s2 := sb.logic(cols...)
	if s1 =="" || s2 == ""{
		return nil

	}
	fmt.Println("aaaaaa", sb.Err)

	return sb.operator(s1, "AND", s2, "", "( %s %s %s )","")

}

func (sb *MySQLSelect) logic(cols ...IBlock) (string, string) {
	s1, e1 := cols[0].Build()
	if e1 != nil{
		sb.Err = e1
		return "",""
	}

	s2,e2 := sb.publicBuild(cols[1:], nil, " OR ")

	if e2 != nil{
		sb.Err = e2
		return "", ""
	}

	return s1, s2
}



// 或
func (sb *MySQLSelect) OR(cols ...IBlock) IBlock {
	if len(cols) == 0{
		return nil
	}else	if len(cols) == 1{
		return cols[0]
	}
	s1,s2 := sb.logic(cols...)
	if s1 =="" || s2 == ""{
		return nil

	}
	fmt.Println(s1)
	fmt.Println(s2)
	return sb.operator(s1, "OR", s2, "", "( %s %s %s )","")


}

// 非
func (sb *MySQLSelect) NOT(cols IBlock) IBlock {
	return nil

}



// 将column生成字符串
func (sb *MySQLSelect) columnBuild(err error) (string, error) {
	return sb.publicBuild(sb.column, err,", ")
}

// 各模块的公共构建方式
func (sb *MySQLSelect) publicBuild(blocks []IBlock, err error, sep string) (string, error) {

	if len(blocks) == 0 {
		return "", err
	}

	var resList []string
	for i := 0; i < len(blocks); i++ {
		if sb.Err != nil{
			return "", sb.Err
		}
		res, err := blocks[i].Build()
		if err != nil {
			return "", err
		}
		resList = append(resList, res)
	}
	return strings.Join(resList, sep), nil

}

// 进行格式化和赋值
func (sb *MySQLSelect) Column(cols ...IBlock) ISelectBuilder {
	if sb.Err != nil {
		return sb
	}
	sb.column = cols
	sb.columnStr, sb.Err = sb.columnBuild(columnNoneErr)

	return sb
}

// 表名构建
func (sb *MySQLSelect) tableBuild(err error) (string, error) {
	return sb.publicBuild(sb.tbName, err, ", ")
}

// 表名
func (sb *MySQLSelect) Table(names ...IBlock) ISelectBuilder {
	if sb.Err != nil {
		return sb
	}
	sb.tbName = names
	sb.tbNameStr, sb.Err = sb.tableBuild(tbNameNoneErr)
	return sb
}

// 条件构建
func (sb *MySQLSelect) whereBuild(err error) (string, error) {
	x,c := sb.publicBuild(sb.where, err, " AND ")
	fmt.Println("123",x,c)
	return sb.publicBuild(sb.where, err, " AND ")
}

var whereNoneErr = errors.New("MySQLSelect params 'where' is None")
var tbNameNoneErr = errors.New("MySQLSelect params 'tbName' is None")
var columnNoneErr = errors.New("MySQLSelect params 'column' is None")

func (sb *MySQLSelect) Where(cons ...IBlock) ISelectBuilder {
	if sb.Err != nil {
		return sb
	}

	if len(cons) == 0 {
		return sb
	}
	sb.where = cons
	sb.whereStr, sb.Err = sb.whereBuild(whereNoneErr)
	return sb
}

// 校验各个流程的值是否合法
func (sb *MySQLSelect) checkBuild() error {
	if sb.Err != nil {
		return sb.Err
	}
	if sb.columnStr == "" {
		return errors.New("Column is None")
	}
	if sb.tbNameStr == "" {
		return errors.New("Table name is None")
	}
	return nil
}

// 最终的构建
func (sb *MySQLSelect) Build() (string, error) {
	fmt.Println("errrrrrr", sb.Err)
	err := sb.checkBuild()
	if err != nil {
		return "", err
	}
	cmd := "SELECT " + sb.columnStr + " FROM " + sb.tbNameStr
	if sb.whereStr != "" {
		cmd = cmd + " WHERE " + sb.whereStr
	}
	return cmd, nil
}

//func Select(column IBlock) ISelectBuilder {
//
//	return &MySQLSelect{}
//}

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
