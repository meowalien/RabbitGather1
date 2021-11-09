package vc_code

import (
	"context"
	"core/src/conf"
	"core/src/lib/errs"
	"core/src/lib/mail"
	"core/src/lib/random"
	"core/src/module/db/redisdb"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const EmailKey = "email"
const VSCode = "vs_code"


func SaveVCCode(recipientType, key, vcCode string) error {
	timeout := time.Duration(conf.GlobalConfig.VCCode.Timeout)
	if conf.DEBUG_MOD {
		timeout = time.Hour *4
	}
	_, err := redisdb.Conn.Set(context.TODO(), redisdb.FormatKey(VSCode,recipientType, key, vcCode), nil, timeout).Result()
	return err
}

func CheckVCCode(recipientType, key, vcCode string) (ok bool, err error) {
	_, err = redisdb.Conn.Get(context.TODO(), redisdb.FormatKey(VSCode,recipientType, key, vcCode)).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}


var MailSender mail.MailSender

func init() {
	MailSender = mail.NewMailSender(conf.GlobalConfig.SMTP.MailAddr,conf.GlobalConfig.SMTP.UserName,conf.GlobalConfig.SMTP.Password)
}



func SendEmailVCCode(email string) (vc string, err error) {
	vc = random.MaxInt0Fill(conf.GlobalConfig.VCCode.Length)
	err = MailSender.SendMail("RabbitGather1 CV code" ,fmt.Sprintf("-- %s --",vc) , email)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	return
}
