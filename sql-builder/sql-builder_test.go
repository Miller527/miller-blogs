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
	bd, err := Builder(MySQL)
	if err != nil {
		panic(err)
	}

	sel := bd.Select()
	xx, e := sel.Column(
		&MySQLColumnStatic{"name"},
		&MySQLColumnFunc{"MAX", "age"},
		&MySQLColumnAlias{"age", "ages"},
		&MySQLColumnAlias{&MySQLColumnFunc{"MIN", "ss"}, "SS"},
	//todo MySQLTableAlias自动生成别名
	).Table(&MySQLTableStatic{"t1"}, &MySQLTableAlias{"t2", "T222222"}).
		Where(
			sel.OR(
				sel.EQ("a", "b", "t1"),
				sel.LE("c", 1, ""),
				sel.AND(
					sel.EQ("d", "e", "t2"),
					sel.LE("f", 2, ""),

				),
			),

			sel.AND(
				sel.EQ("A", "B", ""),
				sel.LE("C", 11, ""),

			),
			sel.IsNull("nnn",""),
			sel.IsNotNull("n1n",""),

			sel.NOT(sel.IN("qqq","",2,3)),
			sel.NotIN("ddd","","a","b"),

		).Build()
	fmt.Println("---", xx, e)

}
