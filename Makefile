
bin_dir:
	@mkdir -p bin

clean:
	@rm -r bin

mysql:bin_dir
	@go build -o bin/mysql-import-client cmd/mysql_sample/main.go

all:bin_dir mysql
	@echo "all build OK"
