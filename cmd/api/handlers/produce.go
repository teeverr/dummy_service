package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/teeverr/dummy_service/internal/database"
	"github.com/teeverr/dummy_service/internal/domain"
)

type WorkloadHandler struct {
	Db     *database.Client
	Config *domain.Config
}

func (wlh *WorkloadHandler) SetCpuWorkload(w http.ResponseWriter, r *http.Request) {
	var ct domain.Workload
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&ct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%s\n", err.Error())
		return
	}
	fmt.Printf("%+v\n", ct)
	err = wlh.Db.SetCpu(&ct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\n", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("target cpu worload was set to %v", ct.TargetCPULoad)
}

func (wlh *WorkloadHandler) ShowTargetCPuWorkload(w http.ResponseWriter, r *http.Request) {
	workload, err := wlh.Db.GetLastCpu()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err.Error())
	}
	if workload.TargetCPULoad == 0 {
		workload.TargetCPULoad = wlh.Config.Workload.Cpu.Min
	}
	err = json.NewEncoder(w).Encode(&workload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err.Error())
	}
}
