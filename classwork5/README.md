#第5次作业

1. 将应用通过gateway发布

   ```shell
   #设置命名空间自动注入边车
   kubectl create ns istio-test
   kubectl label ns istio-test istio-injection=enabled
   ```

   

   ```YAML
   apiVersion: networking.istio.io/v1beta1
   kind: Gateway
   metadata:
     name: httpsserver-app
   spec:
     selector:
       istio: ingressgateway
     servers:
       - hosts:
           - httpsserver-app.cncamp.io
         port:
           name: https-default
           number: 443
           protocol: HTTPS
         tls:
           mode: SIMPLE
           credentialName: cncamp-credential
   ```

2. 七层路由规则，因本身功能不复杂，没有设计过多的条目，其他条目通过bookinfo实例进行了实验，如针对jason用户的故障注入，7s延迟响应

   ```YAML
   ...
      - match:
           - port: 443
         route:
           - destination:
               host: httpserver-app.istio-test.svc.cluster.local
               port:
                 number: 8080
   ...
   ```

   ```YAML
   apiVersion: networking.istio.io/v1beta1
   kind: VirtualService
   ...
   spec:
     hosts:
     - ratings
     http:
     - fault:
         delay:
           fixedDelay: 7s
           percentage:
             value: 100
       match:
       - headers:
           end-user:
             exact: jason
       route:
       - destination:
           host: ratings
           subset: v1
     - route:
       - destination:
           host: ratings
           subset: v1
   ```

3. 安全通过Gateway的tls进行处理

   ```yaml
         tls:
           mode: SIMPLE
           credentialName: cncamp-credential
   ```

   ```shell
   #证书生产
   openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt
   kubectl create -n istio-system secret tls cncamp-credential --key=cncamp.io.key --cert=cncamp.io.crt
   kubectl apply -f istio-specs.yaml -n istio-test
   ```

4.  tracing 的接入

   ```shell
   kubectl apply -f jaeger.yaml
   ```

   

![jaeger](https://github.com/lijiaxing10086/golang-test/blob/main/classwork5/images/jaeger.png)