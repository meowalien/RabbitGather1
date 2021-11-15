package code

import (
	"core/src/conf"
	"core/src/module/log"
	"fmt"
	"github.com/gin-gonic/gin"
)

type StatusCode interface {
	String() string
	Response(c *gin.Context,resp func(d StatusCode, c *gin.Context, data ...interface{}), data ...interface{})
	JsonResponse(c *gin.Context, data ...interface{})
	HTTPCode() int
	MyCode() int
	FlatJsonResponse(c *gin.Context, data ...interface{})
}
type httpStatusCode int

func (d httpStatusCode) Response(c *gin.Context,resp func(d StatusCode, c *gin.Context, data ...interface{}), data ...interface{}) {
	resp(d,c, data...)
}
func (d httpStatusCode) JsonResponse(c *gin.Context, data ...interface{}) {
	d.Response(c ,JsonResponseWrapper(StandardJsonResponse) , data... )
}
func (d httpStatusCode) FlatJsonResponse(c *gin.Context, data ...interface{}) {
	d.Response(c,MapJsonResponseWrapper(standardFlatJsonResponse),data...)
}
func (d httpStatusCode) String() string {
	code, exist := statusText[d]
	if exist {
		return code
	}
	return "undefined code"
}

func (d httpStatusCode) HTTPCode() int {
	return int(d)
}
func (d httpStatusCode) MyCode() int {
	return int(d)
}

type myStatusCode httpStatusCode

func (m myStatusCode) Response(c *gin.Context, resp func(d StatusCode, c *gin.Context, data ...interface{}), data ...interface{}) {
	resp(m,c, data...)
}

func (m myStatusCode) FlatJsonResponse(c *gin.Context, data ...interface{}) {
	m.Response(c,MapJsonResponseWrapper(standardFlatJsonResponse),data...)
}

func (m myStatusCode) String() string {
	code, exist := myStatusText[m]
	if exist {
		return code
	}
	return "undefined code"
}

func (m myStatusCode) JsonResponse(c *gin.Context, data ...interface{}) {
	m.Response(c ,JsonResponseWrapper(StandardJsonResponse) , data... )
}

func (m myStatusCode) HTTPCode() int {
	return 200
}

func (m myStatusCode) MyCode() int {
	return int(m)
}


func JsonResponseWrapper(resp  func (c *gin.Context, statusCode int, myCode int, msg string, data interface{})) func(d StatusCode, c *gin.Context, data ...interface{}){
	return func(d StatusCode, c *gin.Context, data ...interface{}) {
		if data == nil || len(data) == 0 {
			st := d.String()
			if st == "" {
				// 沒有body的回應（如204）
				c.AbortWithStatus(d.HTTPCode())
				return
			} else {
				resp(c, d.HTTPCode(), d.MyCode(), st, nil)
				return
			}
		}

		var msg = ""

		dataLen := len(data)
		var i interface{}

		if dataLen == 1 {
			switch t := data[0].(type) {
			case string:
				if t != "" {
					msg += t
				} else {
					msg = d.String()
				}
			case error:
				if conf.DEBUG_MOD {
					msg = t.Error()
				} else {
					msg = d.String()
				}
				log.Logger.Skip(2).Error(t.Error())
			default:
				// 非信息型態，作為資料處理
				msg = d.String()
				i = data[0]
			}
		} else if dataLen == 2 {
			switch t := data[0].(type) {
			case string:
				if t != "" {
					msg += t
				} else {
					msg = d.String()
				}
			case error:
				if conf.DEBUG_MOD {
					msg = t.Error()
				} else {
					msg = d.String()
				}
				log.Logger.Skip(2).Error(t.Error())
			case nil:
				msg = d.String()
				log.Logger.Skip(2).Error("got nil on first parameter ")
			default:
				panic("the first must be msg( string type)")
			}
			i = data[1]
		} else {
			panic("the data should not be over 2")
		}
		resp(c, d.HTTPCode(), d.MyCode(), msg, i)
	}
}


