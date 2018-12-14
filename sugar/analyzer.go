/*
# __author__ = "Mr.chai"
# Date: 2018/12/14
*/
package sugar

var confTypeList = []string{
	"yml",
	"json",
	"xml",
}

type analyzer interface {
}

type jsonAnalyzer struct {
	data []byte
}

func (ana jsonAnalyzer) dumps() {

}
func (ana jsonAnalyzer) dump() {

}

func (ana jsonAnalyzer) load() {

}

func (ana jsonAnalyzer) loads() {

}

type yamlAnalyzer struct {
	data []byte

}

func (ana yamlAnalyzer) dumps() {

}
func (ana yamlAnalyzer) dump() {

}

func (ana yamlAnalyzer) load() {

}

func (ana yamlAnalyzer) loads() {

}

type xmlAnalyzer struct {
	data []byte
}

func (ana xmlAnalyzer) dumps() {

}
func (ana xmlAnalyzer) dump() {

}

func (ana xmlAnalyzer) load() {

}

func (ana xmlAnalyzer) loads() {

}
