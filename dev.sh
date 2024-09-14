PORT=1323 nodemon --watch './**/*.go' --watch './public/views/**/*.html' --signal SIGTERM --exec 'go' run server.go
