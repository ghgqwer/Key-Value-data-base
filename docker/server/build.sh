docker build --progress=plain -t server:0.0.1 -f ./docker/server/Dockerfile.server .
# . - путь к необходимым файлам в имадже (удаленный)
#GOOS=linux GOARCH=amd64 
#--platform linux/amd64 