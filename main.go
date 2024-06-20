package main

import (
	"MiiContestChannelServer/middleware"
	"MiiContestChannelServer/webpanel"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

var (
	ctx  = context.Background()
	pool *pgxpool.Pool
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("Everybody Votes Channel server has encountered a fatal error! Reason: %v\n", err)
	}
}

func main() {
	// Get config
	config := GetConfig()

	// Start SQL
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Username, config.Password, config.DatabaseAddress, config.DatabaseName)
	dbConf, err := pgxpool.ParseConfig(dbString)
	checkError(err)
	pool, err = pgxpool.ConnectConfig(ctx, dbConf)
	checkError(err)

	// Load salt
	salt, err := os.ReadFile("salt.bin")
	checkError(err)

	// Set up HTTP
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	panel := webpanel.WebPanel{
		Pool:   pool,
		Ctx:    ctx,
		Salt:   salt,
		Config: config,
	}

	r.GET("/panel/login", panel.LoginPage)
	r.POST("/panel/login", panel.Login)

	auth := r.Group("/panel")
	auth.Use(middleware.AuthenticationMiddleware())
	{
		auth.GET("/admin", panel.AdminPage)
		auth.GET("/contests", panel.ViewContests)
		auth.GET("/contests/add", panel.AddContest)
		auth.POST("/contests/add", panel.AddContestPOST)
	}

	r.POST("/cgi-bin/conpost.cgi", conPost)
	r.POST("/cgi-bin/convote2.cgi", conVote)
	r.POST("/cgi-bin/conresult.cgi", conResult)
	r.POST("/cgi-bin/mister.cgi", mister)
	r.POST("/cgi-bin/check.cgi", check)
	r.POST("/cgi-bin/post.cgi", post)
	r.POST("/cgi-bin/vote.cgi", vote)

	fmt.Printf("Starting HTTP connection (%s)...\nNot using the usual port for HTTP?\nBe sure to use a proxy, otherwise the Wii can't connect!\n", config.Address)
	log.Fatalln(r.Run(config.Address))
}
