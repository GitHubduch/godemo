package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	WebServerBase()
}

type BaseJsonBean struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewBaseJsonBean() *BaseJsonBean {
	return &BaseJsonBean{}
}

func WebServerBase() {
	fmt.Println("This is webserver base!")

	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/login", loginTask)

	//服务器要监听的主机地址和端口号
	err := http.ListenAndServe("localhost:8090", nil)

	if err != nil {
		fmt.Println("ListenAndServe error: ", err.Error())
	}
}

func loginTask(w http.ResponseWriter, req *http.Request) {
	fmt.Println("loginTask is running...")

	//模拟延时
	time.Sleep(time.Second * 2)

	//读取到json数据
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))

	//获取客户端通过GET/POST方式传递的参数
	req.ParseForm()
	param_userName, found1 := req.Form["userName"]
	param_passWord, found2 := req.Form["passWord"]
	param_Opt, found3 := req.Form["Opt"]

	if !(found1 && found2 && found3) {
		fmt.Fprint(w, "请勿非法访问")
		return
	}

	userName := param_userName[0]
	passWord := param_passWord[0]

	//db 是一个*sql.DB类型的指针，在后面的操作中，都要用到db
	db, err := sql.Open("mysql", "root:duch123@tcp(127.0.0.1:3306)/?charset=utf8") //第一个参数为驱动名
	checkErr(err)
	fmt.Println("DB open successful!")

	if param_Opt[0] == "login" {
		userlogin(userName, passWord, db, w)
	} else if param_Opt[0] == "regist" {
		userregist(userName, passWord, db, w)
	}
}
func userlogin(userName string, passWord string, db *sql.DB, w http.ResponseWriter) {
	var username, password string
	result := NewBaseJsonBean()
	res := db.QueryRow("select * from godb.login where user=?", userName)

	err := res.Scan(&username, &password)
	checkErr(err)

	if userName == username && passWord == password {
		result.Code = 100
		result.Message = userName + "登录成功"
	} else {
		result.Code = 101
		result.Message = "用户名或密码不正确"
	}
	result.Data = time.Now() //获取系统时间

	//向客户端返回JSON数据
	bytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Fprint(w, string(bytes))
}

func userregist(userName string, passWord string, db *sql.DB, w http.ResponseWriter) {
	result := NewBaseJsonBean()
	var name string

	err := db.QueryRow("SELECT user FROM godb.login where user=?", userName).Scan(&name)
	if name != "" {
		result.Code = 201
		result.Message = "用户名已经存在"
	} else {
		db.Exec("insert into godb.login(user, password) values(?, ?)", userName, passWord)
		result.Code = 200
		result.Message = userName + "注册成功"
	}
	result.Data = time.Now() //获取系统时间

	//向客户端返回JSON数据
	bytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Fprint(w, string(bytes))
}

func checkErr(errMasg error) {
	if errMasg != nil {
		panic(errMasg)
	}
}
