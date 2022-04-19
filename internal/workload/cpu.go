package workload

import (
	"github.com/teeverr/dummy_service/internal/database"
	"github.com/teeverr/dummy_service/internal/domain"
	"log"
	"runtime"
	"time"
)

// RunCPULoad run CPU load in specify cores count and percentage
func RunCPULoad(coresCount int, timeSeconds int, percentage int) {
	runtime.GOMAXPROCS(coresCount)

	// second     ,s  * 1
	// millisecond,ms * 1000
	// microsecond,μs * 1000 * 1000
	// nanosecond ,ns * 1000 * 1000 * 1000

	// every loop : run + sleep = 1 unit

	// 1 unit = 100 ms may be the best
	unitHundresOfMicrosecond := 1000
	runMicrosecond := unitHundresOfMicrosecond * percentage
	sleepMicrosecond := unitHundresOfMicrosecond*100 - runMicrosecond
	for i := 0; i < coresCount; i++ {
		go func() {
			runtime.LockOSThread()
			// endless loop
			for {
				begin := time.Now()
				for {
					// run 100%
					if time.Now().Sub(begin) > time.Duration(runMicrosecond)*time.Microsecond {
						break
					}
				}
				// sleep
				time.Sleep(time.Duration(sleepMicrosecond) * time.Microsecond)
			}
		}()
	}
	// how long
	time.Sleep(time.Duration(timeSeconds) * time.Second)
}

func CpuWorkloadReader(db *database.Client, config *domain.Config) {
	tickerTime := 15 * time.Second
	ticker := time.NewTicker(tickerTime)
	for {
		<-ticker.C
		workload, err := db.GetLastCpu()
		if err != nil {
			log.Printf("%s", err.Error())
		}
		var cpuPercentage int
		switch {
		case workload.TargetCPULoad == 0:
			cpuPercentage = config.Workload.Cpu.Min
		case workload.TargetCPULoad > config.Workload.Cpu.Max:
			cpuPercentage = config.Workload.Cpu.Max
		default:
			cpuPercentage = workload.TargetCPULoad
		}
		log.Printf("TARGET CPU%%  - %v", cpuPercentage)
		RunCPULoad(runtime.NumCPU(), 14, cpuPercentage)
	}
}
