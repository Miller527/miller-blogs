/*
# __author__ = "Miller"
# Date: 2019-02-18
*/
package sqlbuilder

type IChecker interface {
	// 条件体校验
	//Condition(ConditionBody) (string, bool)
	// 条件体字符串校验
	ConditionStr(string) (string, bool)
	// 请求体校验
	Query()
	// 请求体字符串校验
	QueryStr()
}

type MySQLChecker struct {

}


type OracleChecker struct {

}


type SQLiteChecker struct {

}


type PostgreSQLChecker struct {

}




