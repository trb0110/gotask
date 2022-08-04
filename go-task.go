package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	. "gotask/go-task"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var db = make(map[string]string)
func main() {

	NewUserControl()
	InitializeTokens()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
	}))

	authorized.POST("/login", Login)
	authorized.POST("/user", UserPost)

	taskAccess := authorized.Group("/task",TaskAccessAuth())
	taskAccess.POST("/", TaskPost)
	taskAccess.DELETE("/:taskid", TaskDelete)
	taskAccess.PUT("/", TaskUpdate)
	taskAccess.GET("/", TaskGet)


	userAccess := authorized.Group("/user",UserAccessAuth())
	userAccess.GET("/", UserGET)
	userAccess.DELETE("/:userid", UserDelete)
	userAccess.PUT("/", UserUpdate)



	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func UserAccessAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer:=c.Request.Header["Bearer"]

		if(len(bearer)>0){
			accessExists:=LookupTokenKey(bearer[0])
			fmt.Println("Exists User Access		", accessExists)
			if(accessExists){
				session,_ :=TokensMap.Get(bearer[0])
				fmt.Println("Exists User Access		",session)
				if(session.Role=="1"||session.Role=="2"){
					c.Set("token",bearer[0])
					c.Next()
					return
				}else{
					log.Println("Error in User Auth 1")
					c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
				}
			}else{
				log.Println("Error in User Auth 2")
				c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
			}
		}else{
			log.Println("Error in User Auth 3")
			c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
		}
	}
}
func TaskAccessAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer:=c.Request.Header["Bearer"]
		if(len(bearer)>0){
			accessExists:=LookupTokenKey(bearer[0])
			fmt.Println("Exists Task Access		", accessExists)
			if(accessExists){
				session,_ :=TokensMap.Get(bearer[0])
				fmt.Println("Exists Task Access		", session)
				if(session.Role=="1"||session.Role=="3"){
					c.Set("token",bearer[0])
					c.Next()
					return
				}else{
					log.Println("Error in Task Auth 1")
					c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
				}
			}else{
				log.Println("Error in Task Auth 2")
				c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
			}
		}else{
			log.Println("Error in Task Auth 3")
			c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
		}
	}
}
