/*
# __author__ = "Mr.chai"
# Date: 2019-03-01
*/
package sql_generater

// select后的字段
type IColumn interface {
	Build () (result string, err error)
	verify() (err error)
}


type IBuilder interface {
	Select ()
	Update ()
}


type ISelect interface {

}



type IUpdate interface {

}


type IInsert interface {

}


type IDelete interface {

}
