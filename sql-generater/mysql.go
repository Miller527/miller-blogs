/*
# __author__ = "Mr.chai"
# Date: 2019-03-01
*/
package sql_generater


type ColumnStatic struct {
	Column string
}

func (cs ColumnStatic) verify()error{

	return nil
}


func (cs ColumnStatic) Build() (string, error){

	return "", nil
}

type ColumnAlias struct {
	Left string
	Right string
	Column string

}
func (ca *ColumnAlias) verify()error{

	return nil
}


func (ca *ColumnAlias) Build() (string, error){

	return "", nil
}

type ColumnFunc struct {

}
func (cf *ColumnFunc) verify()error{

	return nil
}


func (cf *ColumnFunc) Build() (string, error){

	return "", nil
}



type Column map[string]interface{}

func SelectColumn(args ...interface{}){
	var columnList []IColumn

	for i:=0; i<= len(args);i++{
		col := args[i]
		switch col.(type) {

		case string:
			columnList = append(columnList, ColumnStatic{col.(string)})
		}
	}
}


type MySqlBuilder struct {

}

