package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tochti/gin-gum/gumauth"
	"github.com/tochti/gin-gum/gumspecs"
	"github.com/tochti/gin-gum/gumwrap"
	"github.com/tochti/smem"
	"github.com/tochti/tasks-lib"
	"gopkg.in/gorp.v1"
)

func main() {

	gumspecs.AppName = "tasks"

	srv := gumspecs.ReadHTTPServer()

	mysql := gumspecs.ReadMySQL()
	sqlDB, err := mysql.DB()
	if err != nil {
		log.Fatal(err)
	}

	db := &gorp.DbMap{
		Db: sqlDB,
		Dialect: gorp.MySQLDialect{
			"InnonDB",
			"UTF8",
		},
	}

	router := gin.New()

	s2tore := smem.NewStore()
	signedIn := gumauth.SignedIn(&s2tore)

	tasksAPI := router.Group("/v1/tasks")
	{
		g := gumwrap.Gorp
		tasksAPI.GET("/", signedIn(g(tasks.ReadAll, db)))
		tasksAPI.GET("/:id", signedIn(g(tasks.ReadOne, db)))
		tasksAPI.POST("/", signedIn(g(tasks.Create, db)))
		tasksAPI.PUT("/:id", signedIn(g(tasks.Update, db)))
		tasksAPI.DELETE("/:id", signedIn(g(tasks.Delete, db)))
	}

	uStore := gumauth.SQLUserStore{db.Db}
	router.POST("/signin/:name/:password", gumauth.SignIn(&s2tore, uStore))

	router.POST("/signup", gumauth.CreateUserSQL(db.Db))

	router.Run(srv.String())
}
