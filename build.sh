# Get all library
go mod tidy

# Set environment variables for Windows 64-bit
export GOOS=windows
export GOARCH=amd64

# Build the Windows executable
go build -o  ./dist/rest-api-gorilla.exe ./

export GOOS=linux
export GOARCH=amd64
go build -o  ./dist/rest-api-gorilla-linux ./

export GOOS=darwin
export GOARCH=amd64
go build -o  ./dist/rest-api-gorilla-mac ./



