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


// Home renders index.html
func Home(ctx *gin.Context) {
	name, _ := ctx.GetPostForm("name")
	password, _ := ctx.GetPostForm("password")

	fmt.Println(password)
	hash := sha256.Sum256([]byte(password))
	fmt.Println(hash)
	hashpass:=strings.ToUpper(hex.EncodeToString(hash[:]))

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(hashpass,"\n")

	//Get a task with given ID
	var accounts []database.Account
	err = db.Select(&accounts,"SELECT * FROM accounts WHERE name=? AND password=?", name,hashpass)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "index.html",  gin.H{"Title": "HOME", "Accounts": accounts})
}

//indexに戻る
func BackHome(ctx *gin.Context){
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	uid, _ := ctx.GetPostForm("UID")

	// Get tasks in DB
	var accounts []database.Account
	err = db.Select(&accounts, "SELECT * FROM accounts WHERE accounts.id=?",uid)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.HTML(http.StatusOK, "index.html",  gin.H{"Title": "HOME", "Accounts": accounts})

}
