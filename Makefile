default: stampede

ego.go: *.ego
	ego -package=main

stampede: ego.go *.go
	go build
