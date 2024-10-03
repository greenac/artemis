.PHONY: serve
serve:
	go run cmd/serve/main.go

.PHONY: save
save:
	go run cmd/save/main.go

.PHONY: build
build:
	go build -o ~/go/bin/serveartemis cmd/serve/main.go

.PHONY: addactorstomovies
addactorstomovies:
	go run cmd/addactorstomovies/main.go

.PHONY: secondarypaths
secondarypaths:
	go run cmd/secondarypaths/main.go

.PHONY: swappaths
swappaths:
	go run cmd/swappaths/main.go

.PHONY: missingmovies
missingmovies:
	go run cmd/missingmovies/main.go

.PHONY: images
images:
	go run cmd/images/main.go
