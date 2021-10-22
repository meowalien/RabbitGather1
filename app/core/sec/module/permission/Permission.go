package permission

import (
	"core/sec/conf"
	token_module "core/sec/module/token"
	"core/sec/module/user"
	"fmt"
)

// return true if pass
// skip: path , method
func CheckPermission(stringToken string, path string, method string) (bool, error) {
	if stringToken == "" || path == "" || method == "" {
		return false, fmt.Errorf("input should not be empty string")
	}

	var mapClaim map[string]interface{}
	token, err := token_module.ParseTokenWithClaim(stringToken, conf.GlobalConfig.JWT.SignMethod, &mapClaim)
	if err != nil {
		return false, fmt.Errorf("error when ParseTokenWithClaim. %w", err)
	}
	if !token.Valid {
		return false, nil
	}

	// token有沒有被作廢
	ok, err := token_module.CheckTokenActiveWithClaim(mapClaim)
	if err != nil {
		return false, fmt.Errorf("error when CheckTokenActiveWithClaim. %w", err)
	}
	if !ok {
		return false, nil
	}

	userUUID, err := user.GetUserUUIDInClaim(mapClaim)
	if err != nil {
		return false, fmt.Errorf("error when GetUserUUIDInClaim. %w", err)
	}

	userInst, exist, err := user.GetUserByUUID(userUUID)
	if err != nil {
		return false, fmt.Errorf("error when GetUserByID. %w", err)
	}
	if !exist {
		return false, nil
	}

	if !userInst.IsActive() {
		return false, fmt.Errorf("the user is not ative")
	}

	roleField := mapClaim["roles"]
	if roleField == nil {
		return false, fmt.Errorf("error roles field not exist")
	}

	roles, ok := roleField.([]interface{})
	if !ok {
		//log.Logger.Warning("got roleField not array token: ",req)
		return false, fmt.Errorf("error roles field not exist")
	}
	//fmt.Println("My role : ",roles," ",path," ",method)

	for _, r := range roles {
		role := r.(string)
		//fmt.Println("role : ",role," ",path," ",method)
		allowed, err := Enforcer.Enforce(role, path, method)
		if err != nil {
			return false, fmt.Errorf("enforce error : %w", err)
		}
		if allowed {
			return true, nil
		}
	}
	return false, nil
}
