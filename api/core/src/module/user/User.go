package user

import (
	"context"
	"core/src/conf"
	"core/src/lib/broker"
	"core/src/lib/errs"
	"core/src/lib/password"
	"core/src/lib/uuid"
	"core/src/module/db/mariadb"
	"core/src/module/db/redisdb"
	"database/sql"
	"fmt"
	"time"
)

// 用戶類型
type UserType string

const (
	// 系統
	SystemUserType UserType = "sys"
	// 普通玩家
	PlainUserType UserType = "plain"
)

type CreateNewUserRequest struct {
	Type     UserType
	Account  string
	Password string
	Rule     []string
	NickName string
	Email    string
	Phone    string
	Gender   int8
	PhotoURL string
}

func NewUUID() string {
	return uuid.NewUUID("U")
}

type Rule struct {
	Name        string
	Description string
	Active      bool
}

func CreateNewUser(newUser CreateNewUserRequest) (us User, err error) {
	newUUID := NewUUID()
	passwordHash, passwordSalt, err := password.HashPassword(newUser.Password, conf.GlobalConfig.Pepper, password.RecommendSaltLength)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	fmt.Println("passwordHash: ",passwordHash)
	fmt.Println("passwordSalt: ",passwordSalt)
	var allRule []Rule
	for _, r := range newUser.Rule {
		allRule = append(allRule, Rule{
			Name:   r,
			Active: true,
		})
	}

	var tx *sql.Tx
	tx, err = mariadb.Conn.Begin()
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	defer func(tx *sql.Tx) {
		e := tx.Rollback()
		if e != nil && e != sql.ErrTxDone {
			err = fmt.Errorf(" %w and %s", err, e.Error())
			return
		}
	}(tx)
	var userID int64
	userID, err = insertUser(tx, newUser.Type, newUUID, newUser.Account, passwordHash, passwordSalt)
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	err = insertRule(tx, userID, allRule)
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	if newUser.NickName == "" {
		newUser.NickName = newUser.Account
	}
	err = insertUserInfo(tx, userID, newUser.NickName, newUser.Email, newUser.Phone, newUser.Gender, newUser.PhotoURL)
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	bs := user{
		id:   uint64(userID),
		uuid: newUUID,
	}
	switch newUser.Type {
	case PlainUserType:
		us, err = MakePlainUser(bs)
		if err != nil {
			err = errs.WithLine(err)
			return
		}
		return
	default:
		err = fmt.Errorf("not supported user type: %s", newUser.Type)
		return
	}
}

func MakePlainUser(bs user) (us PlainUser, err error) {
	return &plainUser{user: &bs}, nil
}

type PlainUser interface {
	User
}

type plainUser struct {
	*user
}

func insertUserInfo(tx *sql.Tx, userID int64, nickName string, email string, phone string, gender int8, photoURL string) (err error) {
	//var res sql.Result
	_, err = tx.Exec("insert into user_info(user_id , nick_name,email,phone , gender , photo_url)value (?,?,?,?,?,?);", userID, nickName, email, phone, gender, photoURL)
	return
}

func insertRule(tx *sql.Tx, userID int64, allRule []Rule) (err error) {
	for _, role := range allRule {
		_, err = tx.Exec("insert into user_role(user_id,role_id) value(?,(select id from role where name = ? limit 1));", userID, role.Name)
		if err != nil {
			return
		}
	}
	return
}

func insertUser(tx *sql.Tx, userType UserType, newUUID string, account string, hash string, salt string) (userID int64, err error) {
	var res sql.Result
	res, err = tx.Exec("insert into user (type, uuid, account, password_hash ,password_salt) value (?,?,?,?,?);", string(userType), newUUID, account, hash, salt)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	userID, err = res.LastInsertId()
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	return
}
func GetUserByEmail(email string) (us User, exist bool, err error) {
	var bs user
	err = mariadb.Conn.QueryRow("select type, id, uuid from user u join user_info ui on u.id = ui.user_id where ui.email = ? limit 1", email).Scan(&bs.ustype, &bs.id, &bs.uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errs.WithLine(err)
		return
	}
	exist = true
	us, err = makeUser(bs)
	if err != nil{
		err  = errs.WithLine(err)
		return
	}
	return
}

