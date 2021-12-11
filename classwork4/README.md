#第4次作业

1. 为HTTPserver添加延时metric

   此处参照孟老师的提交进行增加,相关代码也一并提交到了该目录

2. 将HTTPserver部署到业务集群，并增加普罗米修斯的配置

   通过helm部署loki，然后使用其中的prometheus，Grafana，然后给httpserver的deployment增加annotation

   ```YAML
   ...
   spec:
   ...
         annotations:
           prometheus.io/scrape: "true"
           prometheus/port: "8080"
   ...
   ```

   

3. 在prometheus的界面中查询延时指标数据

![prometheus查询延时指标数据](https://github.com/lijiaxing10086/golang-test/blob/main/classwork4/images/prometheus.png)



![grafana查询延时指标数据](https://github.com/lijiaxing10086/golang-test/blob/main/classwork4/images/grafana.png)

4. 使用Grafana创建一个DashBoard来展示延时分配情况

![grafana-DashBoard](https://github.com/lijiaxing10086/golang-test/blob/main/classwork4/images/grafana-dashboard.png)
