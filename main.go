package main

import (
	"database/sql"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"simple-mysql-mgp-hook/moudle"
	"strings"
	"time"
)
import _ "database/sql"
import _ "github.com/go-sql-driver/mysql"

var path string
var conf *moudle.YamlConfig

func CheckErr(err error) {
	if err != nil {
		panic(err)
		fmt.Println("err:", err)
	}
}
func GetTime() string {
	const shortForm = "2006-01-02 15:04:05"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	//fmt.Println(t)
	return str
}
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func readYamlFile(path string) *moudle.YamlConfig {
	conf := new(moudle.YamlConfig)
	yamlFile, err := ioutil.ReadFile(path)
	//log.Println("yamlFile:", yamlFile)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	// err = yaml.Unmarshal(yamlFile, &resultMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	//log.Println("conf", conf)
	return conf
}

/**
构建数据源连接
*/
func createDB() *sql.DB {
	var dataSourceName string
	dataSourceName += conf.Mysql.User + ":"
	dataSourceName += conf.Mysql.Password + "@"
	dataSourceName += "tcp(" + conf.Mysql.Host + ":" + conf.Mysql.Port + ")" + "/" + conf.Mysql.Name + "?charset=utf8"
	log.Printf("[info] 构造连接信息 : %v", dataSourceName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("[error] 建立数据库失败 : %v", err)
	}
	return db

}

/**
执行查询操作
*/
func doQuery(db *sql.DB) bool {
	err := db.Ping()
	if err != nil {
		log.Printf("[error] 发生异常:%v", err)
		fmt.Printf("error:%v", err)
	}
	var member_status string
	query_error := db.QueryRow("select MEMBER_STATE from performance_schema.replication_group_members where MEMBER_HOST = '" + conf.Node.Name + "' ;").Scan(&member_status)
	if query_error != nil {
		log.Printf("[error] 发生错误: %v ,该节点未查询到记录,请检查配置文件", query_error)
		CheckErr(query_error)
	}

	if member_status != "ONLINE" {
		log.Printf("[error] 当前节点状态:%v ,状态异常执行相关操作", member_status)
		return false
	} else {
		log.Printf("[info] 当前节点状态:%v ,状态正常", member_status)
		return true
	}
}

/**
执行配置命令
*/
func execCommand() {
	list := conf.Heartbeat.Command
	for i := 0; i < len(list); i++ {
		log.Printf("[info] 开始执行命令:%v", list[i])
		cmd := exec.Command("bash", "-c", list[i])
		cmd.Start()
	}
}

func init() {
	flag.StringVar(&path, "path", "default", "config yaml path")

}
func main() {
	flag.Parse() //暂停获取参数
	// 如果没有传路径,则直接读取当前执行文件的相对路径下的config.yml
	fmt.Println("路径配置信息:" + path)
	if path == "default" {
		log.Println("[info] 未设置配置文件,寻找同目录下的config.yml文件")
		path = getCurrentDirectory() + "/config.yml"
		log.Println("[info] 解析配置文件路径" + path)
	}

	// 解析文件
	conf = readYamlFile(path)
	if conf.LogPath.Path == "" {
		log.Printf("[info] 未配置日志路径信息,日志将会输出到:%v", getCurrentDirectory()+"/hook.log")
		conf.LogPath.Path = getCurrentDirectory() + "/hook.log"
	}
	file, _ := os.OpenFile(conf.LogPath.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //打开日志文件，不存在则创建
	defer file.Close()
	log.SetOutput(file) //设置输出流
	//log.SetPrefix("[info]")    //日志前缀
	log.SetFlags(log.Ldate | log.Ltime) //日志输出样式
	//log.Println("Hi file")
	log.Println("[info] 程序启动...")
	// 创建数据源
	db := createDB()
	defer db.Close()
	for i := 0; ; i++ {
		// 如果检测失败了
		if !doQuery(db) {
			// do
			execCommand()
		}
		//心跳间隔
		time.Sleep(time.Duration(conf.Heartbeat.Interval) * time.Second)
	}
}
