package main

import (
	"github.com/gDreamcatcher/dockerMonitor/api"
	"github.com/gin-gonic/gin"
)

var (
	containerList = "/container/list"
	containerStat = "/container/stat"
	containerCpu  = "/container/cpu"
	containerMem  = "/container/memory"
)

func server() {
	r := gin.Default()
	r.GET(containerList, api.ContainerList)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func main() {
	server()
}
