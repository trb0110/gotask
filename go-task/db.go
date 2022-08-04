package gotask

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"net/http"
)

var server = flag.String("mssql", "127.0.0.1", "the database server")
var port  = flag.Int("port", 1433, "the database port")
var user = flag.String("user", "sa", "the database user")
var password = flag.String("password", "R13032andompassword", "the database password")
var database = flag.String("database", "gotask", "the database name")

type UserControl struct {
	myconnection *SQLConnection
}

func NewUserControl() {
	UC := &UserControl{
		myconnection: NewDBConnection(),
	}
	userControl=UC
}

type SQLConnection struct {
	originalSession *sql.DB
}

func NewDBConnection() (conn *SQLConnection) {
	conn = new(SQLConnection)
	conn.createLocalConnection()
	return
}


func (c *SQLConnection) createLocalConnection() (err error) {
	log.Println("Connecting to SQL DB server....")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", *server, *user, *password, *port, *database)
	c.originalSession, err = sql.Open("mssql", connString)
	if err == nil {
		log.Println("Connection established to SQL DB server")
		return
	} else {
		log.Println("Error occured while creating SQL DB connection: %s", err.Error())
	}
	// Set maximum number of connections in idle connection pool.
	c.originalSession.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	c.originalSession.SetMaxOpenConns(15)
	return
}

func (c *SQLConnection)  DBUserLogin(creds User)(user *User){

	strSql := `SELECT  [user_id], username,[role_id] ,[password],[prefered_hours] FROM [gotask].dbo.[user_table] where [username] = '`+creds.Username+`'
																							and [password] = '`+creds.Password+`'`

	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	user = &User{}
	for queryResult.Next() {
		err := queryResult.Scan(&user.UserID,&user.Username,&user.Role,&user.Password,&user.PreferredHours)
		if err!=nil{
			//fmt.Println("Check UserName Error	" , err2)
			if err == sql.ErrNoRows{
				return nil
			}

			return nil
		}
	}
	defer queryResult.Close()
	return user
}
func (c *SQLConnection)  DBUserGet(id string)(user *User){

	strSql := `SELECT  [user_id], username,[role_id] ,[prefered_hours] FROM [gotask].dbo.[user_table] where [user_id] = `+id+``

	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	user = &User{}
	for queryResult.Next() {
		err := queryResult.Scan(&user.UserID,&user.Username,&user.Role,&user.PreferredHours)
		if err != nil {
			log.Println(err)
		}
	}
	defer queryResult.Close()
	return user
}
func (c *SQLConnection)  DBUserGetBulk()(user []User){

	strSql := `SELECT  [user_id], username,[role_id] ,[prefered_hours] FROM [gotask].dbo.[user_table]`

	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	users := make([]User,0)
	for queryResult.Next(){
		user := User{}
		err := queryResult.Scan(&user.UserID,&user.Username,&user.Role,&user.PreferredHours)
		if err!=nil{
			if err == sql.ErrNoRows{
				return nil
			}
			return nil
		}
		users= append(users,user)
	}

	defer queryResult.Close()
	return users
}


func (c *SQLConnection)  DBUserCreate(user User)(response int, err error){

	strSql := `insert into [gotask].[dbo].[user_table](username,role_id,password, [prefered_hours])
 				values ('`+user.Username+`','`+user.Role+`','`+user.Password+`','`+user.PreferredHours+`')
		`

	queryResult, rerr := c.originalSession.Query(strSql)
	if rerr != nil {
		return http.StatusInternalServerError, rerr
	}
	defer queryResult.Close()
	return http.StatusOK , nil
}
func (c *SQLConnection)  DBUserUpdate(user User)(response int, err error){

	strSql := `				update [gotask].[dbo].[user_table]
								set
					`

	if user.Username!="" {
		strSql += `[username] = '`+user.Username+`',`
	}

	if user.Password!="" {
		strSql += `[password] = '`+user.Password+`',`
	}
	if user.PreferredHours!="" {
		strSql += `[prefered_hours] = '`+user.PreferredHours+`',`
	}
	if user.Role!="" {
		strSql += `[role_id] = '`+user.Role+`',`
	}

	strSql = strSql[:len(strSql)-1]
	strSql+= ` where [user_id] = `+user.UserID+``

	queryResult, rerr := c.originalSession.Query(strSql)

	if rerr != nil {
		return http.StatusInternalServerError, rerr
	}
	defer queryResult.Close()
	return http.StatusOK , nil

}

func (c *SQLConnection)  DBUserDelete(id string)(response int, err error){
	strSql := `
	DELETE [gotask].[dbo].[task] where [user_id] = `+id+`
	DELETE [gotask].[dbo].[user_table] where [user_id] = `+id+`
`
	queryResult, rerr := c.originalSession.Query(strSql)
	if rerr != nil {
		return http.StatusInternalServerError, rerr
	}
	defer queryResult.Close()
	return http.StatusOK , nil
}


