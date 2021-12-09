# fcc-mock-restaurant-backend
A mock restaurant backend with user authenticated ordering created by a FreeCodeCamp Dallas community team.



## Windows setup notes:
####  Installing gorm.io/driver/sqlite
* To overcome the issue with 'cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%' I found the following [Stack Overflow](https://stackoverflow.com/questions/43580131/exec-gcc-executable-file-not-found-in-path-when-trying-go-build) page.
* Basically you just need to make sure you install mingw-64 and ensure it's setup for x86_64 
* Then include the <installation_path>\bin in your environmental variables

#### WSL notes:
* When using WSL through goland I had to export `CGO_ENABLED=1` as well as compiling on remote