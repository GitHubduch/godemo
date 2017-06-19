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
	var username, password string
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

	if !(found1 && found2) {
		fmt.Fprint(w, "请勿非法访问")
		return
	}

	result := NewBaseJsonBean()
	userName := param_userName[0]
	passWord := param_passWord[0]

	db, err := sql.Open("mysql", "root:duch123@tcp(127.0.0.1:3306)/?charset=utf8") //第一个参数为驱动名
	checkErr(err)
	fmt.Println("DB open successful!")

	res := db.QueryRow("select * from godb.login where user=?", userName)

	err = res.Scan(&username, &password)
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
func checkErr(errMasg error) {
	if errMasg != nil {
		panic(errMasg)
	}
}
