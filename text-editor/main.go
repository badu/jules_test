package main

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/gltext"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype/truetype"
)

const (
	width  = 640
	height = 480
)

var (
	textBuffer strings.Builder
	filePath   string
)

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Text Editor", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize gl:", err)
	}

	gl.Viewport(0, 0, int32(width), int32(height))

	// Set up callbacks
	window.SetCharCallback(charCallback)
	window.SetKeyCallback(keyCallback)

	// Read font file
	fontBytes, err := ioutil.ReadFile("NotoSans-Regular.ttf")
	if err != nil {
		log.Fatalln(err)
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatalln(err)
	}

	scale := float32(24)
	runes := make([]rune, 0)
	for i := 32; i < 127; i++ {
		runes = append(runes, rune(i))
	}

	textRenderer, err := gltext.NewTruetype(font, scale, runes)
	if err != nil {
		log.Fatalln(err)
	}
	defer textRenderer.Release()

	projection := mgl32.Ortho2D(0, float32(width), 0, float32(height))
	textRenderer.SetProjection(projection)

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	// File I/O
	if len(os.Args) > 1 {
		filePath = os.Args[1]
		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Println("Could not read file, starting with empty buffer:", err)
		} else {
			textBuffer.WriteString(string(fileBytes))
		}
	} else {
		textBuffer.WriteString("Hello, World!") // Initial text
	}

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		textRenderer.SetColor(0.0, 0.0, 0.0, 1.0)
		err = textRenderer.Printf(10, float32(height)-scale, textBuffer.String())
		if err != nil {
			log.Println(err)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func charCallback(w *glfw.Window, r rune) {
	textBuffer.WriteRune(r)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		if mods == glfw.ModControl && key == glfw.KeyS {
			if filePath == "" {
				log.Println("No file specified to save to.")
				return
			}
			err := ioutil.WriteFile(filePath, []byte(textBuffer.String()), 0644)
			if err != nil {
				log.Println("Failed to save file:", err)
			} else {
				log.Println("File saved successfully.")
			}
			return
		}

		switch key {
		case glfw.KeyBackspace:
			if textBuffer.Len() > 0 {
				currentText := textBuffer.String()
				textBuffer.Reset()
				textBuffer.WriteString(currentText[:len(currentText)-1])
			}
		case glfw.KeyEnter:
			textBuffer.WriteRune('\n')
		}
	}
}
