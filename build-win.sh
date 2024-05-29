# Get all library
go mod tidy

# Set environment variables for Windows 64-bit
# shellcheck disable=SC2121
set GOOS=windows
# shellcheck disable=SC2121
set GOARCH=amd64

# Build the Windows executable
go build -o  ./dist/rest-api-gorilla.exe ./

# shellcheck disable=SC2121
set GOOS=linux
# shellcheck disable=SC2121
set GOARCH=amd64
go build -o  ./dist/rest-api-gorilla-linux ./

# shellcheck disable=SC2121
set GOOS=darwin
# shellcheck disable=SC2121
set GOARCH=amd64
go build -o  ./dist/rest-api-gorilla-mac ./



