package g3

import (
	"fmt"
	"image"
	"gl"
)

type openGLGraphicsDevice struct{}

type openGLTexture2D struct {
	tex            gl.Texture
}

type openGLShader struct {
	vertexShader   gl.Shader
	fragmentShader gl.Shader
	program        gl.Program
}

type openGLVertexBuffer struct {
	// TODO: use VBOs
	vertices2 []Vec2
	vertices3 []Vec3
}

type openGLIndexBuffer struct {
	// TODO: use VBOs
	indices []uint32
}

func NewOpenGLGraphicsDevice() GraphicsDevice {
	gl.Enable(gl.DEPTH_TEST)
	return &openGLGraphicsDevice{}
}

func (gd *openGLGraphicsDevice) NewTexture2D(img image.Image) Texture2D {
	rect := img.Bounds()
	width, height := rect.Max.X - rect.Min.X, rect.Max.Y - rect.Min.Y
	id := gl.GenTexture()
	id.Bind(gl.TEXTURE_2D)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
	switch i := img.(type) {
		case *image.RGBA:
			gl.TexImage2D(gl.TEXTURE_2D, 0, 4, width, height, 0, gl.RGBA, &i.Pix[0].R)
			//fmt.Println("teximg2d", i.Pix[0:4], i.Stride, width, height, gl.GetError(), id)
		default:
			panic("unknown format")
	}
	return &openGLTexture2D{id}
}

func (gd *openGLGraphicsDevice) NewShader(vertexShaderStr, fragmentShaderStr string) Shader {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	vertexShader.Source(vertexShaderStr)
	vertexShader.Compile()
	fmt.Println(vertexShader.GetInfoLog())

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragmentShader.Source(fragmentShaderStr)
	fragmentShader.Compile()
	fmt.Println(fragmentShader.GetInfoLog())

	program := gl.CreateProgram()
	program.AttachShader(vertexShader)
	program.AttachShader(fragmentShader)
	program.Link()
	program.Validate()
	fmt.Println(program.GetInfoLog())

	return &openGLShader{vertexShader, fragmentShader, program}
}

func (gd *openGLGraphicsDevice) NewVertexBufferVec2(vertices []Vec2) VertexBuffer {
	return &openGLVertexBuffer{vertices, nil}
}

func (gd *openGLGraphicsDevice) NewVertexBufferVec3(vertices []Vec3) VertexBuffer {
	return &openGLVertexBuffer{nil, vertices}
}

func (gd *openGLGraphicsDevice) NewIndexBuffer(indices []uint32) IndexBuffer {
	return &openGLIndexBuffer{indices}
}

func (gd *openGLGraphicsDevice) SetFillMode(mode int) {
	switch mode {
	case FillSolid:
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	case FillWireFrame:
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	default:
		panic("invalid fill mode")
	}
}

func (gd *openGLGraphicsDevice) SetViewport(x, y, w, h int) {
	gl.Viewport(x, y, w, h)
}

func (gd *openGLGraphicsDevice) SetMatrix(mtype int, m *Matrix4x4) {
	switch mtype {
	case MatrixProjection:
		gl.MatrixMode(gl.PROJECTION)
	case MatrixModelView:
		gl.MatrixMode(gl.MODELVIEW)
	default:
		panic("invalid matrix type")
	}
	//TODO: use glLoadTransposeMatrix for correct memory layout
	tm := m.Transposed()
	gl.LoadMatrixf(&tm.M11)
}

func (dev *openGLGraphicsDevice) SetTexture2D(texture Texture2D, unit uint) {
	gltexture := texture.(*openGLTexture2D)
	gl.ActiveTexture(gl.TEXTURE0 + gl.GLenum(unit))
	if gltexture != nil {
		gltexture.tex.Bind(gl.TEXTURE_2D)
	} else {
		gl.Texture(0).Bind(gl.TEXTURE_2D)
	}
}

func (gd *openGLGraphicsDevice) SetShader(shader Shader) {
	glshader := shader.(*openGLShader)
	glshader.program.Use()
}

func (gd *openGLGraphicsDevice) SetTexCoords(buffer VertexBuffer, index uint) {
	glbuffer := buffer.(*openGLVertexBuffer)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY) // TODO: DisableClientState
	gl.TexCoordPointer(2, 3*4, &glbuffer.vertices2[0].X)
}

func (gd *openGLGraphicsDevice) SetNormals(buffer VertexBuffer) {
	glbuffer := buffer.(*openGLVertexBuffer)
	gl.EnableClientState(gl.NORMAL_ARRAY) // TODO: DisableClientState
	gl.NormalPointer(3*4, &glbuffer.vertices3[0].X)
}

func (gd *openGLGraphicsDevice) SetVertices(buffer VertexBuffer) {
	glbuffer := buffer.(*openGLVertexBuffer)
	gl.EnableClientState(gl.VERTEX_ARRAY) // TODO: DisableClientState
	gl.VertexPointer(3, 3*4, &glbuffer.vertices3[0].X)
}

func (gd *openGLGraphicsDevice) DrawIndexed(buffer IndexBuffer) {
	glBuffer := buffer.(*openGLIndexBuffer)
	gl.DrawElements(gl.TRIANGLES, len(glBuffer.indices), &glBuffer.indices[0])
}

func (gd *openGLGraphicsDevice) Clear() {
	gl.ClearColor(0.0, 0.0, 1.0, 0.5)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (t *openGLTexture2D) Release() {
	t.tex.Delete()
}

func (vb *openGLVertexBuffer) Release() {
}

func (ib *openGLIndexBuffer) Release() {
}

func (p *openGLShader) GetUniformLocation(name string) uint {
	return uint(p.program.GetUniformLocation(name))
}

func (p *openGLShader) SetVec3(location uint, v *Vec3) {
	p.program.Use()
	gl.UniformLocation(location).Uniform3f(v.X, v.Y, v.Z)
}

func (p *openGLShader) SetTexture(location uint, unit uint) {
	p.program.Use();
	gl.UniformLocation(location).Uniform1i(int(unit))
}

func (sh *openGLShader) Release() {
	sh.vertexShader.Delete()
	sh.fragmentShader.Delete()
	sh.program.Delete()
}
