dist-update:
	git submodule foreach 'git checkout master; git pull'
	go-bindata -pkg=main -o=dist.go ./ciste-web-content/dist/*
