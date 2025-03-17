package services

//go:generate powershell -Command "Remove-Item -Recurse -Force mocks -ErrorAction Ignore; New-Item -ItemType Directory -Path mocks"

///go:generate bash -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i UserService -o ./mocks/ -s "_minimock.go"
