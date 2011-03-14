package g3

type Frustum struct {
	Left, Right Plane
	Top, Bottom Plane
	Near, Far   Plane
}

func MakeFrustumFromMatrix(m *Matrix4x4) *Frustum {
	left := Plane{Vec3{m.M41 + m.M11, m.M42 + m.M12, m.M43 + m.M13}, m.M44 + m.M14}
	left.Normalize()

	right := Plane{Vec3{m.M41 - m.M11, m.M42 - m.M12, m.M43 - m.M13}, m.M44 - m.M14}
	right.Normalize()

	top := Plane{Vec3{m.M41 - m.M21, m.M42 - m.M22, m.M43 - m.M23}, m.M44 - m.M24}
	top.Normalize()

	bottom := Plane{Vec3{m.M41 + m.M21, m.M42 + m.M22, m.M43 + m.M23}, m.M44 + m.M24}
	bottom.Normalize()

	near := Plane{Vec3{m.M41 + m.M31, m.M42 + m.M32, m.M43 + m.M33}, m.M44 + m.M34}
	near.Normalize()

	far := Plane{Vec3{m.M41 - m.M31, m.M42 - m.M32, m.M43 - m.M33}, m.M44 - m.M34}
	far.Normalize()
/*
	left := Plane{Vec3{m.M14 + m.M11, m.M24 + m.M21, m.M34 + m.M31}, m.M44 + m.M41}
	left.Normalize()

	right := Plane{Vec3{m.M14 - m.M11, m.M24 - m.M11, m.M34 - m.M31}, m.M44 - m.M41}
	right.Normalize()

	top := Plane{Vec3{m.M14 - m.M12, m.M24 - m.M22, m.M34 - m.M32}, m.M44 - m.M42}
	top.Normalize()

	bottom := Plane{Vec3{m.M14 + m.M12, m.M24 + m.M22, m.M34 + m.M32}, m.M44 + m.M42}
	bottom.Normalize()

	near := Plane{Vec3{m.M14 + m.M31, m.M24 + m.M23, m.M34 + m.M33}, m.M44 + m.M43}
	near.Normalize()

	far := Plane{Vec3{m.M14 - m.M13, m.M24 - m.M23, m.M34 - m.M33}, m.M44 - m.M43}
	far.Normalize()
*/
	return &Frustum{left, right, top, bottom, near, far}
}

func (frustum *Frustum) ClipBoundingVolume(bvol BoundingVolume) bool {
	if bvol.ClassifyPlane(&frustum.Left) == 0 {
		return false
	}
	if bvol.ClassifyPlane(&frustum.Right) == 0 {
		return false
	}
	if bvol.ClassifyPlane(&frustum.Top) == 0 {
		return false
	}
	if bvol.ClassifyPlane(&frustum.Bottom) == 0 {
		return false
	}
	if bvol.ClassifyPlane(&frustum.Near) == 0 {
		return false
	}
	if bvol.ClassifyPlane(&frustum.Far) == 0 {
		return false
	}
	return true
}

func (frustum *Frustum) ClipPoint(v *Vec3) bool {
	if frustum.Left.DistanceToPoint(v) < 0.0 {
		return false
	}
	if frustum.Right.DistanceToPoint(v) < 0.0 {
		return false
	}
	if frustum.Top.DistanceToPoint(v) < 0.0 {
		return false
	}
	if frustum.Bottom.DistanceToPoint(v) < 0.0 {
		return false
	}
	if frustum.Near.DistanceToPoint(v) < 0.0 {
		return false
	}
	if frustum.Far.DistanceToPoint(v) < 0.0 {
		return false
	}
	return true
}

// TODO: ClipSphere

