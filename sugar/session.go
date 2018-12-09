/*
# __author__ = "Mr.chai"
# Date: 2018/12/9
*/
package sugar

type Session interface {
	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
}


type CacheSession struct {

}

func (cs *CacheSession) Read(p []byte)(n int, err error){
	return 1,nil
}


func (cs *CacheSession) Write(p []byte)(n int, err error){
	return 1,nil
}


type RedisSession struct {

}

func (rs *RedisSession) Read(p []byte)(n int, err error){
	return 1,nil
}

func (rs *RedisSession) Write(p []byte)(n int, err error){
	return 1,nil
}

type FileSession struct {

}

func (fs *FileSession) Read(p []byte)(n int, err error){
	return 1,nil
}

func (fs *FileSession) Write(p []byte)(n int, err error){
	return 1,nil
}


func InitSession(session Session)  {

}