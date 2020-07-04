package main

import (
	"encoding/json"
        "fmt"
        "io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// data struct
type Todo struct {
	gorm.Model
	Tag       string
	Content   string
	TimeLimit string
}

// database connect

type DbConfig struct {
        User     string `json:"user"`
        Password string `json:"password"`
        Host     string `json:"host"`
        Port     string `json:"port"`
        Dbname   string `json:"dbname"`
        Sslmode  string `json:"sslmode"`
}

var db *gorm.DB

// connect db

func (d DbConfig) Connect() string {
        return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", d.User, d.Password, d.Host, d.Port, d.Dbname, d.Sslmode)
}

func db_connect() *gorm.DB {
        raw, err := ioutil.ReadFile("test.json")
        if err != nil {
                panic(err)
        }

        db_info := new(DbConfig)
        if err = json.Unmarshal([]byte(raw), db_info); err != nil {
                panic(err)
        }

        db, err := gorm.Open("postgres", db_info.Connect())
        if err != nil {
                panic(err)
        }

        return db
}


// init db
func db_init(db *gorm.DB) {
	if err := db.DB().Ping(); err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todo{})

	//defer db.Close()
}

// create
func db_create(db *gorm.DB, tag string, content string, timelimit string) {
	db.Create(&Todo{Tag: tag, Content: content, TimeLimit: timelimit})
	//defer db.Close()
}

// query
func db_query_all(db *gorm.DB) []Todo {
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	//defer db.Close()

	return todos
}

func db_query(db *gorm.DB, id int) Todo {
	var todo Todo
	db.First(&todo, id)
	//db.Close()

	return todo
}

// update
func db_update(db *gorm.DB, id int, tag string, content string, timelimit string){
	var todo Todo
	db.First(&todo, id)
	todo.Tag = tag
	todo.Content = content
	todo.TimeLimit = timelimit
	db.Save(&todo)

	//defer db.Close()
}


// delete
func db_delete(db *gorm.DB, id int) {
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)

	//defer db.Close()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	db = db_connect()
	db_init(db)

	// default page
	r.GET("/", func(ctx *gin.Context) {
		todos := db_query_all(db)
		ctx.HTML(200, "index.html", gin.H{
			"todos": todos,
		})
	})

	// create
	r.POST("/new", func(ctx *gin.Context) {
		tag := ctx.PostForm("tag")
		content := ctx.PostForm("content")
		timelimit := ctx.PostForm("timelimit")
		db_create(db, tag, content, timelimit)

		ctx.Redirect(302, "/")
	})


	//update
	r.GET("/update/:id", func(ctx *gin.Context){
		param := ctx.Param("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			panic(err)
		}
		todo := db_query(db, id)
		ctx.HTML(200, "update.html", gin.H{"todo": todo})
	})

	r.POST("/update/:id", func(ctx *gin.Context){
		param := ctx.Param("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			panic(err)
		}
		tag := ctx.PostForm("tag")
		content := ctx.PostForm("content")
		timelimit := ctx.PostForm("timelimit")
		db_update(db, id, tag, content, timelimit)

		ctx.Redirect(302, "/")
	})


	// delete
	r.GET("/delete/:id", func(ctx *gin.Context) {
		param := ctx.Param("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			panic(err)
		}

		todo := db_query(db, id)
		ctx.HTML(200, "delete.html", gin.H{"todo": todo})
	})

	r.POST("/delete/:id", func(ctx *gin.Context) {
		param := ctx.Param("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			panic(err)
		}

		db_delete(db, id)

		ctx.Redirect(302, "/")
	})

	r.Run()
}
