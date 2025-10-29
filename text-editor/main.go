package main

import (
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Text Editor", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		glfw.PollEvents()
	}
}
