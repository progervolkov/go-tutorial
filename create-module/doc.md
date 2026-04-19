go mod edit -replace example.com/greetings=../greetings

go tidy

$env:Path += ";C:\Users\mycod\go\bin"
go install