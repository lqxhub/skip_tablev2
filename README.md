# skiptable

#### 介绍

**skiptable** 是一个使用go 模仿redis的zset 的跳表，实现了redis zset的大多数功能 需要 **go 版本>=1.18**

#### 软件架构
在`skip_list.go` 文件中实现了一个跳表，并实现了跳表的基本功能（增删改查）函数

在`sort_sort.go` 文件中实现了一个类似redis zset的有序集合,并实现了redis的绝大多数功能

#### 安装教程

* 通过go get 安装 `go get https://github.com/lqxhub/skip_tablev2`

#### 使用说明

没有任何外部依赖，直接导入代码即可使用

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request
