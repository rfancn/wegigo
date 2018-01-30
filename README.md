# go-wegigo

[![Build Status](https://travis-ci.org/rfancn/wegigo.svg?branch=master)](https://travis-ci.org/rfancn/wegigo)

Wegigo is a distribute application develope framework. GIGO means gust in gust out, Let's get gust and output some gust too. :)
As we all know this is joke, Gust acturally means some kind of info, we receive some kind of info, process it and respond with another kind info, this is normally what we do in real world.
Wegigo is designed to develop some application to proceed in/out messages rapidly and efficiently.

As an exampke, Now wegigo intergrated with a Wechat Mediaplatform Package, I will separate to another project in the future.

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

                         |--- queue -- app...app
message-> Wegigo Server -|--- quque -- app...app
                         |--- queue -- app...app
#### Target

- scalable
- high performance
- high concurrency
- failover
- rapid development
- easy to deploy

#### Wxmp Module Workflow
- For synchronized reply message

								  |---> enabled appA's queue --> WxmpRequest -> appA process --|
	HttpRequest -> WxmpRequest -> |---> enabled appB's queue --> WxmpRequest -> appB process --|--> WxmpReply -> Wxmp's reply queue -> WxmpReply -> HttpResponse
								  |--- >enabled appC's queue --> WxmpRequest -> appC process --|

- For asynchronized reply message


								  |---> enabled appA's queue --> appA process --|
	HttpRequest -> WxmpRequest -> |---> enabled appB's queue --> appB process --|--> Response Success HttpResponse
								  |--- >enabled appC's queue --> appC process --|






