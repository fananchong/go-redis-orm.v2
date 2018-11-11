set CURDIR=%~dp0
set BASEDIR=%CURDIR%;D:\temp
set GOPATH=%BASEDIR%
set GOBIN=%CURDIR%
go install ./...
redis2go.exe --input_dir=%CURDIR%\..\..\example\redis_def --output_dir=%CURDIR%\..\..\example\ --package=main