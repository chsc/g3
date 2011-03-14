package g3

import (
	"image"
)

const (
	FillSolid = iota
	FillWireFrame
)

const (
	MatrixProjection = iota
	MatrixModelView
)

type GraphicsDevice interface {
	NewTexture2D(img image.Image) Texture2D
	NewShader(vertexShader, fragmentShader string) Shader
	NewVertexBufferVec2(vertices []Vec2) VertexBuffer
	NewVertexBufferVec3(vertices []Vec3) VertexBuffer
	NewIndexBuffer(indices []uint32) IndexBuffer

	SetFillMode(mode int)
	SetViewport(x, y, w, h int)

	SetMatrix(mtype int, m *Matrix4x4)
	SetTexture2D(texture Texture2D, unit uint)
	SetShader(shader Shader)
	SetTexCoords(buffer VertexBuffer, index uint)
	SetNormals(buffer VertexBuffer) 
	SetVertices(buffer VertexBuffer)
	DrawIndexed(buffer IndexBuffer)

	Clear()
}

type Texture2D interface {
	Release()
}

type VertexBuffer interface {
	Release()
}

type IndexBuffer interface {
	Release()
}

type Shader interface {
	GetUniformLocation(name string) uint
	SetVec3(location uint, v *Vec3)
	SetTexture(location uint, unit uint)
	Release()
}
