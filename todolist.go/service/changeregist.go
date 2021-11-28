package service

import (
	"net/http"
	"crypto/sha256"
	//"fmt"
	"encoding/hex"
	"strings"
	//"strconv"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)


//登録変更画面
func ReRegist(ctx *gin.Context) {
	uid, _ := ctx.GetPostForm("UID")
	ctx.HTML(http.StatusOK, "changeaccount.html", gin.H{"Title": "REGIST","UID":uid})
}

//生年月日の変更
func ChangeBirthRegist(ctx *gin.Context) {

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	birth, _ := ctx.GetPostForm("birth")

	uid, _ := ctx.GetPostForm("UID")


	data := map[string]interface{}{"id" :uid,"birth":birth} 
		db.NamedExec(`UPDATE accounts SET birth=:birth  WHERE id=:id`,data)//生年月日を変更
		ctx.HTML(http.StatusOK,  "changeaccount.html", gin.H{"Title": "REGIST","UID":uid,"Message":"生年月日が変更されました"})
		return


}

//ユーザ名の変更
func ChangeNameRegist(ctx *gin.Context) {

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	namepre, _ := ctx.GetPostForm("namepre")
	name, _ := ctx.GetPostForm("name")

	uid, _ := ctx.GetPostForm("UID")

	var account database.Account
	err = db.Get(&account,"SELECT * FROM accounts WHERE id=? AND name=?", uid,namepre)
	if err != nil {//idとユーザ名の組が一致しない
		ctx.HTML(http.StatusOK, "changeaccount.html", gin.H{"Title": "REGIST","UID":uid,"Message":"元ユーザ名が違います。"})
		return
	}

	if namepre == name {//もとのものと同じ
		ctx.HTML(http.StatusOK, "changeaccount.html", gin.H{"Title": "REGIST","UID":uid,"Message":"元のユーザ名と同じです"})
		return
	}

	err = db.Get(&account,"SELECT * FROM accounts WHERE name=?", name)
	if err == nil {//すでにaccountsの中にユーザ名が存在する
		ctx.HTML(http.StatusOK, "changeaccount.html", gin.H{"Title": "REGIST","UID":uid,"Message":"そのユーザ名はすでに存在します。"})
		return
	} else {
		data := map[string]interface{}{"id" :uid,"name":name} 
		db.NamedExec(`UPDATE accounts SET name=:name  WHERE id=:id`,data)//ユーザ名を変更
		ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Accout list","Re":"ユーザ名が変更されました"})
		return
	}

}

//パスワードの変更
func ChangeRegist(ctx *gin.Context) {

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	passwordpre, _ := ctx.GetPostForm("passwordpre")
	password, _ := ctx.GetPostForm("password1")

	hashpre := sha256.Sum256([]byte(passwordpre))
	hashpasspre:=strings.ToUpper(hex.EncodeToString(hashpre[:]))

	hash := sha256.Sum256([]byte(password))
	hashpass:=strings.ToUpper(hex.EncodeToString(hash[:]))

	uid, _ := ctx.GetPostForm("UID")

	var account database.Account
	err = db.Get(&account,"SELECT * FROM accounts WHERE id=? AND password=?", uid,hashpasspre)
	if err != nil {//idとパスワードの組が一致しない
		ctx.HTML(http.StatusOK, "changeaccount.html", gin.H{"Title": "REGIST","UID":uid,"Message":"元パスワードが違います。"})
		return
	}

	if hashpasspre == hashpass {//パスワードがもとのものと同じ
		ctx.HTML(http.StatusOK, "changeaccount.html", gin.H{"Title": "REGIST","UID":uid,"Message":"元のパスワードと同じです"})
		return
	}

	data := map[string]interface{}{"id" :uid,"password":hashpass} 
	db.NamedExec(`UPDATE accounts SET password=:password  WHERE id=:id`,data)//パスワードを変更

	ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Accout list","Re":"パスワードが変更されました"})
	return
}

//アカウントの削除
func DeleteRegist(ctx *gin.Context) {

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	name, _ := ctx.GetPostForm("name")
	password, _ := ctx.GetPostForm("password")

	hash := sha256.Sum256([]byte(password))
	hashpass:=strings.ToUpper(hex.EncodeToString(hash[:]))

	uid, _ := ctx.GetPostForm("UID")

	var account database.Account
	err = db.Get(&account,"SELECT * FROM accounts WHERE id=? AND name=? AND password=?", uid,name,hashpass)
	if err != nil {//id、ユーザ名とパスワードの組が一致しない
		ctx.HTML(http.StatusOK, "changeaccount.html", gin.H{"Title": "REGIST","UID":uid,"Message":"ユーザ名あるいはパスワードが違います。"})
		return
	}

	data := map[string]interface{}{"uid" :uid} 
	db.NamedExec(`DELETE FROM accounts WHERE id=:uid`,data)

	ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Accout list","Re":"ユーザを削除しました。"})

}
