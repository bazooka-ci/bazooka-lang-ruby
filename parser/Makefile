default: image

image: gox
	docker build -t bazooka/parser-ruby .

gox:
	gox -osarch="linux/amd64" -output="main"
