# 策略模式（Strategy Pattern）,又叫作政策模式（Policy Pattern）
![](process.png)

> Define a family of algorithms, encapsulate each one, and make them interchangeable. Strategy lets the algorithm vary independently from clients that use it.

翻译成中文就是：定义一族算法类，将每个算法分别封装起来，让它们可以互相替换。策略模式可以使算法的变化独立于使用它们的客户端（这里的客户端代指使用算法的代码）。


## 场景

在我们写的程序中，大多有if else的条件语句基本上都适合Strategy 模式，但是if else 条件的情况是不变的，则不适合此模式，例如一周7天。