/*
# __author__ = "Miller"
# Date: 2019-03-02
*/
package sqlbuilder

import (
	"fmt"
	"testing"
)

func TestMySQLBuilder(t *testing.T) {
	bd,err := Builder(MySQL)
	if err != nil{
		panic(err)
	}

	xx, e := bd.Select().Column(
		&MySQLColumnStatic{"name"},
		&MySQLColumnFunc{"MAX","age"},
		&MySQLColumnAlias{"age","ages"},
		&MySQLColumnAlias{&MySQLColumnFunc{"MIN","ss"},"SS"},
	).Table(&MySQLTableStatic{"t1"}, &MySQLTableAlias{"t2","T2"}).Build()
	fmt.Println("---",xx, e)
}