func (c *SQLConnection)  DBTaskGetOne(taskId string)(task *Task){

	strSql := `SELECT  
						[task_id]
							  ,[user_id]
							  ,[task_stamp]
							  ,[task_description]
							  ,[task_duration]
						  FROM [gotask].[dbo].[task]
						  where [task_id] = `+taskId+`
						  
  `


	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	task = &Task{}
	for queryResult.Next() {
		err := queryResult.Scan(&task.TaskId,&task.UserID,&task.Timestamp,&task.Description,&task.Duration)
		if err != nil {
			log.Println(err)
		}
	}
	defer queryResult.Close()
	return task
}

func (c *SQLConnection)  DBTaskGetBulk(startdate string, enddate string)(task []Task){

	strSql := `SELECT 		[task_id]
							  ,[user_id]
							  ,[task_stamp]
							  ,[task_description]
							  ,[task_duration]
					,(select [username]from [gotask].[dbo].[user_table] where[user_id]=t.[user_id]) as [username]

				FROM [gotask].[dbo].[task] t `
	//fmt.Println(strSql)
	if(len(startdate)>0&&len(enddate)==0){
		strSql+= ` where task_stamp>'`+startdate+`'`
	}
	if(len(enddate)>0&&len(startdate)==0){
		strSql+= ` where task_stamp<'`+enddate+`'`
	}
	if(len(enddate)>0&&len(startdate)>0){
		strSql+= ` where task_stamp>'`+startdate+`' and task_stamp<'`+enddate+`'`
	}
	strSql+=`order by task_stamp desc`

	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	tasks := make([]Task,0)
	for queryResult.Next(){
		task := Task{}
		err := queryResult.Scan(&task.TaskId,&task.UserID,&task.Timestamp,&task.Description,&task.Duration,&task.Username)
		if err!=nil{
			if err == sql.ErrNoRows{
				return nil
			}
			return nil
		}
		tasks= append(tasks,task)
	}

	defer queryResult.Close()
	return tasks
}
func (c *SQLConnection)  DBTaskUserGetBulk(userId string, startdate string, enddate string)(task []Task){

	strSql := `SELECT  [task_id]
							  ,[user_id]
							  ,[task_stamp]
							  ,[task_description]
							  ,[task_duration]
					,(select [username]from [gotask].[dbo].[user_table] where[user_id]=t.[user_id]) as [username]

				FROM [gotask].[dbo].[task] t 
				where [user_id] = `+userId+`
		`
	if(len(startdate)>0&&len(enddate)==0){
		strSql+= ` and task_stamp>'`+startdate+`'`
	}
	if(len(enddate)>0&&len(startdate)==0){
		strSql+= ` and task_stamp<'`+enddate+`'`
	}
	if(len(enddate)>0&&len(startdate)>0){
		strSql+= ` and task_stamp>'`+startdate+`' and task_stamp<'`+enddate+`'`
	}

	strSql+=`order by user_id,task_stamp desc`

	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	tasks := make([]Task,0)
	for queryResult.Next(){
		task := Task{}
		err := queryResult.Scan(&task.TaskId,&task.UserID,&task.Timestamp,&task.Description,&task.Duration,&task.Username)
		if err!=nil{
			if err == sql.ErrNoRows{
				return nil
			}
			return nil
		}
		tasks= append(tasks,task)
	}

	defer queryResult.Close()
	return tasks
}
func (c *SQLConnection)  DBTaskCreate(task Task)(response int, err error){

	timestamp := ""
	if(len(task.Timestamp)>0){
		timestamp="'"+task.Timestamp+"'"
	}else{
		timestamp="getdate()"
	}
	strSql := `
				insert into [gotask].[dbo].[task] ([user_id],[task_stamp]
													 ,[task_description]
												     ,[task_duration])	
 				values ( `+ task.UserID+`,`+ timestamp+`,'`+ task.Description+`',`+task.Duration+`)
				`

	//fmt.Println(strSql)
	queryResult, rerr := c.originalSession.Query(strSql)
	if rerr != nil {
		return http.StatusInternalServerError, rerr
	}

	defer queryResult.Close()
	return http.StatusOK , nil
}

func (c *SQLConnection)  DBTaskUpdate(task Task)(response int, err error){

	strSql := `				update [gotask].[dbo].[task]
							set
				`

	if task.Duration!="" {
		strSql += `[task_duration] = `+task.Duration+`,`
	}

	if task.Description!="" {
		strSql += `[task_description] = '`+task.Description+`',`
	}
	if task.Timestamp!="" {
		strSql += `[task_stamp] = '`+task.Timestamp+`',`
	}

	strSql = strSql[:len(strSql)-1]
	strSql+= ` where [task_id] = `+task.TaskId+``

	queryResult, rerr := c.originalSession.Query(strSql)

	if rerr != nil {
		return http.StatusInternalServerError, rerr
	}
	defer queryResult.Close()
	return http.StatusOK , nil
}


func (c *SQLConnection)  DBTaskDelete(id string)(response int, err error){

	strSql := `DELETE  [gotask].[dbo].[task] where [task_id] = `+id+``

	queryResult, rerr := c.originalSession.Query(strSql)
	if rerr != nil {
		return http.StatusInternalServerError, rerr
	}

	defer queryResult.Close()
	return http.StatusOK , nil
}
