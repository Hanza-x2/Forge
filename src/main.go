package main

import (
	"GoForge/graphics"
	"GoForge/graphics/viewports"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
)

type App struct {
	driver   *Driver
	batch    *graphics.Batch
	rock     *graphics.Texture
	viewport viewports.Viewport
}

func (app *App) Create(driver *Driver) {
	app.driver = driver

	shader, err := graphics.NewShader(
		`#version 130

in vec4 a_position;
in vec4 a_color;
in vec2 a_texCoord;

uniform mat4 u_projection;

out mediump vec4 v_color;
out highp vec2 v_texCoords;

void main() {
    v_color = a_color;
    v_color.a *= (255.0 / 254.0);
    v_texCoords = a_texCoord;
    gl_Position = u_projection * a_position;
}`,
		`#version 130

in mediump vec4 v_color;
in highp vec2 v_texCoords;

uniform highp sampler2D u_texture;

out vec4 fragColor;

void main() {
    fragColor = v_color * texture(u_texture, v_texCoords);
}`,
	)
	if err != nil {
		log.Fatalf("Failed to create shader: %v", err)
	}

	app.rock, err = graphics.NewTexture("assets/rock.png")
	if err != nil {
		log.Fatalf("Failed to load texture: %v", err)
	}
	app.rock.Bind(0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	app.batch = graphics.NewBatch(shader)
	app.viewport = viewports.NewFitViewport(16, 9, driver.Width, driver.Height)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	log.Println("Application created")
}

func (app *App) Render(driver *Driver, delta float32) {
	gl.ClearColor(0.1, 0.2, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	app.viewport.Apply(true)
	projection := app.viewport.GetCamera().Matrix
	app.batch.SetProjection(projection)

	app.batch.Begin()
	app.batch.Draw(app.rock, 6, 2.5, 4, 4)
	app.batch.End()

	if driver.Input.IsKeyPressed(glfw.KeyEscape) {
		driver.Stop()
	}
}

func (app *App) Resize(driver *Driver, width, height int) {
	app.viewport.Update(width, height, true)
	log.Printf("Application resized to %d x %d", width, height)
}

func (app *App) Destroy(driver *Driver) {
	app.rock.Dispose()
	app.batch.Dispose()
	log.Println("Application destroyed")
}

func init() {
	runtime.LockOSThread()
}

func main() {
	configuration := DefaultConfig()
	configuration.Width = 800
	configuration.Height = 450
	if err := RunSafe(&App{}, configuration); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
