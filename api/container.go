package api

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"sync"
	"time"
)

func init(){
	once.Do(New)
}

var d *docker

var once sync.Once

type docker struct {
	dockerClient *client.Client
}

func New(){
	d = &docker{}
	dockerClient, _ := client.NewEnvClient()
	d.dockerClient = dockerClient
}

type Container struct {
	ID    string   `json:"id"`
	CPUUsage float64 `json:cpu_usage`
	MemoryUsage uint64 `json: memory_usage`
}

func ContainerList(c *gin.Context) {
	preStat, err := getStatJson()
	if err != nil {
		c.JSON(-1, err.Error())
		return
	}
	time.Sleep(1 * time.Millisecond)
	afterStat, err := getStatJson()
	if err != nil {
		c.JSON(-1, err.Error())
		return
	}
	var crs []Container
	for id, stats := range afterStat {
		cr := Container{
			ID:   stats.ID,
			MemoryUsage: stats.MemoryStats.Stats["total_rss"],
		}
		if v, ok := preStat[id]; ok {
			cpuUsage := stats.CPUStats.CPUUsage.TotalUsage - v.CPUStats.CPUUsage.TotalUsage
			sysUsage := stats.CPUStats.SystemUsage - v.CPUStats.SystemUsage
			cr.CPUUsage = float64(cpuUsage) / float64(sysUsage) / float64(stats.CPUStats.OnlineCPUs) * 100.0
		} else {
			cr.CPUUsage = -1
		}
		crs = append(crs, cr)
	}

	c.JSON(0, crs)
}

func getStatJson() (map[string]types.StatsJSON, error) {
	statsJsons := make(map[string]types.StatsJSON, 0)
	var opt types.ContainerListOptions
	ctx := context.Background()
	containerLists, err := d.dockerClient.ContainerList(ctx, opt)
	if err != nil {
		return statsJsons, errors.Wrap(err, "get container list failed")
	}
	for _, container := range containerLists {
		stats, err := d.dockerClient.ContainerStats(ctx, container.ID, false)
		if err != nil {
			err = errors.Wrap(err, container.ID + " stats failed")
			return nil, err
		}
		var statsJSON types.StatsJSON
		err = json.NewDecoder(stats.Body).Decode(&statsJSON)
		if err != nil {
			err = errors.Wrap(err, container.ID + " decode failed")
			return nil, err
		}
		statsJsons[container.ID] = statsJSON
	}
	return statsJsons, nil
}
