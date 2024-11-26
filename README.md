# Firefly Chat Room

## 目录

- [简介](#简介)
- [环境要求](#环境要求)
- [安装](#安装)
- [运行](#运行)
- [功能](#功能)
- [使用说明](#使用说明)
- [开发指南](#开发指南)
- [贡献](#贡献)
- [许可证](#许可证)

## 简介

Firefly Chat Room 是一个简单的聊天室应用，支持多用户在线聊天、私聊和设置用户名等功能。项目分为服务器端和客户端两部分，使用 Go 语言编写。

## 环境要求

- Go 1.14 或更高版本
- 支持 TCP 协议的网络环境

## 安装

1. 克隆项目仓库：
   ```sh
   git clone
   cd firefly-chat-room
   ```

## 运行

- 启动服务器

1. 导航到服务器目录：

   ```sh
   cd server

   ```

2. 编译并运行服务器：
   ```sh
    go run server.go
   ```

- 启动客户端

1. 导航到客户端目录：

   ```sh
   cd client

   ```

2. 编译并运行客户端：
   ```sh
    go run client.go
   ```

## 功能

- 用户登录
- 广播消息
- 私聊消息
- 设置用户名

## 使用说明

1. 启动服务器后，客户端可以通过连接 localhost:8080 加入聊天室。

2. 客户端启动后，会显示帮助信息：

   1. set:your name
   2. all:your msg -- broadcast
   3. anyone:your msg -- private msg

3. 使用命令进行操作：

   - 设置用户名：set:yourname
   - 发送广播消息：all:your message
   - 发送私聊消息：anyone:your message
