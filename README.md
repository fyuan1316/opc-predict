# opcdata-predict-server


```bash
make clean && make build-linux



# run websocket server
./opc-predict-server-linux start --wsport 9000 --brokers tbds-172-27-0-6:6668 --topic testfy --username kafka --passwd kafka@Tbds.com --sasl true
```

# 展示结果数据
## 1 在界面中查看
  
浏览器中打开 http://localhost:8080
(如果在浏览器中查看，请确保websockets.html 页面和 opc-predict-server-linux 在同一个路径下，并且如果服务启动时指定了wsport，websockets.html页面中的port 地址要一致。)

## 2 在命令行查看
可以安装wscat 在命令行查看结果数据
```bash
npm install -g wscat


wc -c wscat -c ws://172.27.0.5:9000/echo
(/echo 是websocket router的path)
```


## 结合数据生成程序一起使用
```bash
java -Djava.security.auth.login.config=/opt/tbds/kafka/config/kafka_client_for_ranger_jaas.conf -jar simpleclient-4.3.0-1075.jar --topic testfy --protocol SASL_PLAINTEXT

```