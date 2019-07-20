package master

import (
	"encoding/json"
	"fmt"
	"github.com/c0ding/crontab/common"

	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

//配置单例 4-1， 因为是模块化开发，所以用单例来访问模块
var (
	G_apiServer *ApiServer
)

// 保存任务到etcd接口 6 。那么要定义 任务类型； 连接etcd，把任务传给etcd 保存
// eg:  POST job={"name": "job1", "command": "echo hello", "cronExpr": "* * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	// 2, 取表单中的job字段
	postJob = req.PostForm.Get("job")
	// 3, 反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	fmt.Println(job)

	// 4, 保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {

		resp.Write(bytes)
	}
	return

ERR:
	// 6, 返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux //路由对象
		listener   net.Listener
		httpServer *http.Server
	)

	// 配置路由 1
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	// 启动 TCP 监听 2
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建一个 HTTP 服务 3 ，把这个服务配置到自己定义的结构体中
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Second,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Second,
		Handler:      mux,
	}

	// 赋值单例 4
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 启动了服务端 5
	go httpServer.Serve(listener)
	return

}
