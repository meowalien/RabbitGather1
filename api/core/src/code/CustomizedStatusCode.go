package code

import "github.com/gin-gonic/gin"

type CustomizedStatusCode struct {
	MyStatusCode int
	HttpCode int
	Message string
	ResponseStruct func(d StatusCode, c *gin.Context, data ...interface{})
}

func (s CustomizedStatusCode) Error() string {
	return s.String()
}

func (s CustomizedStatusCode) String() string {
	return s.Message
}

func (s CustomizedStatusCode) ResponseAsStruct(c *gin.Context, resp func(d StatusCode, c *gin.Context, data ...interface{}), data ...interface{}) {
	resp(s, c, data...)
}
func (s CustomizedStatusCode) Response(c *gin.Context, data ...interface{}) {
	s.ResponseAsStruct(c,s.ResponseStruct,  data...)
}
func (s CustomizedStatusCode) HTTPCode() int {
	return s.HttpCode
}

func (s CustomizedStatusCode) MyCode() int {
	return s.MyStatusCode
}


