include $(GOROOT)/src/Make.inc

all:
	cd src/pkg/g3 && gomake
	cd src/pkg/g3/geomipmapping && gomake

install:
	cd src/pkg/g3 && gomake install
	cd src/pkg/g3/geomipmapping && gomake install

clean:
	cd src/pkg/g3 && gomake clean
	cd src/pkg/g3/geomipmapping && gomake clean

example: install
	cd src/cmd/geomipmapping && gomake
