#第三次作业
http-deployment.yaml中包含全量的yaml，含service和ingress
##优雅启动
当前通过readinessProbe对健康检查接口进行可读性检查，确保其ready后及是可接受流量的状态，检查的成功、失败次数等设置采用默认配置即可，不做额外修改
```YAML
readinessProbe:
  httpGet:  
    path: /healthz 
    port: 8080
  periodSeconds: 10
```

##优雅终止
我当前的程序中没加程序逻辑上的优雅退出，而是通过spec.lifecycle.preStop去做一个简易的优雅退出。preStop的执行过程中pod已经被从endpoints列表中摘除，所以采取preStop进行10s的sleep确保当前的即时访问都已经返回，然后主动kill进行退出的方式完成优雅退出。
优雅终止可由程序逻辑、开发框架等实现。业务逻辑的实现孟老师的例子中已经实现，不做拷贝了。
```YAML
lifecycle:
  preStop:
    exec:
      command: ["/bin/sh"，"-c"，"sleep 10s && kill -QUIT 1"]
```

##资源需求和Qos保证
测试用的httpserver程序对资源的需求相对小，300m的cpu基本就能撑起数百的tps，基本没有io也没有复杂的业务逻辑，所以对磁盘和内存的要求几乎没有，网络上不做限速。
所以仅保证其有一定资源的限制，没有额外的需求。可随意驱逐，以Burstable的模式进行资源配额。
```yaml
resources:
  limits:
     memory: "128Mi"
     cpu: "300m"
```

##探活
作为web应用，选择最基本的探活手段，对容器所使用的8080端口进行tcpSocket检查，检查的成功、失败次数等设置采用默认配置即可，不做额外修改。
```YAML
livenessProbe:
  tcpSocket:
    port: 8080
  initialDelaySeconds: 10
```

##运维日志
引入glog去处理基本的日志需求，实现的httpsever没有过于复杂的业务，所以没有去在代码上做日志分级，在此处仅说明了思路。
对于进行日志分级，运维人员可通过-v参数来调整，查看更全面的日志。
同时可制定日志额外保存的路径，并配合volume进行持久化保存（此做法考虑的是集群中没有部署ELK进行日志持久化，并且docker有日志大小的限制，此处通过手动制定日志目录，后续出错时，可能会有收集跨度较大日志的需求）。
```YAML
command: ["/bin/sh","-c"]
args: ["/httpserver -log_dir=/tmp -alsologtostderr"]
```

##配置和代码分离
本次以configmap展示配置和代码分离的思路，将基本的配置信息配置到configmap中，然后go程序按照环境变量中指定的目录去读取解析配置文件，并应用到程序中，后续修改配置的操作在configmap中进行，实现配置和代码的分离。
```YAML
apiVersion: v1
kind: ConfigMap
metadata:
  name: httpserver-config
data:
  config.json: |-
    {
      "timeout": "3s",
      "httpport": "8080"
    }
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpserver-app
spec:
...
...
        env:
          - name: CONFIG_PATH
            value: "/config/myconfig/config.json"
        volumeMounts:
        - mountPath: /config/myconfig/config.json
          name: config
          subPath: config.json
          readOnly: true
      volumes:
        - name: config
          configMap:
            name: httpserver-config
...
```

##应用高可用
最基本的高可用，无状态应用采用多副本，通过service暴露提供服务，提供运行时的高可用。
```yaml
replicas: 4
```

对于deployment，更新的默认策略既是滚动更新，能保证基本的更新时或者故障时的高可用，因为没有额外需求，此处采用默认值，没有在spec中进行重新配置。
如果需要重新配置，需要变更以下的条目
```YAML
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 2  
```

##TLS和HTTPS
提供简单的配置说明：
1. 生成自签的证书并创建对应的secret
 ```shell
$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ${KEY_FILE} -out ${CERT_FILE} -subj "/CN=${HOST}/O=${HOST}"

kubectl create secret tls ${CERT_NAME} --key ${KEY_FILE} --cert ${CERT_FILE}
```
2. 在nginx-ingress-controller的启动参数中指定默认证书的secret
--default-ssl-certificate secretname

当然ingress上也可指定后端service所使用的证书，单本次作业中没有使用
