rd /s/q release
md release
::go build -ldflags "-H windowsgui" -o chat.exe
go build -o chat.exe
COPY chat.exe release\
COPY favicon.ico release\favicon.ico
XCOPY asset\*.* release\asset\  /s /e
XCOPY views\*.* release\view\  /s /e
XCOPY config\*.* release\config\  /s /e
XCOPY index.html release\ /s /e
