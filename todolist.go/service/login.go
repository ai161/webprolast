package service

import (
	"net/http"
	"crypto/sha256"
	"fmt"
	"encoding/hex"
	"strings"
	//"strconv"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)


// Login renders login.html
func Login(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "LOGIN"})
}

//登録画面に遷移
func Register(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "newregist.html", gin.H{"Title": "REGIST"})
}

//パスワード忘れた画面へ
func ForgetRegister(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "forget.html", gin.H{"Title": "FORGET"})
}


//パスワード忘れのときの変更
func Changepass(ctx *gin.Context) {

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	birth, _ := ctx.GetPostForm("birth")
	name, _ := ctx.GetPostForm("name")
	password,_ :=ctx.GetPostForm("password")
	//ハッシュ化
	hash := sha256.Sum256([]byte(password))
	//ハッシュ化したパスワードをstringにする
	hashpass:=strings.ToUpper(hex.EncodeToString(hash[:]))

	var account database.Account
	err = db.Get(&account,"SELECT * FROM accounts WHERE name=? AND birth=?", name,birth)
	if err == nil {//一致する
		data := map[string]interface{}{"password" : hashpass ,"name": name } 
		db.NamedExec(`UPDATE accounts SET password=:password WHERE name=:name`,data)
		ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Accout list","Re":"パスワードを変更しました。"})
		return
	}

	ctx.HTML(http.StatusOK, "forget.html", gin.H{"Title": "Forget list","Text":"生年月日あるいはユーザ名が違います。"})
	return


}

//新規登録
func AddRegister(ctx *gin.Context) {
	//ユーザ名とパスワードを取得
	name, _ := ctx.GetPostForm("name")
	password, _ := ctx.GetPostForm("password")
	//ハッシュ化
	hash := sha256.Sum256([]byte(password))
	//ハッシュ化したパスワードをstringにする
	hashpass:=strings.ToUpper(hex.EncodeToString(hash[:]))
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(hashpass,"\n")
	data := map[string]interface{}{"name": name,"password":hashpass} 
	var account database.Account
	err = db.Get(&account,"SELECT * FROM accounts WHERE name=?", name)
	if err == nil {//すでにaccountsの中にユーザ名が存在する
		ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Accout list","Re":"そのユーザ名はすでに存在します。"})
		return
	}
	db.NamedExec(`INSERT INTO accounts (name,password) VALUES (:name,:password)`,data)

	// Render tasks
	ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Accout list","Re":"ユーザが登録されました"})

}
