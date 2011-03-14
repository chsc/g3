package g3

type Vec2 struct {
	X, Y float32
}

type Vec3 struct {
	X, Y, Z float32
}

type Vec4 struct {
	X, Y, Z, W float32
}

func (self Vec3) Add(a Vec3) Vec3 {
	return Vec3{X: self.X + a.X, Y: self.Y + a.Y, Z: self.Z + a.Z}
}

func (self Vec3) Sub(a Vec3) Vec3 {
	return Vec3{X: self.X - a.X, Y: self.Y - a.Y, Z: self.Z - a.Z}
}

func (self Vec3) Mul(a float32) Vec3 {
	return Vec3{X: self.X * a, Y: self.Y * a, Z: self.Z * a}
}

func (self Vec3) Dot(a Vec3) float32 {
	return self.X*a.X + self.Y*a.Y + self.Z*a.Z
}

func (self Vec3) Cross(a Vec3) Vec3 {
	return Vec3{X: self.Y*a.Z - self.Z*a.Y, Y: self.Z*a.X - self.X*a.Z, Z: self.X*a.Y - self.Y*a.X}
}

func (self Vec3) Inverted() Vec3 {
	return Vec3{X: -self.X, Y: -self.Y, Z: -self.Z}
}

func (self Vec3) Scaled(s float32) Vec3 {
	return Vec3{X: self.X*s, Y: self.Y*s, Z: self.Z*s}
}

func (self Vec3) Normalized() Vec3 {
	l := 1.0 / self.Length()
	return Vec3{X: self.X*l, Y: self.Y*l, Z: self.Z*l}
}

func (v *Vec3) Set(x, y, z float32) {
	v.X = x
	v.Y = y
	v.Z = z
}

func (self *Vec3) Accumulate(a Vec3) {
	self.X += a.X
	self.Y += a.Y
	self.Z += a.Z
}

func (self *Vec3) Scale(a float32) {
	self.X *= a
	self.Y *= a
	self.Z *= a
}

func (v *Vec3) Normalize() {
	l := 1.0 / v.Length()
	v.X *= l
	v.Y *= l
	v.Z *= l
}

func (v *Vec3) LengthSq() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vec3) Length() float32 {
	return Sqrt(v.LengthSq())
}

func (v *Vec3) TransformToWorld(x, y, z *Vec3) Vec3 {
	return Vec3{
		X: x.X*v.X + y.X*v.Y + z.X*v.Z,
		Y: x.Y*v.X + y.Y*v.Y + z.Y*v.Z,
		Z: x.Z*v.X + y.Z*v.Y + z.Z*v.Z}
}

func (v *Vec3) TransformToLocal(x, y, z *Vec3) Vec3 {
	return Vec3{
		X: x.X*v.X + x.Y*v.Y + x.Z*v.Z,
		Y: y.X*v.X + y.Y*v.Y + y.Z*v.Z,
		Z: z.X*v.X + z.Y*v.Y + z.Z*v.Z}
}
