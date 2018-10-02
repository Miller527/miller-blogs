//
// __author__ = "Miller"
// Date: 2018/9/30
//
package public

import (
	"errors"
	"fmt"
	"os"
	"path"
)

// 验证目录是否存在
func DirVerify(dirName string) (string, error) {
	currentDir, _ := os.Getwd()
	tmpPath := path.Join(currentDir, dirName)
	tmpFileInfo, err := os.Stat(tmpPath)
	if err == nil {
		if ! tmpFileInfo.IsDir() {
			err = errors.New("File of the same name")
		}
	} else {
		err = os.Mkdir(tmpPath, 0755)

	}
	return tmpPath, nil
}



// 验证元素是否存在于列表, 并返回索引
func ElementInList(li []string, ele string) (int, bool) {
	for index, val := range li{
		fmt.Println(val, ele)
		if val == ele {
			return index, true
		}
	}
	return -1, false
}

