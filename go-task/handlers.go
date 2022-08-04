package gotask

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var userControl *UserControl

func Login(c *gin.Context) {
	var creds User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := userControl.myconnection.DBUserLogin(creds)
	if(len(user.UserID)==0){
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong username or password"})
	}else{
		tokenK,err := GenerateToken(user)
		if err !=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		user.Password=""
		user.Token=tokenK

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func UserGET(c *gin.Context) {

	token := c.MustGet("token").(string)
	session,_:= TokensMap.Get(token)
	if(session.Role=="3"){
		userId := c.Query("userid")
		if(session.userid==userId) {
			user := userControl.myconnection.DBUserGet(userId)
			c.JSON(http.StatusOK, gin.H{"user": user})
		}else{
			c.String(http.StatusUnauthorized, "unauthorized")
		}
	}else{
		user := userControl.myconnection.DBUserGetBulk()
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
func UserPost(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp,_ := userControl.myconnection.DBUserCreate(user)

	c.JSON(resp, gin.H{"resp": resp})
}

func UserUpdate(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp,_ := userControl.myconnection.DBUserUpdate(user)

	c.JSON(resp, gin.H{"resp": resp})
}
func UserDelete(c *gin.Context) {
	userId := c.Params.ByName("userid")
	resp,_ := userControl.myconnection.DBUserDelete(userId)
	c.JSON(resp, gin.H{"resp": resp})
}

func TaskPost(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := c.MustGet("token").(string)
	session,_:= TokensMap.Get(token)
	if(session.Role=="1"){
		resp,_ := userControl.myconnection.DBTaskCreate(task)
		c.JSON(resp, gin.H{"tasks": task})
	}else{
		if(session.userid==task.UserID) {
			resp,_ := userControl.myconnection.DBTaskCreate(task)
			c.JSON(resp, gin.H{"tasks": task})
		}else{
			fmt.Println("IDs not matching")
			c.String(http.StatusUnauthorized, "unauthorized")
		}
	}
}


func TaskDelete(c *gin.Context) {
	taskid := c.Params.ByName("taskid")
	token := c.MustGet("token").(string)
	session,_:= TokensMap.Get(token)
	if(session.Role=="1"){
		resp,_ := userControl.myconnection.DBTaskDelete(taskid)
		c.String(resp, "%d",resp)
	}else{
		task := userControl.myconnection.DBTaskGetOne(taskid)
		if(session.userid==task.UserID) {
			resp,_ := userControl.myconnection.DBTaskDelete(taskid)
			c.String(resp, "%d",resp)
		}else{
			c.String(http.StatusUnauthorized, "unauthorized")
		}
	}
}
func TaskUpdate(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := c.MustGet("token").(string)
	session,_:= TokensMap.Get(token)
	if(session.Role=="1"){
		resp,_ := userControl.myconnection.DBTaskUpdate(task)
		c.JSON(resp, gin.H{"tasks": task})
	}else{
		if(session.userid==task.UserID) {
			resp,_ := userControl.myconnection.DBTaskUpdate(task)
			c.JSON(resp, gin.H{"tasks": task})
		}else{
			c.String(http.StatusUnauthorized, "unauthorized")
		}
	}
}


func TaskGet(c *gin.Context) {
	userId := c.Query("userid")
	taskId := c.Query("taskid")
	startdate := c.Query("startdate")
	enddate := c.Query("enddate")

	token := c.MustGet("token").(string)
	session,_:= TokensMap.Get(token)
	//fmt.Println(session)
	if(session.Role=="1"){
		tasks := userControl.myconnection.DBTaskGetBulk(startdate,enddate)
		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	}else{
		if(len(userId)>0 && len(taskId)==0){
			if _, err := strconv.Atoi(userId); err == nil {
				if(session.userid==userId){
					tasks := userControl.myconnection.DBTaskUserGetBulk(userId,startdate,enddate)
					c.JSON(http.StatusOK, gin.H{"tasks": tasks})
				}else{
					c.String(http.StatusUnauthorized, "unauthorized")
				}
			}else{
				c.String(http.StatusBadRequest, "Please pass a valid param")
			}
		}else if(len(taskId)>0 && len(userId)>0){
			if _, err := strconv.Atoi(taskId); err == nil {
				if(session.userid==userId) {
					task := userControl.myconnection.DBTaskGetOne(taskId)
					c.JSON(http.StatusOK, gin.H{"tasks": task})
				}else{
					c.String(http.StatusUnauthorized, "unauthorized")
				}
			}else{
				c.String(http.StatusBadRequest, "Please pass a valid param")
			}
		}else{
			c.String(http.StatusBadRequest, "Please pass a valid param")
		}
	}

}
