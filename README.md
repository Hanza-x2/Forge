# Forge

A minimal 2D game framework for desktop, written in Go.
Forge is designed to be simple and easy to use, while
still providing the necessary tools for game development.

**It's a work in progress, and is not yet feature complete.**

Forge aims to be a lightweight alternative to larger game engines,
and is built with performance in mind. It's meant to be similar to
the legendary Java framework [LibGDX](https://libgdx.com/) and doesn't
force you to use a specific architecture or design pattern.

## Features

While Forge is still in its early stages, it already has a number
of features, such as:

- Simple and intuitive APIs
- Cross-platform desktop support (Windows, Linux, macOS)
- Advanced 2D rendering using OpenGL (with a simple Batch implementation)
- Advanced input handling methods
- Basic audio support (both Music and Sound)
- Minimal geometry and math library
- Minimal scene graph implementation
- Advanced camera support
- Advanced viewport support
- Bitmap font support
- Support for GL-related entities such as:
    - Shaders
    - Textures
    - FrameBuffers

## Goal

The goal of Forge is to allow both beginners and experienced developers
to create 2D games quickly and easily, without being bound to a specific
architecture or design pattern. The code is written in a way that allows
you to use it as a library, or as a base for your own game engine.

Forge's main focus is on performance and simplicity, we may also introduce
support for more platforms in the future, such as mobile (Android, IOS) and
web (using WebAssembly).

## Getting Started

While nobody is stopping you from using Forge in your own projects,
we suggest you to wait until we reach a more stable version for production-grade
code. The framework is still in its early stages, and we are constantly
adding new features and breaking changes, so the API is not yet stable.
As of now, the framework has reached a point where it is usable, but we
are still working on polishing the API. Please check the [CHANGELOG](CHANGELOG.md)
for more information on the latest changes. The current version of the
framework is `0.0.8`, you may check the tags for other versions.

To get started with Forge, you can use (clone/fork) the [starter-template
repository](https://github.com/ForgeLeaf/ForgeStarterTemplate)
which contains a simple example of how to use the framework.

You may clone the template using this command:

```bash
git clone https://github.com/ForgeLeaf/ForgeStarterTemplate.git
```

### Preparations & Building

In order to run the example, you need to have Go installed on your machine.
You can download Go from the official website: [golang.org](https://golang.org/dl/)

Once you have Go installed, you can CD into the directory and
build the template using:

```bash
go build .
```

Then you may run the executable using the following command
**If you're on Linux or macOS**:

```bash
./ForgeStarterTemplate
```

or, **if you're on Windows**:

```bash
./ForgeStarterTemplate.exe
```

*If you'd like to build the final release executable, we'd recommend
the following command:*

```bash
go build -ldflags='-s -w -H=windowsgui' .
```

*This will build a smaller executable without the console window.*

### The Code

The code doesn't force you to use a specific style, all you
have to supply is an implementation of the `Forge.Application`
interface and an instance of the `DesktopConfiguration` struct.

The main method may look like this:

```go
func main() {
if err := Forge.RunSafe(&Application{}, Forge.DefaultDesktopConfig()); err != nil {
panic(err)
}
}
```

The `RunSafe` method will take care of initializing the framework,
creating the window, and running the main loop. Please make sure to
add this snippet of code to run the context on the main thread:

```go
func init() {
runtime.LockOSThread()
}
```

The final code may look like this:

```go
package main

import (
	"github.com/ForgeLeaf/Forge"
	"runtime"
)

type Application struct{}

func (application *Application) Create(driver *Forge.Driver) {
	// Initialization logic goes here
}

func (application *Application) Render(driver *Forge.Driver, delta float32) {
	// Render logic goes here
}

func (application *Application) Resize(driver *Forge.Driver, width, height float32) {
	// Resize logic goes here
}

func (application *Application) Destroy(driver *Forge.Driver) {
	// Cleanup logic goes here
}

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := Forge.RunSafe(&Application{}, Forge.DefaultDesktopConfig()); err != nil {
		panic(err)
	}
}
```

## Issues & Suggestions

If you encounter any issues while using Forge, please report
them here on GitHub. You can also suggest new features or improvements
to the framework. We are always looking for ways to improve
the framework and make it better for everyone.

## Contributing

If you'd like to contribute to the project, please feel free to
fork the repository and submit a pull request. We welcome
any contributions, whether it's bug fixes, new features, or
documentation improvements. As of now, we don't have a fixed
contribution guide, but we will be adding one in the future.