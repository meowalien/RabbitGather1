package user

import (
	"core/sec/conf"
	"core/sec/db/mariadb"
	"core/sec/lib/password"
	"core/sec/log"
	"core/sec/module/token"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"time"
)

type User struct {
	Name               string
	ID                 uint
	UUID               string
	passwordHash       string
	salt               string
	roles              []string
	passwordEncryption bool
}

//func (u *User) GetRoles() ([]string, error) {
//	if u.roles != nil {
//		return u.roles, nil
//	}
//	panic("GetRoles  is not implemented")
//}

// 檢查密碼
func (u *User) CheckPassword(passwordGiven string) (bool, error) {
	passwordHashInStorage, salt, passwordEncryption, err := u.GetPasswordHashAndSalt()
	if err != nil {
		return false, fmt.Errorf("error when GetPasswordHashAndSalt: %w", err)
	}
	//fmt.Println(passwordHashInStorage, salt, passwordEncryption, err)
	ok := password.CheckPasswordHash(passwordGiven, passwordHashInStorage, conf.GlobalConfig.Pepper, salt, passwordEncryption)

	return ok, nil
}

const KEYUserUUID = "user_uuid"

// UserIdentificationClaim will be the payload in the JWT token
type UserIdentificationClaim struct {
	jwt.StandardClaims
	UserUUID string   `mapstructure:"user_uuid"`
	Roles    []string `mapstructure:"roles"`
}

func (u UserIdentificationClaim) ToMap() map[string]interface{} {
	var mp map[string]interface{}
	err := mapstructure.Decode(u, &mp)
	if err != nil {
		log.Logger.Error("error when decode UserIdentificationClaim to map")
		return nil
	}
	return mp
}

// 創建身份識別信息
func (u *User) CreateIdentificationClaim() map[string]interface{} {
	var result map[string]interface{}
	input := UserIdentificationClaim{
		StandardClaims: token.CreateDefaultStandardClaims(u.Name),
		UserUUID:       u.UUID,
		Roles:          u.roles,
	}
	err := mapstructure.Decode(input, &result)
	if err != nil {
		log.Logger.Error("error when Decode: ",err.Error())
	}
	return result
}

func (u *User) GetPasswordHashAndSalt() (hash string, salt string, passwordEncryption bool, err error) {
	if u.passwordHash != "" && u.salt != "" {
		return u.passwordHash, u.salt, u.passwordEncryption, nil
	}
	var res struct {
		PasswordHash       string
		PasswordEncryption bool
		PasswordSalt       string
	}

	//fmt.Println("ID: ",u.ID)
	err = mariadb.Conn.QueryRow("SELECT password_hash , password_salt , password_encryption FROM `user` WHERE `user`.`id` = ? AND `user`.`deleted_at` = 0  LIMIT 1;", u.ID).
		Scan(&res.PasswordHash, &res.PasswordSalt, &res.PasswordEncryption)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", false, fmt.Errorf("got nothing from db")
		} else {
			return "", "", false, fmt.Errorf("error when query: %w", err)
		}
	}
	return res.PasswordHash, res.PasswordSalt, res.PasswordEncryption, nil
}

// 儲存最後一次豋入時間
func (u *User) SaveLastLogin(now time.Time) {
	_, err := mariadb.Conn.Exec("update user set last_login = ? where id = ?;", now, u.ID)
	if err != nil {
		log.Logger.Warningf("fail to log last login time: userName=%s, userID=%d , now=%s", u.Name, u.ID, now.String())
		return
	}
}

// 確認使用者是否被停權
func (u *User) IsActive() bool {
	_, err := mariadb.Conn.Exec("select id from user where id = ?  and deleted_at =0 and active = 1", u.ID)
	if err != nil {
		log.Logger.Error("fail to get user active status")
		return false
	}
	return true
}

// 假刪除使用者
func (u *User) RemoveUser() error {
	_, err := mariadb.Conn.Exec("UPDATE `user`SET deleted_at = ? WHERE id  = ? ;", time.Now(), u.ID)
	if err != nil {
		return fmt.Errorf("error when remove user in db: %w", err)
	}
	return nil
}

