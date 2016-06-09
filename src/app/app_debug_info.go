package app

import (
	log "github.com/Sirupsen/logrus"
	_ "net/http/pprof"

	"net/http"
	"runtime"
	"time"
)

func (a *app) debugInfo() {

	go func() {

		log.Info("Worker running in debug mode")
		log.Infof("pprof: http://127.0.0.1:8081/debug/pprof")
		log.Error(http.ListenAndServe("127.0.0.1:8081", nil))
	}()

	tick := time.Tick(time.Minute)

	for {

		var memStats runtime.MemStats

		runtime.ReadMemStats(&memStats)

		log.Debugf(
			"gorutines: %d, num gc: %d, alloc: %d, mallocs: %d, frees: %d, heap alloc: %d, stack inuse: %d",
			runtime.NumGoroutine(),
			memStats.NumGC,
			memStats.Alloc,
			memStats.Mallocs,
			memStats.Frees,
			memStats.HeapAlloc,
			memStats.StackInuse,
		)

		<-tick
	}
}
