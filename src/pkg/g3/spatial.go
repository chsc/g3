package g3

// common functions/structs for generic spatial data structures like octrees, quadtrees ...

type NewSpatElementFunc func(depth uint, parent SpatElement) SpatElement

type TraverseFunc func(element SpatElement, leaf bool) bool

type FrustumTraversable interface {
	TraverseFrustum(frustum *Frustum, traverseFunc TraverseFunc)
}

type SpatElement interface {
	HasBoundingVolume
	FrustumTraversable
	GetParent() SpatElement
	GetChildren() []SpatElement
	GetData() interface{}
}

type SpatElementData struct {
	BVolume BoundingVolume
	Parent  SpatElement
	Data    interface{}
}

func (e *SpatElementData) GetBoundingVolume() BoundingVolume {
	return e.BVolume
}

func (e *SpatElementData) GetData() interface{} {
	return e.Data
}

func (e *SpatElementData) GetParent() SpatElement {
	return e.Parent
}

type SpatNode struct {
	SpatElementData
	Children []SpatElement
}

func (node *SpatNode) GetChildren() []SpatElement {
	return node.Children
}

func (node *SpatNode) TraverseFrustum(frustum *Frustum, traverseFunc TraverseFunc) {
	if frustum.ClipBoundingVolume(node.BVolume) {
		if traverseFunc(node, false) {
			for _, child := range node.Children {
				child.TraverseFrustum(frustum, traverseFunc)
			}
		}
	}
}

type SpatLeaf struct {
	SpatElementData
}

func (leaf *SpatLeaf) GetChildren() []SpatElement {
	return nil
}

func (leaf *SpatLeaf) TraverseFrustum(frustum *Frustum, traverseFunc TraverseFunc) {
	if frustum.ClipBoundingVolume(leaf.BVolume) {
		traverseFunc(leaf, true)
	}
}
