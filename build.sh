echo "create new key.."
go run generate.go
cd generated_code
go build
mv gen ../build/tool
cd ..
#echo "hide last key, generate new.."
#go run generate.go
