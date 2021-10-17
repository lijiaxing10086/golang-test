#构建镜像
docker build -t httpserver:v1.0 .

#推送镜像
docker tag httpserver:v1.0 lijiaxing10086/httpserver:v1.0
docker push lijiaxing10086/httpserver:v1.0

#运行镜像
docker run -it --name httpserver -p 8080:8080 -d httpserver:v1.0 

#通过nsenter进入容器查看ip
docker inspect httpserver | grep Pid  #获取pid
nsenter -t PID -n ip a                #根据pid进入容器查看ip
