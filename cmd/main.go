package main

import (
	"irj/internal"
	"irj/pkg/framework"
)

func main() {
	framework.Run(internal.InitEnv)
}
