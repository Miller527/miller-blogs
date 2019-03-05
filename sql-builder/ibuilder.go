/*
# __author__ = "Miller"
# Date: 2019-03-01
*/
package sqlbuilder

// select后的字段和表名
type IBlock interface {
	Build() (result string, err error)
	verify() (err error)
}

type ITable interface {
	GetName() (result string)
	GetAlias() (result string)

}

// 表名字
type ITableBlock interface {
	IBlock
	ITable
}

type IBuilder interface {
	Select() ISelectBuilder
	BuilderJson() (string, error)
}

type ICondition interface {
	// 包含
	IN(column, tb string, values ...interface{}) IBlock
	NotIN(column, tb string, values ...interface{}) IBlock
	// 是否为空
	IsNull(column, tb string) IBlock
	IsNotNull(column, tb string) IBlock
	// less than 小于
	LT(column string, value interface{}, tb string) IBlock
	// less than or equal to 小于等于
	LE(column string, value interface{}, tb string) IBlock
	// equal to 等于
	EQ(column string, value interface{}, tb string) IBlock
	// not equal to 不等于
	NE(column string, value interface{}, tb string) IBlock
	// greater than or equal to 大于等于
	GE(column string, value interface{}, tb string) IBlock
	// greater than 大于
	GT(column string, value interface{}, tb string) IBlock
	//like
	LIKE(column, value, tb string) IBlock
	// 与
	AND(cols ...IBlock) IBlock
	// 或
	OR(cols ...IBlock) IBlock
	// 非
	NOT(cols IBlock) IBlock
}

type ISelectBuilder interface {
	Table(names ...IBlock) ISelectBuilder
	Column(cols ...IBlock) ISelectBuilder
	Where(cons ...IBlock) ISelectBuilder
	Build() (string, error)
	ICondition
}

type IUpdateBuilder interface {
	Table() IUpdateBuilder

	Where() IUpdateBuilder
	Build() (string, error)
	ICondition
}

type IInsertBuilder interface {
	Table() IUpdateBuilder
	Column() IUpdateBuilder
	Value() IUpdateBuilder
	Build() (string, error)
}

type IDeleteBuilder interface {
	Table() IDeleteBuilder

	Where() IDeleteBuilder
	ICondition

	Build() (string, error)
}
