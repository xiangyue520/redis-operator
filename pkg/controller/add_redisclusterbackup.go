package controller

import (
	"github.com/xiangyue520/redis-operator/pkg/controller/redisclusterbackup"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, redisclusterbackup.Add)
}