func StandardJsonResponse(c *gin.Context, statusCode int, myCode int, msg string, data interface{}) {
	if data == nil {
		c.AbortWithStatusJSON(statusCode, gin.H{
			"code": myCode,
			"msg":  msg,
		})
	} else {
		c.AbortWithStatusJSON(statusCode, gin.H{
			"code": myCode,
			"msg":  msg,
			"data": data,
		})
	}
}


func MapJsonResponseWrapper(resp  func (c *gin.Context, statusCode int, myCode int, msg string, data interface{})) func(d StatusCode, c *gin.Context, data ...interface{}) {
	return func(d StatusCode, c *gin.Context, data ...interface{}) {
		if data == nil || len(data) == 0 {
			st := d.String()
			if st == "" {
				// 沒有body的回應（如204）
				c.AbortWithStatus(d.HTTPCode())
				return
			} else {
				resp(c, d.HTTPCode(), d.MyCode(), st, nil)
				return
			}
		}

		var msg = ""

		dataLen := len(data)
		var i map[string]interface{}

		if dataLen == 1 {
			switch t := data[0].(type) {
			case string:
				if t != "" {
					msg += t
				} else {
					msg = d.String()
				}
			case error:
				if conf.DEBUG_MOD {
					msg = t.Error()
				} else {
					msg = d.String()
				}
				log.Logger.Skip(2).Error(t.Error())
			case map[string]interface{}:
				msg = d.String()
				i = t
			case nil:
				msg = d.String()
				log.Logger.Skip(2).Error("got nil on first parameter ")

			default:
				panic(fmt.Sprintf("not supported type:%v", t))
			}
		} else if dataLen == 2 {
			switch t := data[0].(type) {
			case string:
				if t != "" {
					msg += t
				} else {
					msg = d.String()
				}
			case error:
				if conf.DEBUG_MOD {
					msg = t.Error()
				} else {
					msg = d.String()
				}
				log.Logger.Skip(2).Error(t.Error())
			default:
				panic("the first must be msg( string type)")
			}

			mp, ok := data[1].(map[string]interface{})
			if !ok {
				panic("the second parameter should be map[string]interface{}")
			}
			i = mp
		} else {
			panic("the data should not be over 2")
		}
		resp(c, d.HTTPCode(), d.MyCode(), msg, i)
	}
}

func standardFlatJsonResponse(c *gin.Context, code int, myCode int, msg string, m interface{}) {
	mp := m.(map[string]interface{})
	if mp == nil {
		mp = map[string]interface{}{}
	}
	mp["code"] = myCode
	mp["msg"] = msg
	c.AbortWithStatusJSON(code, mp)
}






