# Wegigo

[![Build Status](https://travis-ci.org/rfancn/wegigo.svg?branch=master)](https://travis-ci.org/rfancn/wegigo)

Wegigo is a distributed application develope framework. GIGO means garbage in garbage out, Let's get garbage and output some garbage too. :)
As we all know this is a joke, Garbage acturally means some kind of info, we receive some kind of info, process it and respond with another kind info, this is normally what we do in real world.
Wegigo is designed to develop some application to proceed in/out messages rapidly and efficiently.

As an exampke, Now wegigo intergrated with a Wechat Mediaplatform Package, I will separate it to another project in the future.

### Concept

Package:
- codes: [src]/pkg/[package name]
- roles: it represents a server which can be configured to receive/distribute messages and manage plugin based applications

Application:
- codes: [src]/apps/[application name]
- roles: responsible for processing message and execute corresponding actions

Sdk:
- code: [src]/sdk
- roles: provide utility libraires to help develop package/application

#### Architecture

- Using docker container as the basic service unit
- Using kubernetes for resources orchestration
- Using RabbitMQ as the message broker to separate business logic and messages
- Using Etcd as cluster level config db
- Using monodb as persistence data storage(not implemented yet)
- Using redis/memcache to do cache stuff(not implemented yet)
```

input <-> Wegigo Framework <-> output
                 |
  sever dockerA... server dockerB
                 |
            message broker

```

#### Advantages

- scalable
- high performance
- high concurrency
- failover
- rapid development
- easy to deploy

#### Wxmp Server Workflow

- For synchronized reply message
```

               |---> enabled appA's queue --> appA process --|
HttpRequest -> |---> enabled appB's queue --> appB process --|--> WxmpReply -> Wxmp's reply queue -> HttpResponse
               |--- >enabled appC's queue --> appC process --|
```

- For asynchronized reply message
```
							  |---> enabled appA's queue --> appA process --|
HttpRequest -> WxmpRequest -> |---> enabled appB's queue --> appB process --|--> HttpResponse
							  |--- >enabled appC's queue --> appC process --|
```






