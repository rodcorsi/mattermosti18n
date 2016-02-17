BINARY=./bin

all:
	go build -o ${BINARY}/i18n2po i18n2po.go
	go build -o ${BINARY}/po2i18n po2i18n.go

clean:
	if [ -d ${BINARY} ] ; then rm -fr ${BINARY} ; fi