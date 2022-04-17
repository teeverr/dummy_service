package domain

import "gorm.io/gorm"

type Workload struct {
	gorm.Model
	TargetCPULoad int `json:"target_cpu_load"`
}
