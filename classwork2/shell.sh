#��������
docker build -t httpserver:v1.0 .

#���;���
docker tag httpserver:v1.0 lijiaxing10086/httpserver:v1.0
docker push lijiaxing10086/httpserver:v1.0

#���о���
docker run -it --name httpserver -p 8080:8080 -d httpserver:v1.0 

#ͨ��nsenter���������鿴ip
docker inspect httpserver | grep Pid  #��ȡpid
nsenter -t PID -n ip a                #����pid���������鿴ip