// 停權使用者
func (u *User) FreezeUser() error {
	_, err := mariadb.Conn.Exec("UPDATE `user`SET active = 0 WHERE id  = ? ;", u.ID)
	if err != nil {
		return fmt.Errorf("error when inactive user in db: %w", err)
	}
	return nil
}

// 變更密碼
func (u *User) ChangePassword(newPassword string) error {
	// 是否加密密碼
	passwordEncryptionOn := conf.GlobalConfig.PasswordEncryption

	var passwordHash, passwordSalt string
	var err error
	if passwordEncryptionOn {
		passwordHash, passwordSalt, err = password.HashPassword(newPassword, conf.GlobalConfig.Pepper, password.RecommendSaltLength)
		if err != nil {
			return  fmt.Errorf("error when hash password %w", err)
		}
	} else {
		passwordHash, passwordSalt = newPassword, ""
	}

	// 預設加密
	pe := 1

	if !passwordEncryptionOn{
		pe = 0
	}
	_, err = mariadb.Conn.Exec("update user set password_hash = ? , password_salt = ? , password_encryption = ? where id = ?;" , passwordHash , passwordSalt,pe , u.ID)
	if err != nil {
		return fmt.Errorf("error when Exec %w", err)
	}
	return nil
}

// GetUserByID will return *User if the given userid exist in system.
func GetUserByID(userid string) (*User, bool, error) {
	res := User{}
	rs, err := mariadb.Conn.Query(`select a.id , a.name as username , c.name as rule
from user a join user_role b on a.deleted_at =0 and a.active = 1 and a.id = b.user_id
    join role c on a.deleted_at = 0 and c.active = 1 and c.id = b.role_id where a.id = ?;`, userid)

	if err != nil {
		return nil, false, fmt.Errorf("error when Query: %w", err)
	}

	for rs.Next() {
		var r string
		err := rs.Scan(&res.ID, &res.Name, &r)
		if err != nil {
			return nil, false, fmt.Errorf("error when scan: %w", err)
		}
		res.roles = append(res.roles, r)
	}

	return &res, true, nil
}

// GetUserByUUID will return *User if the given userid exist in system.
func GetUserByUUID(userUUID string) (*User, bool, error) {
	res := User{}
	rs, err := mariadb.Conn.Query(`select a.id , a.name as username , c.name as rule
from user a join user_role b on a.deleted_at =0 and a.active = 1 and a.id = b.user_id
    join role c on a.deleted_at = 0 and c.active = 1 and c.id = b.role_id where a.uuid = ?;`, userUUID)

	if err != nil {
		return nil, false, fmt.Errorf("error when Query: %w", err)
	}

	for rs.Next() {
		var r string
		err := rs.Scan(&res.ID, &res.Name, &r)
		if err != nil {
			return nil, false, fmt.Errorf("error when scan: %w", err)
		}
		res.roles = append(res.roles, r)
	}
	if res.roles == nil || len(res.roles) ==0{
		return nil, false, nil
	}

	res.UUID = userUUID
	return &res, true, nil
}

type CreateNewUserRequest struct {
	Id          uint
	Name        string
	Role        []string
	Password    string
	Introducer  string
	InitBalance int64
}

