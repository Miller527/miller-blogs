/*
# __author__ = "Miller"
# Date: 2019-03-01
*/
package sqlbuilder

// select后的字段
type IColumn interface {
	Build () (result string, err error)
	verify() (err error)
}
// 表名字
type ITableName interface {
	IColumn
}





type IBuilder interface {
	Select () ISelectBuilder
	BuilderJson () (string, error)
}


type ISelectBuilder interface {

Table(names ...ITableName)ISelectBuilder
Column (cols ...IColumn)ISelectBuilder
Where() ISelectBuilder
Build () (string, error)
}



type IUpdateBuilder interface {
	Table() IUpdateBuilder

	Where() IUpdateBuilder
	Build () (string, error)
}


type IInsertBuilder interface {
	Table() IUpdateBuilder
	Column()IUpdateBuilder
	Value()IUpdateBuilder

	Build () (string, error)
}


type IDeleteBuilder interface {
	Table() IDeleteBuilder

	Where() IDeleteBuilder
	Build () (string, error)
}
