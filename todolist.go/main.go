package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")

	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/", service.Login)//ログイン画面へ
	engine.POST("/", service.AddRegister)//新規登録画面へ
	engine.GET("/newregist", service.Register)//新規登録画面からログイン画面へ
	engine.POST("/changepass", service.Changepass)//パスワード忘れ画面からログイン画面は
	engine.GET("/forgetregist", service.ForgetRegister)//パスワード忘れ画面へ
	engine.POST("/home", service.Home)//index.html画面へ
	engine.POST("/backhome", service.BackHome)//task_list.htmlからindex.htm画面へ
	engine.POST("/confirm", service.ReRegist)//登録変更画面へ
	engine.POST("/homebirthchange", service.ChangeBirthRegist)//生年月日変更
	engine.POST("/homenamechange", service.ChangeNameRegist)//ユーザ名変更
	engine.POST("/homechange", service.ChangeRegist)//パスワード変更
	engine.POST("/homedel", service.DeleteRegist)//削除
	engine.POST("/list", service.TaskList)//task_list.html
	engine.POST("/listchange/:tid", service.Change)//タスクの状態を編集
	engine.POST("/listadd", service.Add)//タスクを追加
	engine.POST("/listdel/:tid", service.Delete)//タスクを削除
	engine.POST("/listsearch", service.Search)//タスクを検索
	engine.GET("/task/:tid", service.ShowTask) // ":id" is a parameter

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
