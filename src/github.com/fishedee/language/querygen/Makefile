.PHONY:gen
bench:gen
	go test . -bench=. -benchmem -v 
gen:
	go install github.com/fishedee/language/querygen
	rm -rf testdata/testdata_querygen.go
	querygen -r github.com/fishedee/language/querygen/testdata
test:gen
	go test . -v
