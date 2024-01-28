<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [regexp 正则表达式 Regular Expression](#regexp-%E6%AD%A3%E5%88%99%E8%A1%A8%E8%BE%BE%E5%BC%8F-regular-expression)
  - [基本概念](#%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5)
  - [分类](#%E5%88%86%E7%B1%BB)
    - [BRE 和 ERE 的区别](#bre-%E5%92%8C-ere-%E7%9A%84%E5%8C%BA%E5%88%AB)
  - [Linux 中常用文本工具与正则表达式的关系](#linux-%E4%B8%AD%E5%B8%B8%E7%94%A8%E6%96%87%E6%9C%AC%E5%B7%A5%E5%85%B7%E4%B8%8E%E6%AD%A3%E5%88%99%E8%A1%A8%E8%BE%BE%E5%BC%8F%E7%9A%84%E5%85%B3%E7%B3%BB)
    - [1 grep , egrep 正则表达式特点：](#1-grep--egrep-%E6%AD%A3%E5%88%99%E8%A1%A8%E8%BE%BE%E5%BC%8F%E7%89%B9%E7%82%B9)
    - [2 sed 正则表达式特点](#2-sed-%E6%AD%A3%E5%88%99%E8%A1%A8%E8%BE%BE%E5%BC%8F%E7%89%B9%E7%82%B9)
    - [3 Awk（gawk）正则表达式特点](#3-awkgawk%E6%AD%A3%E5%88%99%E8%A1%A8%E8%BE%BE%E5%BC%8F%E7%89%B9%E7%82%B9)
  - [golang regexp](#golang-regexp)
    - [通过正则判断是否匹配](#%E9%80%9A%E8%BF%87%E6%AD%A3%E5%88%99%E5%88%A4%E6%96%AD%E6%98%AF%E5%90%A6%E5%8C%B9%E9%85%8D)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# regexp 正则表达式 Regular Expression

在计算机科学中，是指一个用来描述或者匹配一系列符合某个句法规则的字符串的单个字符串。在很多文本编辑器或其他工具里，正则表达式通常被用来检索和/或替换那些符合某个模式的文本内容。
许多程序设计语言都支持利用正则表达式进行字符串操作。

正则表达式这个概念最初是由Unix中的工具软件（例如sed和grep）普及开的。正则表达式通常缩写成“regex”，单数有regexp、regex，复数有regexps、regexes、regexen。这些是正则表达式的定义。
由于起源于unix系统，因此很多语法规则一样的。但是随着逐渐发展，后来扩展出以下几个类型。


其实字符串处理我们可以使用strings包来进行搜索(Contains、Index)、替换(Replace)和解析(Split、Join)等操作，
但是这些都是简单的字符串操作，他们的搜索都是大小写敏感，而且固定的字符串，如果我们需要匹配可变的那种就没办法实现了，当然如果strings包能解决你的问题，那么就尽量使用它来解决。
因为他们足够简单、而且性能和可读性都会比正则好
## 基本概念

- leftmost-first： 在匹配文本时，该正则表达式会尽可能早的开始匹配，并且在匹配过程中选择回溯搜索到的第一个匹配结果
```shell
func Compile(expr string) (*Regexp, error)
```

- leftmost longest（最左最长匹配 ）:默认正则为贪婪匹配Lazy Repetition。
```shell
func CompilePOSIX(expr string) (*Regexp, error) 
```




## 分类

POSIX 标准中定义了两种正则标准：Basic Regular Expressions (BREs)和 Extended Regular Expressions (EREs)。

1、基本的正则表达式（Basic Regular Expression 又叫 Basic RegEx  简称 BREs）

2、扩展的正则表达式（Extended Regular Expression 又叫 Extended RegEx 简称 EREs）

3、Perl 的正则表达式（Perl Regular Expression 又叫 Perl RegEx 简称 PREs）

### BRE 和 ERE 的区别
- 在基本的正则表达式中， {, |, (, ) 仅是其代表的字面字符，要进行转义（前面加反斜线 \ ）才是正则的元字符；在基本的正则表达式中不支持元字符 ?, + 。
- 在扩展的正则表达式中， ?, +, {, |, (, ) 都是正则元字符，这样写书正则表达式更方便


## Linux 中常用文本工具与正则表达式的关系 

### 1 grep , egrep 正则表达式特点：

1）grep 支持：BREs、EREs、PREs 正则表达式

grep 指令后不跟任何参数，则表示要使用 ”BREs“

grep 指令后跟 ”-E" 参数，则表示要使用 “EREs“

grep 指令后跟 “-P" 参数，则表示要使用 “PREs"



2）egrep 支持：EREs、PREs 正则表达式

egrep 指令后不跟任何参数，则表示要使用 “EREs”

egrep 指令后跟 “-P" 参数，则表示要使用 “PREs"


### 2 sed 正则表达式特点
sed 文本工具支持：BREs、EREs

- sed 指令默认是使用"BREs"
- sed 命令参数 “-r ” ，则表示要使用“EREs"

### 3 Awk（gawk）正则表达式特点
Awk 文本工具支持：EREs

- awk 指令默认是使用 “EREs"


## golang regexp 

采用 [re2](https://github.com/google/re2/wiki/Syntax) 语法（除了\c、\C），和Perl、Python等语言的正则基本一致。
```shell
➜ go doc regexp/syntax
```

### 通过正则判断是否匹配
```shell
func Match(pattern string, b []byte) (matched bool, error error)
func MatchReader(pattern string, r io.RuneReader) (matched bool, error error)
func MatchString(pattern string, s string) (matched bool, error error)
```


## 参考
- [正则使用](https://cloud.tencent.com/developer/article/2173761)
