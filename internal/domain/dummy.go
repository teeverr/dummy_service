package domain

import "gorm.io/gorm"

type Workload struct {
	gorm.Model    `json:"-"`
	TargetCPULoad int `json:"target_cpu_load"`
}
