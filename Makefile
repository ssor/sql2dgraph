
bin_dir:
	@mkdir -p bin

clean:
	@rm -r bin
	@rm -r data/output

mysql:bin_dir
	@go build -o bin/mysql-client cmd/mysql_sample/main.go

resetdb:bin_dir
	@go build -o bin/reset cmd/resetdb/main.go

all:bin_dir mysql resetdb
	@echo "all build OK"
