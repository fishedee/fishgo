.PHONY:test_macro test_reflect test_all
test_all:test_reflect test_macro 
	echo "test finish"
test_macro:
	go install github.com/fishedee/language/querygen
	rm -rf testdata/testdata_querygen.go
	querygen -r github.com/fishedee/language/testdata
	go test . -v -args macro
test_reflect:
	rm -rf testdata/testdata_querygen.go
	go test . -v -args reflect