func GetUserByAccount(account string) (us User, exist bool, err error) {
	var bs user
	err = mariadb.Conn.QueryRow("select type, id, uuid from user where account = ?", account).Scan(&bs.ustype, &bs.id, &bs.uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errs.WithLine(err)
		return
	}
	exist = true
	us, err = makeUser(bs)
	if err != nil{
		err  = errs.WithLine(err)
		return
	}
	return
}

//const Root = "root"

func makeUser(bs user) (User, error) {
	switch bs.ustype {
	case string(PlainUserType):
		return &plainUser{user: &bs}, nil
	default:
		return nil, errs.WithLine(fmt.Errorf("not supported user type : %s", bs.ustype))
	}
}

func LogoutUserByUUID(uuid string) error {
	KickOnline(uuid)
	return nil
}

// 踢掉線上的使用者
func KickOnline(uuid string) {
	EventBroker.Publish(UserEvent{UserUUID: uuid})
}

type UserEvent struct {
	UserUUID string
	Msg      interface{}
}

var EventBroker = broker.NewBroker(nil)

func init() {
	go EventBroker.Start()
	fmt.Println("start UserEventBroker")
}

const UserOwnToken = "UserOwnToken"

func RemoveUserOwnTokenByUUID(uuid string) error {
	_, err := redisdb.Conn.Del(context.TODO(), redisdb.FormatKey(UserOwnToken, uuid)).Result()
	if err != nil {
		return fmt.Errorf("error whem set DeleteUserOwnTokenByUUID on redis: %w", err)
	}
	return nil
}

// 紀錄使用者 - token關係
func saveUserOwnToken(userUUID string, tokenID string) error {
	_, err := redisdb.Conn.Set(context.TODO(), redisdb.FormatKey(UserOwnToken, userUUID), tokenID, time.Duration(conf.GlobalConfig.JWT.TokenExpiresAt)).Result()
	if err != nil {
		return fmt.Errorf("error whem set UserOwnToken to redis: %w", err)
	}
	return nil
}
func (u *user) UUID() string {
	return u.uuid
}

func (u *user) UserClaim() map[string]interface{} {
	return map[string]interface{}{
		"user_uuid": u.uuid,
	}
}
func (u *user) CleanUserOwnToken() error {
	return RemoveUserOwnTokenByUUID(u.uuid)
}

func (u *user) SaveUserOwnToken(tokenID string) error {
	return saveUserOwnToken(u.uuid, tokenID)
}

//type rootUser struct {
//	user
//}

//type RootUser interface {
//	User
//}

type user struct {
	id     uint64
	uuid   string
	ustype string
}

func (u *user) Frozen() (exist bool, err error) {
	var fs int
	err = mariadb.Conn.QueryRow("select frozen from user where id = ?", u.id).Scan(&fs)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errs.WithLine(err)
		return
	}
	return fs == 1, nil
}

func (u *user) CheckPassword(ps string) (ok bool, err error) {
	var passHash string
	var passSalt string
	err = mariadb.Conn.QueryRow("select password_hash , password_salt from user where id = ?", u.id).Scan(&passHash, &passSalt)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errs.WithLine(err)
		return
	}

	fmt.Println("passHash: ",passHash)
	fmt.Println("passSalt: ",passSalt)
	ok = password.CheckPasswordHash(ps, passHash, conf.GlobalConfig.Pepper, passSalt)
	return
}

func (u *user) Logout() error {
	return LogoutUserByUUID(u.uuid)
}

type User interface {
	Frozen() (bool, error)
	CheckPassword(password string) (bool, error)
	Logout() error
	UUID() string
	UserClaim() map[string]interface{}
	CleanUserOwnToken() error
	SaveUserOwnToken(tokenID string) error
}
