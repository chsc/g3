include $(GOROOT)/src/Make.inc

TARG=g3
GOFILES=math.go bbox.go frustum.go plane.go vector.go matrix.go utils.go \
	fileutils.go \
	spatial.go \
	graphics.go ogl_graphics.go \
	engine.go sdl_engine.go 

include $(GOROOT)/src/Make.pkg

format:
	gofmt -w *.go
