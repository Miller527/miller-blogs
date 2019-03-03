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

	sel := bd.Select()
	xx, e := sel.Column(
		&MySQLColumnStatic{"name"},
		&MySQLColumnFunc{"MAX","age"},
		&MySQLColumnAlias{"age","ages"},
		&MySQLColumnAlias{&MySQLColumnFunc{"MIN","ss"},"SS"},
	).Table(&MySQLTableStatic{"t1"}, &MySQLTableAlias{"t2","T2"}).
		Where(
			sel.OR(
				sel.EQ("a", "b",""),
				sel.LE("c", 1,""),
				sel.AND(
					sel.EQ("d", "e",""),
					sel.LE("f", 2,""),

				),
				),

		sel.AND(
				sel.EQ("A", "B",""),
				sel.LE("C", 11, ""),

			),
		).Build()
	fmt.Println("---",xx, e)
}
