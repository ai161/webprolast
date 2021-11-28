package service

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	uid, _ := ctx.GetPostForm("UID")

	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, "SELECT tasks.* FROM tasks , task_owners WHERE task_owners.task_id=tasks.id AND task_owners.user_id=?", uid) // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks,"UID":uid})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	tid, err := strconv.Atoi(ctx.Param("tid"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", tid) // Use DB#Get for one entry
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Render task
	ctx.String(http.StatusOK, task.Title)
}


//タスクの状態の変化
func Change(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	uid, _ := ctx.GetPostForm("UID")

	tid, err := strconv.Atoi(ctx.Param("tid"))
	if err != nil {
 		ctx.String(http.StatusBadRequest, err.Error())
 		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", tid)

	data := map[string]interface{}{"is_done" : !task.IsDone ,"id": tid } 
	db.NamedExec(`UPDATE tasks SET is_done=:is_done WHERE tasks.id=:id`,data)



	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, "SELECT tasks.* FROM tasks , task_owners WHERE task_owners.task_id=tasks.id AND task_owners.user_id=?", uid) // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks,"UID":uid})
}

//タスクの削除
func Delete(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	uid, _ := ctx.GetPostForm("UID")

	tid, err := strconv.Atoi(ctx.Param("tid"))
	if err != nil {
 		ctx.String(http.StatusBadRequest, err.Error())
 		return
	}

	data := map[string]interface{}{"tid" :tid} 
	db.NamedExec(`DELETE FROM tasks WHERE id=:tid`,data)
	db.NamedExec(`DELETE FROM task_owners WHERE task_id=:tid`,data)



	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, "SELECT tasks.* FROM tasks , task_owners WHERE task_owners.task_id=tasks.id AND task_owners.user_id=?", uid) // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks,"UID":uid})
}

//タスクの追加
func Add(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	uid, _ := ctx.GetPostForm("UID")

	
	title, _ := ctx.GetPostForm("title")
	data := map[string]interface{}{"title": title } 
	res,_ :=db.NamedExec(`INSERT INTO tasks (title) VALUES (:title)`,data)
	tid,_ :=res.LastInsertId()
	data2 := map[string]interface{}{"uid": uid ,"tid":tid } 
	db.NamedExec(`INSERT INTO task_owners (user_id,task_id) VALUES (:uid,:tid)`,data2)

	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, "SELECT tasks.* FROM tasks , task_owners WHERE task_owners.task_id=tasks.id AND task_owners.user_id=?", uid) // Use DB#Select for multiple entries
	if err != nil {
		fmt.Println("!!!")
		ctx.String(http.StatusInternalServerError, err.Error())
		fmt.Println("!!!")
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks , "UID":uid})
}

//タスクを検索
func Search(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	uid, _ := ctx.GetPostForm("UID")

	//search_id , _:= ctx.GetQuery("search--id")

	search_word , _ := ctx.GetPostForm("search-word")
	//search_date ,  _ := ctx.GetQuery("search-date")
	search_dostate , nodostate := ctx.GetPostForm("search-dostate")
	search_undostate , noundostate := ctx.GetPostForm("search-undostate")

	var isdo bool = true
	if(search_dostate=="do"){isdo=true}
	if(search_undostate=="undo"){isdo=false}


	// Get tasks in DB
		//search_state := map[string]interface{}{"isdo":isdo, "undo":!undo} 
	var tasks []database.Task
	
	//data := map[string]interface{}{"ida":search_id,"search-word":search_word,"search-date":search_date,"isdo":isdo } 
	if (search_word=="") && ((!nodostate&&!noundostate)||(nodostate&&noundostate)){
		err = db.Select(&tasks,"SELECT tasks.* FROM tasks , task_owners WHERE task_owners.task_id=tasks.id AND task_owners.user_id=?", uid)
	} else if (search_word==""){
		err = db.Select(&tasks,"SELECT tasks.* FROM tasks, task_owners WHERE is_done=? AND task_owners.task_id=tasks.id AND task_owners.user_id=?",isdo,uid)
	} else if ((!nodostate&&!noundostate)||(nodostate&&noundostate)){
		err = db.Select(&tasks,"SELECT tasks.* FROM tasks, task_owners WHERE title like  CONCAT('%', ?, '%') AND task_owners.task_id=tasks.id AND task_owners.user_id=?",search_word, uid)
	} else {
		err = db.Select(&tasks,"SELECT tasks.* FROM tasks, task_owners WHERE is_done=? AND title like  CONCAT('%', ?, '%') AND task_owners.task_id=tasks.id AND task_owners.user_id=?",isdo,search_word,uid)
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}


	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks ,"UID":uid})
}

