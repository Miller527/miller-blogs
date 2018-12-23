/*
# __author__ = "Mr.chai"
# Date: 2018/12/9
*/
package sugar

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"miller-blogs/sugar/utils"
)

type Session interface {
	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
}

type CacheSession struct {
}

func (cs *CacheSession) Read(p []byte) (n int, err error) {
	return 1, nil
}

func (cs *CacheSession) Write(p []byte) (n int, err error) {
	return 1, nil
}

type RedisSession struct {
}

func (rs *RedisSession) Read(p []byte) (n int, err error) {
	return 1, nil
}

func (rs *RedisSession) Write(p []byte) (n int, err error) {
	return 1, nil
}

type FileSession struct {
}

func (fs *FileSession) Read(p []byte) (n int, err error) {
	return 1, nil
}

func (fs *FileSession) Write(p []byte) (n int, err error) {
	return 1, nil
}

// session config
type SessionConfig struct {
	Type string

	Size     int
	Protocol string
	Address  string
	Password string

	keys []string
}

var sessionStore sessions.Store

func InitSession(sessConf SessionConfig) {
	bytesKeys := [][]byte{}
	for _, k := range sessConf.keys {
		bytesKeys = append(bytesKeys, []byte(k))
	}
	store, err := redis.NewStore(sessConf.Size, sessConf.Protocol,
		sessConf.Address, sessConf.Password, bytesKeys...)
	utils.PanicCheck(err)

	sessionStore = store

}

