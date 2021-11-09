package member

import (
	"context"
	"core/src/lib/decode_encode"
	"core/src/lib/errs"
	"core/src/module/files"
	"core/src/module/permission"
	"core/src/module/response"
	"core/src/module/token"
	"core/src/module/user"
	"core/src/module/vc_code"
	"fmt"
	"github.com/gin-gonic/gin"
)

type HTTP struct {
}

func (m *HTTP) Mount(ctx context.Context, engine *gin.Engine) error {
	router := engine.Group("/member")
	// 新建用戶
	router.POST("/signup", m.signup)
	// 豋入
	router.POST("/login", m.login)
	// 發送驗證碼
	router.POST("/send_vc", m.sendVerificationCode)

	return nil
}

func (h *HTTP) login(c *gin.Context) {
	type Request struct {
		Username string `form:"account" json:"account" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}
	var req Request
	type Response struct {
		UUID         string `form:"uuid" json:"uuid" binding:"required"`
		Token        string `form:"token" json:"token" binding:"required"`
		RefreshToken string `json:"refresh_token"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.StatusMissingParameters.JsonResponse(c, err)
		return
	}

	us, exist, err := user.GetUserByAccount(req.Username)
	if err != nil {
		response.ServerError.JsonResponse(c, err)
		return
	}
	if !exist {
		response.StatusAccountNotFound.JsonResponse(c)
		return
	}

	if yes, e := us.Frozen(); e != nil {
		response.ServerError.JsonResponse(c, e)
	} else if yes {
		response.StatusAccountFrozen.JsonResponse(c)
		return
	}

	passOK, err := us.CheckPassword(req.Password)
	if err != nil {
		response.ServerError.JsonResponse(c, err)
		return
	}
	if !passOK {
		response.StatusPasswordWrong.JsonResponse(c)
		return
	}

	newAPIToken, refreshToken, err := h.reCreateToken(us)
	if err != nil {
		response.ServerError.JsonResponse(c)
		return
	}

	response.OK.JsonResponse(c, Response{
		UUID:         us.UUID(),
		Token:        newAPIToken,
		RefreshToken: refreshToken,
	})
}

type File struct {
	FileBase64    string `form:"file_base64" json:"file_base64"`
	ExtensionName string `form:"extension_name" json:"extension_name"`
}

func (h HTTP) signup(c *gin.Context) {
	type Request struct {
		Account  string `form:"account" json:"account"`
		Password string `form:"password" json:"password"  binding:"required"`
		NickName string `form:"nick_name" json:"nick_name"`
		Email    string `form:"email" json:"email" binding:"required"`
		VCCode   string `form:"vc_code" json:"vc_code"  binding:"required"`
		VCType   string `form:"vc_type" json:"vc_type"  binding:"required"`
		Type     string `form:"type" json:"type" binding:"required"`
		Photo    File   `form:"photo" json:"photo"`
		Phone    string `form:"phone" json:"phone"`
		Gender   int8   `form:"gender" json:"gender"`
	}
	var req Request
	type Response struct {
		UUID         string `form:"uuid" json:"uuid" binding:"required"`
		Token        string `form:"token" json:"token" binding:"required"`
		RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.StatusMissingParameters.JsonResponse(c, err)
		return
	}

	_, exist, err := user.GetUserByAccount(req.Account)
	if err != nil {
		response.ServerError.JsonResponse(c,err)
		return
	}
	if exist{
		response.StatusAccountAlreadyExist.JsonResponse(c,"account duplicate")
		return
	}
	_, exist, err = user.GetUserByEmail(req.Email)
	if err != nil {
		response.ServerError.JsonResponse(c,err)
		return
	}
	if exist{
		response.StatusAccountAlreadyExist.JsonResponse(c,"email duplicate")
		return
	}


	var vcKey string

	switch req.VCType {
	case vc_code.EmailKey:
		vcKey = req.Email
	default:
		response.StatusWrongInput.JsonResponse(c, fmt.Sprintf("not supported vc_type: %s", req.VCType))
	}

	ok, err := vc_code.CheckVCCode(req.VCType, vcKey, req.VCCode)
	if err != nil {
		response.ServerError.JsonResponse(c, err)
		return
	}
	if !ok {
		response.StatusVerificationCodeWrong.JsonResponse(c)
		return
	}

	switch req.Type {
	case string(user.PlainUserType):

	default:
		response.StatusWrongInput.JsonResponse(c, "not supported user type")
		return
	}
	fileURL := ""
	if req.Photo.FileBase64 != "" {
		var bytes []byte
		bytes, err = decode_encode.DecodeBase64ToBytes(req.Photo.FileBase64)
		if err != nil {
			response.ServerError.JsonResponse(c, err)
			return
		}

		fileURL, err = files.TakeUploadFileURL(files.File{
			Bin:           bytes,
			ExtensionName: req.Photo.ExtensionName,
		})
		if err != nil {
			response.ServerError.JsonResponse(c, err)
			return
		}
	}

	_, err = user.CreateNewUser(user.CreateNewUserRequest{
		Type:     user.UserType(req.Type),
		Account:  req.Account,
		Rule:     []string{permission.LOGIN},
		PhotoURL: fileURL,
		Phone:    req.Phone,
		Email:    req.Email,
		Gender:   req.Gender,
		NickName: req.NickName,
		Password: req.Password,
	})
	if err != nil {
		response.ServerError.JsonResponse(c, err)
		return
	}

	//newAPIToken, refreshToken, err := h.reCreateToken(us)
	//if err != nil {
	//	response.ServerError.JsonResponse(c)
	//	return
	//}
	response.StatusNoContent.JsonResponse(c)
}

// 刷新重新頒發token
func (h HTTP) reCreateToken(us user.User) (tk string, reTk string, err error) {

	// 登出現有連線
	err = us.Logout()
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	err = us.CleanUserOwnToken()
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	refreshClaim := token.StandardRefreshClaims()

	// 發新的refresh_token
	reTk, err = token.NewTokenWithClaims(refreshClaim, us.UserClaim())
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	apiClaim := token.StandardAPIClaims()

	// 發新的token
	tk, err = token.NewTokenWithClaims(apiClaim, us.UserClaim())
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	// 將API token 與user綁定
	theID, _ := token.GetTokenIDInClaim(apiClaim)
	err = us.SaveUserOwnToken(theID)
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	return
}

func (h HTTP) sendVerificationCode(c *gin.Context) {
	type Request struct {
		VCType string `form:"vc_type" json:"vc_type"  binding:"required"`
		Email  string `from:"email"`
		Phone  string `from:"phone"`
	}
	var req Request

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.StatusMissingParameters.JsonResponse(c, err)
		return
	}

	var recipient string
	var vc string
	switch req.VCType {
	case vc_code.EmailKey:
		recipient = req.Email
		if req.Email == "" {
			response.StatusMissingParameters.JsonResponse(c, "email is empty")
			return
		}
		vc, err = vc_code.SendEmailVCCode(req.Email)
		if err != nil {
			response.ServerError.JsonResponse(c, err)
			return
		}
		fmt.Println("vc: ", vc)
	default:
		response.StatusWrongInput.JsonResponse(c, fmt.Sprintf("not supported vc_type:%s", req.VCType))
		return
	}
	err = vc_code.SaveVCCode(req.VCType, recipient, vc)
	if err != nil {
		response.ServerError.JsonResponse(c, err)
		return
	}

	response.StatusNoContent.JsonResponse(c)
}