func CreateNewUser(newUser CreateNewUserRequest) (*User, error) {
	if newUser.InitBalance < 0 {
		return nil, fmt.Errorf("the InitBalance should not be negative number")
	}

	// 是否加密密碼
	passwordEncryptionOn := conf.GlobalConfig.PasswordEncryption

	var passwordHash, passwordSalt string
	var err error
	if passwordEncryptionOn {
		passwordHash, passwordSalt, err = password.HashPassword(newUser.Password, conf.GlobalConfig.Pepper, password.RecommendSaltLength)
		if err != nil {
			return nil, fmt.Errorf("error when hash password %w", err)
		}
	} else {
		passwordHash, passwordSalt = newUser.Password, ""
	}

	allRoles := []tables.Role{}
	for _, r := range newUser.Role {
		allRoles = append(allRoles, tables.Role{
			Name:   r,
			Active: true,
		})
	}

	tx, err := mariadb.Conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("error when open db transaction %w", err)
	}
	defer tx.Rollback()


	pe := 1
	if !passwordEncryptionOn {
		pe = 0
	}


	uuid := NewUUID()
	var userID int64

	if newUser.Id != 0 {
		_, err := tx.Exec("insert into user (id,name ,password_hash,password_salt,password_encryption,introducer,uuid)value(?,?,?,?,?,?,?);", newUser.Id, newUser.Name,
			passwordHash, passwordSalt, pe, newUser.Introducer, uuid)
		if err != nil {
			return nil, fmt.Errorf("error when insert user record %w", err)
		}
		userID = int64(newUser.Id)


	} else {
		res, err := tx.Exec("insert into user (name ,password_hash,password_salt,password_encryption,introducer , uuid)value(?,?,?,?,?,?);", newUser.Name, passwordHash, passwordSalt, pe, newUser.Introducer,uuid)
		if err != nil {
			return nil, fmt.Errorf("error when insert user record %w", err)
		}
		userID, err = res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("error when get LastInsertId %w", err)
		}


	}

	for _, role := range allRoles {
		_, err = tx.Exec("insert into user_role(user_id,role_id) value(?,(select id from role where name = ? limit 1));", userID, role.Name)
		if err != nil {
			return nil, fmt.Errorf("error when insert user_role: %w", err)
		}
	}

	wt, err := wallet.NewWallet(tx,newUser.InitBalance, wallet.TWDCurrencyName, newUser.Id)
	if err != nil {
		return nil, fmt.Errorf("error when insert user record %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return nil, fmt.Errorf("error when commit: %w", err)
	}
	return &User{
		ID:           uint(userID),
		UUID:         uuid,
		Name:         newUser.Name,
		passwordHash: passwordHash,
		salt:         passwordSalt,
		roles:        newUser.Role,
		wallet:       wt,
	}, nil
}

func NewUUID() string {
	return fmt.Sprintf("U%s", math.Snowflake().Base58())
}

func GetUserByUsername(username string) (*User, bool, error) {
	res, err := mariadb.Conn.Query("select a.id ,a.uuid , a.name as username , c.name as rule from user a join user_role b on a.deleted_at =0 and a.active = 1 and a.id = b.user_id join role c on a.deleted_at = 0 and c.active = 1 and c.id = b.role_id where a.name = ?;", username)
	if err != nil {
		return nil, false, fmt.Errorf("error Query: %w", err)
	}
	resUser := User{}

	for res.Next() {
		var r string
		err := res.Scan(&resUser.ID, &resUser.UUID, &resUser.Name, &r)

		if err != nil {
			return nil, false, fmt.Errorf("error Scan: %w", err)
		}
		resUser.roles = append(resUser.roles, r)
	}
	if resUser.roles == nil || len(resUser.roles) ==0 {
		return nil , false, nil
	}

	return &resUser, true, nil
}

func LogoutUserByUUID(userUUID string) error {
	//if userID <= 0 {
	//	return errors.New("userID should should be <= 1")
	//}
	/*
	   移除狀態
	*/
	return nil
}

func GetUserUUIDInClaim(claim map[string]interface{}) (string, error) {
	u, ok := claim[KEYUserUUID]
	if !ok {
		return "", fmt.Errorf("fail to get userID , user_id field is not exist")
	}
	return fmt.Sprint(u), nil
}

func GetUUIDByUserID(userID int64) (string, bool, error) {
	var uuid string
	err := mariadb.Conn.QueryRow("select uuid from user where id = ? limit 1;", userID).Scan(&uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		} else {
			return "", false, fmt.Errorf("error when QueryRow: ", err.Error())
		}
	}
	return uuid, true, nil
}
