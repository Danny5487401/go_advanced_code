<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Trace 可视化分析goroutine的调度](#trace-%E5%8F%AF%E8%A7%86%E5%8C%96%E5%88%86%E6%9E%90goroutine%E7%9A%84%E8%B0%83%E5%BA%A6)
  - [运行](#%E8%BF%90%E8%A1%8C)
    - [Goroutine analysis](#goroutine-analysis)
    - [view trace](#view-trace)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Trace 可视化分析goroutine的调度

## 运行
![](.trace_images/trace_main_menu.png)
第一行View trace（可视化整个程序的调度流程）和第二行Gorutine analysis。

### Goroutine analysis
![](.trace_images/goroutine_analysis.png)
- Goroutines:
- main.main.func1 N=3     
- runtime.main N=1
- runtime/trace.Start.func1 N=1
- N=4

解析： 程序一共有5个goroutine，分别是三个for循环里启动的匿名go func()、一个trace.Start.func1和runtime.main。

### view trace 
![](.trace_images/view_trace.png)
当主goroutine中的for循环逻辑已经走完并阻塞于wg.Wait（）一段时间后，go func的goroutine才启动准备（准备资源，挂载M线程等）完毕。
那么，此时三个goroutine中获取的url都是指向的最后一次for循环的url。

