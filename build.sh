go run generate.go
cd generated_code
go build
mv generated_code ../build/tool
cd ..
go run generate.go
