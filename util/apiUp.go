package util

import (
	"fund/entity"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func parse(b []byte, start, end, sep string) (sli []string) {
	str := string(b)
	pos := strings.Index(str, start)
	if -1 != pos {
		str = str[pos+len(start):]
	}
	pos = strings.Index(str, end)
	if -1 != pos {
		str = str[:pos]
	}
	return strings.Split(str, sep)
}

type APIUP interface {
	getUptime() string
	getUrl() string
	getWait() int
	parse([]byte) error
}

var (
	g_consumerTask  []APIUP
	g_consumerMutex sync.Mutex
	g_getTaskCond   *sync.Cond
	g_setTaskCond   *sync.Cond
)

func getTask() (ret []APIUP) {
	g_consumerMutex.Lock()
	defer g_consumerMutex.Unlock()
	for 0 == len(g_consumerTask) {
		g_getTaskCond.Wait()
	}
	ret = g_consumerTask
	g_consumerTask = nil
	g_setTaskCond.Signal()
	return ret
}

func setTask(sli []APIUP) {
	g_consumerMutex.Lock()
	defer g_consumerMutex.Unlock()
	for 0 != len(g_consumerTask) {
		g_setTaskCond.Wait()
	}
	g_consumerTask = sli
	g_getTaskCond.Signal()
}

func Download(task APIUP) {
	url := task.getUrl()
	entity.GetLog().Println("url:", url, "uptime:", task.getUptime())
	time.Sleep(time.Duration(rand.Intn(5)+rand.Intn(5)+2) * time.Second)
	b := entity.Get(url)
	if err := task.parse(b); nil != err {
		entity.GetLog().Println(err)
	}
}
func RUN() {
	g_getTaskCond = sync.NewCond(&g_consumerMutex)
	g_setTaskCond = sync.NewCond(&g_consumerMutex)
	go func() {
		l := len("2020-12-12")
		for {
			tasks := getTask()
			entity.GetLog().Println("update start.")
			for i := 0; i < len(tasks); i++ {
				for j := i + 1; j < len(tasks); j++ {
					if tasks[i].getUptime() > tasks[j].getUptime() {
						tasks[i], tasks[j] = tasks[j], tasks[i]
					}
				}
			}

			now := time.Now()
			for i, task := range tasks {
				upTime := task.getUptime()
				if i < 300 && entity.GetConf().GetForceUpdate() && now.AddDate(0, 0, -10).String()[:l] > upTime || now.AddDate(0, 0, task.getWait()).String()[:l] > upTime {
					Download(task)
				} else {
					break
				}
			}
			entity.GetLog().Println("update end.")
		}
	}()

	go func() {
		ch := time.Tick(4 * time.Hour)
		new(Update).AutoUp()
		for {
			select {
			case <-ch:
				new(Update).AutoUp()
			}
		}
	}()
}