// RFC code
const (
	StatusContinue           httpStatusCode = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols httpStatusCode = 101 // RFC 7231, 6.2.2
	StatusProcessing         httpStatusCode = 102 // RFC 2518, 10.1
	StatusEarlyHints         httpStatusCode = 103 // RFC 8297

	OK                         httpStatusCode = 200 // RFC 7231, 6.3.1
	Created                    httpStatusCode = 201 // RFC 7231, 6.3.2
	Accepted                   httpStatusCode = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo httpStatusCode = 203 // RFC 7231, 6.3.4
	NoContent                  httpStatusCode = 204 // RFC 7231, 6.3.5
	StatusResetContent         httpStatusCode = 205 // RFC 7231, 6.3.6
	StatusPartialContent       httpStatusCode = 206 // RFC 7233, 4.1
	StatusMultiStatus          httpStatusCode = 207 // RFC 4918, 11.1
	StatusAlreadyReported      httpStatusCode = 208 // RFC 5842, 7.1
	StatusIMUsed               httpStatusCode = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   httpStatusCode = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  httpStatusCode = 301 // RFC 7231, 6.4.2
	StatusFound             httpStatusCode = 302 // RFC 7231, 6.4.3
	StatusSeeOther          httpStatusCode = 303 // RFC 7231, 6.4.4
	StatusNotModified       httpStatusCode = 304 // RFC 7232, 4.1
	StatusUseProxy          httpStatusCode = 305 // RFC 7231, 6.4.5
	StatusTemporaryRedirect httpStatusCode = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect httpStatusCode = 308 // RFC 7538, 3

	BadRequest httpStatusCode = 400 // RFC 7231, 6.5.1
	// 未豋入之類
	Unauthorized          httpStatusCode = 401 // RFC 7235, 3.1
	StatusPaymentRequired httpStatusCode = 402 // RFC 7231, 6.5.2
	// 已豋入但權限不族
	StatusForbidden                    httpStatusCode = 403 // RFC 7231, 6.5.3
	NotFound                           httpStatusCode = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             httpStatusCode = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                httpStatusCode = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            httpStatusCode = 407 // RFC 7235, 3.2
	StatusRequestTimeout               httpStatusCode = 408 // RFC 7231, 6.5.7
	StatusConflict                     httpStatusCode = 409 // RFC 7231, 6.5.8
	StatusGone                         httpStatusCode = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               httpStatusCode = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           httpStatusCode = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        httpStatusCode = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            httpStatusCode = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         httpStatusCode = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable httpStatusCode = 416 // RFC 7233, 4.4
	StatusExpectationFailed            httpStatusCode = 417 // RFC 7231, 6.5.14
	StatusTeapot                       httpStatusCode = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest           httpStatusCode = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity          httpStatusCode = 422 // RFC 4918, 11.2
	StatusLocked                       httpStatusCode = 423 // RFC 4918, 11.3
	StatusFailedDependency             httpStatusCode = 424 // RFC 4918, 11.4
	StatusTooEarly                     httpStatusCode = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              httpStatusCode = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         httpStatusCode = 428 // RFC 6585, 3
	StatusTooManyRequests              httpStatusCode = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  httpStatusCode = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   httpStatusCode = 451 // RFC 7725, 3

	ServerError                   httpStatusCode = 500 // RFC 7231, 6.6.1
	StatusNotImplemented          httpStatusCode = 501 // RFC 7231, 6.6.2
	StatusBadGateway              httpStatusCode = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable      httpStatusCode = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout          httpStatusCode = 504 // RFC 7231, 6.6.5
	HTTPVersionNotSupported       httpStatusCode = 505 // RFC 7231, 6.6.6
	VariantAlsoNegotiates         httpStatusCode = 506 // RFC 2295, 8.1
	InsufficientStorage           httpStatusCode = 507 // RFC 4918, 11.5
	LoopDetected                  httpStatusCode = 508 // RFC 5842, 7.2
	NotExtended                   httpStatusCode = 510 // RFC 2774, 7
	NetworkAuthenticationRequired httpStatusCode = 511 // RFC 6585, 6
)

