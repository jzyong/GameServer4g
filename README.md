Intro
====
&emsp;&emsp;[GameServer4j](https://github.com/jzyong/GameServer4j) for go version, distributed framework. Clients and gateways use
TCP custom protocols, Intranet message forwarding use GRPC forwarding, all stateless services can be horizontally
extended, and stateful services can be horizontally extended through partition, state binding and other rules. The basic
architecture of the project is shown below:



![项目架构图](https://raw.githubusercontent.com/jzyong/mmo-server/master/mmo-res/img/mmo%E6%9C%8D%E5%8A%A1%E5%99%A8.png)

Module
====

Navmesh pathfinding [server](https://github.com/jzyong/GameAI4j) [client](https://github.com/jzyong/NavMeshDemo)

Project                     |Description
--------------------------- |------------------------------              
game-api                    |login,charge logic
game-common                 |common logic,local tool etc
game-gate                   |network handle,message dispatcher
game-hall                   |game logic server
game-manager                |maintain background http services
game-message                |protobuf message
game-res                    |document,script etc
game-service                |game micro service
game-world                  |world logic



### TODO
* refactoring src
* mongoDB
* 多网关测试,消息分发
* 编写测试客户端（测试tcp通信，模拟登录流程） （转发大厅返回消息给client）
* add document


discuss
---------
* QQ交流群：143469012



