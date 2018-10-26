TAR?=tar
VERSION=$(shell git describe --tags --long --abbrev=4 --dirty=-D)
DISTDIR=e212-$(VERSION)

.PHONY: e212 dist e212-$(VERSION) e212_cmd clean all
all: e212 e212_cmd

e212:
	go build -o $@ --ldflags "-X main.gVersion=$(VERSION)" web/main.go 
e212_cmd:
	go build -o $@ --ldflags "-X main.gVersion=$(VERSION)" cmd/main.go

e212-$(VERSION):
	mkdir -p $(DISTDIR)
	mkdir -p $(DISTDIR)/etc/
	cp e212 $(DISTDIR)/e212
	cp e212_cmd $(DISTDIR)/e212
	mkdir -p $(DISTDIR)/public/
	mkdir -p $(DISTDIR)/templates/
	cp -ap e212.service $(DISTDIR)/etc/
	cp -ap public/* $(DISTDIR)/public/
	cp -ap templates/* $(DISTDIR)/templates/
	${TAR} --owner=nobody --group=nobody -cvzf $(DISTDIR).tar.gz $(DISTDIR)
	rm -rf $(DISTDIR)

dist: e212 e212_cmd e212-$(VERSION)

clean:
	rm -f e212 e212_cmd	
	rm -rf *tar.gz
	rm -rf e212-*
