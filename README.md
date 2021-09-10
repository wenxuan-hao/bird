# 简介
用golang实现的并发下载器。

还没写完

# quick start
```
// 下载&编译
go get 
cd bird 
go install .

// 设置bird下载器的临时存储目录
export BIRD_PATH = '~/.bird' (默认地址)

// 下载，-c表示并发数量，默认为cpu核数  
bird download url -c n

// 恢复下载中断的任务
bird continue url

// 查看当前执行的下载任务
bird task
```

# 特性
- 下载进度条显示
- 支持并发下载
- 支持下载中断恢复
- 支持查看目前所有的下载任务