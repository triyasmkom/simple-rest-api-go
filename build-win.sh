# Get all library
go mod tidy

# Set environment variables for Windows 64-bit
set GOOS=windows
set GOARCH=amd64

# Build the Windows executable
go build -o  ./dist/rest-api-gorilla.exe ./

set GOOS=linux
set GOARCH=amd64
go build -o  ./dist/rest-api-gorilla-linux ./

set GOOS=darwin
set GOARCH=amd64
go build -o  ./dist/rest-api-gorilla-mac ./



