echo "create new key.."
go run generate.go
cd generated_code
set GOOS=windows
set GOARCH=amd64
go build
mv gen.exe ../build/tool.exe
cd ..
#echo "hide last key, generate new.."
#go run generate.go
