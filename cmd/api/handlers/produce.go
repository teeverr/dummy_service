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
	Db *database.Client
}

func (wh *WorkloadHandler) SetCpuWorkload(w http.ResponseWriter, r *http.Request) {
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
	err = wh.Db.SetCpu(&ct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\n", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("target cpu worload was set to %v", ct.TargetCPULoad)
}
