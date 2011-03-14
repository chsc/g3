package g3

import (
	"fmt"
)

// Represents a 4x4 Matrix
type Matrix4x4 struct {
	M11, M12, M13, M14 float32
	M21, M22, M23, M24 float32
	M31, M32, M33, M34 float32
	M41, M42, M43, M44 float32
}

func MakeMatrixFromSlice(m []float32) Matrix4x4 {
	return Matrix4x4{
		m[0], m[1], m[2], m[3],
		m[4], m[5], m[6], m[7],
		m[8], m[9], m[10], m[11],
		m[12], m[13], m[14], m[15]}
}

func MakeIdentityMatrix() Matrix4x4 {
	return Matrix4x4{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0}
}

func MakeScaleMatrix(x, y, z float32) Matrix4x4 {
	return Matrix4x4{
		x, 0.0, 0.0, 0.0,
		0.0, y, 0.0, 0.0,
		0.0, 0.0, z, 0.0,
		0.0, 0.0, 0.0, 1.0}
}

func MakeTranslationMatrix(x, y, z float32) Matrix4x4 {
	return Matrix4x4{
		1.0, 0.0, 0.0, x,
		0.0, 1.0, 0.0, y,
		0.0, 0.0, 1.0, z,
		0.0, 0.0, 0.0, 1.0}
}

func MakeXRotationMatrix(theta float32) Matrix4x4 {
	cosTheta := Cos(theta)
	sinTheta := Sin(theta)
	return Matrix4x4{
		1.0, 0.0, 0.0, 0.0,
		0.0, cosTheta, -sinTheta, 0.0,
		0.0, sinTheta, cosTheta, 0.0,
		0.0, 0.0, 0.0, 1.0}
}

func MakeYRotationMatrix(theta float32) Matrix4x4 {
	cosTheta := Cos(theta)
	sinTheta := Sin(theta)
	return Matrix4x4{
		cosTheta, 0.0, sinTheta, 0.0,
		0.0, 1.0, 0.0, 0.0,
		-sinTheta, 0.0, cosTheta, 0.0,
		0.0, 0.0, 0.0, 1.0}
}

func MakeZRotationMatrix(theta float32) Matrix4x4 {
	cosTheta := Cos(theta)
	sinTheta := Sin(theta)
	return Matrix4x4{
		cosTheta, -sinTheta, 0.0, 0.0,
		sinTheta, cosTheta, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0}
}

func MakePerspectiveMatrix(fovy, aspect, zNear, zFar float32) Matrix4x4 {
	f := 1.0 / Tan(fovy/2.0)
	a := 1.0 / (zNear - zFar)
	return Matrix4x4{
		f / aspect, 0.0, 0.0, 0.0,
		0.0, f, 0.0, 0.0,
		0.0, 0.0, (zFar + zNear) * a, 2.0 * zFar * zNear * a,
		0.0, 0.0, -1.0, 0.0}
}

func MakeLookAtMatrix(eye, center, up *Vec3) Matrix4x4 {
	f := center.Sub(*eye).Normalized()
	u := up.Normalized()
	s := f.Cross(u)
	u = s.Cross(f)
	t := MakeTranslationMatrix(-eye.X, -eye.Y, -eye.Z)
	return Matrix4x4{
		s.X, s.Y, s.Z, 0.0,
		u.X, u.Y, u.Z, 0.0,
		-f.X, -f.Y, -f.Z, 0.0,
		0.0, 0.0, 0.0, 1.0}.Multiply(&t)
}

func (m Matrix4x4) Transposed() Matrix4x4 {
	return Matrix4x4{
		m.M11, m.M21, m.M31, m.M41,
		m.M12, m.M22, m.M32, m.M42,
		m.M13, m.M23, m.M33, m.M43,
		m.M14, m.M24, m.M34, m.M44}
}

func (m1 Matrix4x4) Multiply(m2 *Matrix4x4) Matrix4x4 {
	return Matrix4x4{
		m1.M11*m2.M11 + m1.M12*m2.M21 + m1.M13*m2.M31 + m1.M14*m2.M41,
		m1.M11*m2.M12 + m1.M12*m2.M22 + m1.M13*m2.M32 + m1.M14*m2.M42,
		m1.M11*m2.M13 + m1.M12*m2.M23 + m1.M13*m2.M33 + m1.M14*m2.M43,
		m1.M11*m2.M14 + m1.M12*m2.M24 + m1.M13*m2.M34 + m1.M14*m2.M44,

		m1.M21*m2.M11 + m1.M22*m2.M21 + m1.M23*m2.M31 + m1.M24*m2.M41,
		m1.M21*m2.M12 + m1.M22*m2.M22 + m1.M23*m2.M32 + m1.M24*m2.M42,
		m1.M21*m2.M13 + m1.M22*m2.M23 + m1.M23*m2.M33 + m1.M24*m2.M43,
		m1.M21*m2.M14 + m1.M22*m2.M24 + m1.M23*m2.M34 + m1.M24*m2.M44,

		m1.M31*m2.M11 + m1.M32*m2.M21 + m1.M33*m2.M31 + m1.M34*m2.M41,
		m1.M31*m2.M12 + m1.M32*m2.M22 + m1.M33*m2.M32 + m1.M34*m2.M42,
		m1.M31*m2.M13 + m1.M32*m2.M23 + m1.M33*m2.M33 + m1.M34*m2.M43,
		m1.M31*m2.M14 + m1.M32*m2.M24 + m1.M33*m2.M34 + m1.M34*m2.M44,

		m1.M41*m2.M11 + m1.M42*m2.M21 + m1.M43*m2.M31 + m1.M44*m2.M41,
		m1.M41*m2.M12 + m1.M42*m2.M22 + m1.M43*m2.M32 + m1.M44*m2.M42,
		m1.M41*m2.M13 + m1.M42*m2.M23 + m1.M43*m2.M33 + m1.M44*m2.M43,
		m1.M41*m2.M14 + m1.M42*m2.M24 + m1.M43*m2.M34 + m1.M44*m2.M44}
}

func (m Matrix4x4) Transform(v Vec3) Vec3 {
	return Vec3 {
		m.M11*v.X+m.M12*v.Y+m.M13*v.Z+m.M14,
		m.M21*v.X+m.M22*v.Y+m.M23*v.Z+m.M24,
		m.M31*v.X+m.M32*v.Y+m.M33*v.Z+m.M34}
}

func (m *Matrix4x4) String() string {
	return fmt.Sprintf(
		"/%f %f %f %f\\\n|%f %f %f %f|\n|%f %f %f %f|\n\\%f %f %f %f/",
		m.M11, m.M12, m.M13, m.M14,
		m.M21, m.M22, m.M23, m.M24,
		m.M31, m.M32, m.M33, m.M34,
		m.M41, m.M42, m.M43, m.M44)
}
