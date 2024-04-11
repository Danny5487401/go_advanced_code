<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go-cmp](#go-cmp)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# go-cmp

reflect.DeepEqual 的替代品

特点
- reflect.DeepEqual不够灵活，无法提供选项实现我们想要的行为，例如允许浮点数误差，对test 不友好 
- 其他类型可以通过 equal 方法扩展
- 不会比较未导出字段（即字段名首字母小写的字段）。遇到未导出字段，cmp.Equal()直接panic


## 参考资料

1 [每日一库 go-cmp](https://darjun.github.io/2020/03/20/godailylib/go-cmp/)