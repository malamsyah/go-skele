package main

import (
	"github.com/malamsyah/go-skele/internal/handler"
	"github.com/malamsyah/go-skele/pkg/config"
)

func main() {
	r := handler.SetupRouter()
	err := r.Run(":" + config.Instance().AppPort)
	if err != nil {
		panic(err)
	}
}