var statusText = map[httpStatusCode]string{
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "Switching Protocols",
	StatusProcessing:         "Processing",
	StatusEarlyHints:         "Early Hints",

	OK:                         "OK",
	Created:                    "Created",
	Accepted:                   "Accepted",
	StatusNonAuthoritativeInfo: "Non-Authoritative Information",
	NoContent:                  "No Content",
	StatusResetContent:         "Reset Content",
	StatusPartialContent:       "Partial Content",
	StatusMultiStatus:          "Multi-Status",
	StatusAlreadyReported:      "Already Reported",
	StatusIMUsed:               "IM Used",

	StatusMultipleChoices:   "Multiple Choices",
	StatusMovedPermanently:  "Moved Permanently",
	StatusFound:             "Found",
	StatusSeeOther:          "See Other",
	StatusNotModified:       "Not Modified",
	StatusUseProxy:          "Use Proxy",
	StatusTemporaryRedirect: "Temporary Redirect",
	StatusPermanentRedirect: "Permanent Redirect",

	BadRequest:                         "Bad Request",
	Unauthorized:                       "Unauthorized",
	StatusPaymentRequired:              "Payment Required",
	StatusForbidden:                    "Forbidden",
	NotFound:                           "Not Found",
	StatusMethodNotAllowed:             "Method Not Allowed",
	StatusNotAcceptable:                "Not Acceptable",
	StatusProxyAuthRequired:            "Proxy Authentication Required",
	StatusRequestTimeout:               "Request Timeout",
	StatusConflict:                     "Conflict",
	StatusGone:                         "Gone",
	StatusLengthRequired:               "Length Required",
	StatusPreconditionFailed:           "Precondition Failed",
	StatusRequestEntityTooLarge:        "Request Entity Too Large",
	StatusRequestURITooLong:            "Request URI Too Long",
	StatusUnsupportedMediaType:         "Unsupported Media Type",
	StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	StatusExpectationFailed:            "Expectation Failed",
	StatusTeapot:                       "I'm a teapot",
	StatusMisdirectedRequest:           "Misdirected Request",
	StatusUnprocessableEntity:          "Unprocessable Entity",
	StatusLocked:                       "Locked",
	StatusFailedDependency:             "Failed Dependency",
	StatusTooEarly:                     "Too Early",
	StatusUpgradeRequired:              "Upgrade Required",
	StatusPreconditionRequired:         "Precondition Required",
	StatusTooManyRequests:              "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	ServerError:                   "Internal Server Error",
	StatusNotImplemented:          "Not Implemented",
	StatusBadGateway:              "Bad Gateway",
	StatusServiceUnavailable:      "Service Unavailable",
	StatusGatewayTimeout:          "Gateway Timeout",
	HTTPVersionNotSupported:       "HTTP Version Not Supported",
	VariantAlsoNegotiates:         "Variant Also Negotiates",
	InsufficientStorage:           "Insufficient Storage",
	LoopDetected:                  "Loop Detected",
	NotExtended:                   "Not Extended",
	NetworkAuthenticationRequired: "Network Authentication Required",
}

const StatusCodeMagnification = 1000

const (
	StatusAccountNotFound      = myStatusCode(NotFound*StatusCodeMagnification) + 1 + iota
	StatusAccountFrozen
	StatusAccountAlreadyExist
	StatusWalletNotFound
	NoCorrespondingOrderNumber
	StatusOrderNotExist
)

const (
	StatusMissingParameters     = myStatusCode(BadRequest*StatusCodeMagnification )+ 1 + iota
	MissingInputValue
	StatusVerificationCodeWrong
	StatusWalletBalanceOverflow
	StatusSignatureIllegal
	StatusOrderDuplicated
	OrderAlreadyCanceled
	StatusWrongInput
)

const (
	StatusPasswordWrong = myStatusCode(StatusForbidden*StatusCodeMagnification) + 1 + iota
)

var myStatusText = map[myStatusCode]string{
	StatusAccountNotFound:       "account not found",
	StatusAccountFrozen:         "account frozen",
	StatusMissingParameters:     "missing parameters",
	MissingInputValue:           "missing input value",
	StatusVerificationCodeWrong: "verification code wrong",
	StatusAccountAlreadyExist:   "account already exist",
	StatusPasswordWrong:         "password wrong",
	StatusWalletNotFound:        "wallet not found",
	StatusSignatureIllegal:      "signature illegal",
	StatusOrderDuplicated:       "order duplicated",
	StatusWalletBalanceOverflow: "status wallet balance overflow",
	NoCorrespondingOrderNumber:  "no corresponding order number",
	StatusOrderNotExist:         "order not exist",
	OrderAlreadyCanceled:        "order already canceled",
	StatusWrongInput:            "status wrong input",
}
