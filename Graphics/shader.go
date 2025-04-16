package Graphics

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Shader struct {
	Program          uint32
	uniformLocations map[string]int32
}

func NewShader(vertexSource, fragmentSource string) (*Shader, error) {
	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return nil, fmt.Errorf("vertex shader: %v", err)
	}

	fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, fmt.Errorf("fragment shader: %v", err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength+1)
		gl.GetProgramInfoLog(program, logLength, nil, &log[0])

		return nil, fmt.Errorf("shader link failed: %v", string(log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	shader := &Shader{
		Program:          program,
		uniformLocations: make(map[string]int32),
	}

	shader.uniformLocations["u_projection"] = shader.GetUniformLocation("u_projection")
	shader.uniformLocations["u_texture"] = shader.GetUniformLocation("u_texture")

	return shader, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength+1)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])

		return 0, fmt.Errorf("shader compile failed: %v", string(log))
	}

	return shader, nil
}

func (s *Shader) Bind() {
	gl.UseProgram(s.Program)
}

func (s *Shader) Unbind() {
	gl.UseProgram(0)
}

func (s *Shader) GetUniformLocation(name string) int32 {
	if loc, ok := s.uniformLocations[name]; ok {
		return loc
	}
	loc := gl.GetUniformLocation(s.Program, gl.Str(name+"\x00"))
	s.uniformLocations[name] = loc
	return loc
}

func (s *Shader) SetUniform1i(name string, value int32) {
	loc := s.GetUniformLocation(name)
	gl.Uniform1i(loc, value)
}

func (s *Shader) SetUniform1f(name string, value float32) {
	loc := s.GetUniformLocation(name)
	gl.Uniform1f(loc, value)
}

func (s *Shader) SetUniformMatrix4fv(name string, value *float32) {
	loc := s.GetUniformLocation(name)
	gl.UniformMatrix4fv(loc, 1, false, value)
}

func (s *Shader) Dispose() {
	gl.DeleteProgram(s.Program)
}
