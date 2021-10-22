package permission

import (
	"core/sec/conf"
	"core/sec/db/mariadb"
	"database/sql"
	"fmt"
	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"net/http"
)

const (
	LOGIN = "login"
)

var Enforcer *MyEnforcer

func InitRBAC(db *sql.DB) {
	fmt.Println("InitRBAC ...")
	if db == nil {
		panic("the db is nil")
	}

	a, err := sqladapter.NewAdapter(db, conf.GlobalConfig.RBAC.DriverName, conf.GlobalConfig.RBAC.PolicyTableName)
	if err != nil {
		panic(err)
	}

	enforcer, err := casbin.NewEnforcer(conf.GlobalConfig.RBAC.ModelConfFile, a)
	Enforcer = &MyEnforcer{Enforcer: enforcer}
	if err != nil {
		panic(err)
	}

	// Load the policy from DB.
	if err = Enforcer.LoadPolicy(); err != nil {
		panic("LoadPolicy failed, err: " + err.Error())
	}

	//if conf.DEBUG_MOD {
	//	addPermissionForUser()
	//}

}

// 測試時，自動為角色加上存取權
func addPermissionForUser() {
	_, err := Enforcer.AddPermissionForUser(LOGIN, "/member/name", http.MethodGet)
	if err != nil {
		panic(err)
	}

	_, err = Enforcer.AddPermissionForUser(LOGIN, "/lobby", http.MethodGet)
	if err != nil {
		panic(err)
	}

}

type MyEnforcer struct {
	*casbin.Enforcer
}

func (m *MyEnforcer) AddPermissionForUser(user string, permission ...string) (bool, error) {
	err := AddUser(user)
	if err != nil {
		return false, err
	}
	return m.Enforcer.AddPermissionForUser(user, permission...)
}

func AddUser(user string) error {
	_, err := mariadb.Conn.Exec("insert into role (name , active ) value (? , @a:= ?) on duplicate key UPDATE active=@a;", user, 1)

	if err != nil {
		return err
	}
	return nil
}
