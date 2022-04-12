# 模版模式
定义一个算法骨架，将一些步骤延迟到子类进行。模板模式使得子类可以不改变一个算法的结构，即可重新定义该算法的某些特定步骤

## 优点
- 封装不变的部分，扩展可变部分
- 提取公共代码，便于维护
- 行为由父类控制，子类实现

## 缺点
每一个不同的实例，都需要一个子类来实现，导致类的个数增加。

## 案例
1. 做饭，打开煤气，开火，（做饭）， 关火，关闭煤气。除了做饭其他步骤都是相同的，抽到抽象类中实现
2. 让我们来考虑一个一次性密码功能 （OTP One Time Password） 的例子。

一次性密码功能 :又称“一次性口令”，是指只能使用一次的密码。 一次性密码是根据专门算法、每隔60秒生成一个不可预测的随机数字组合，iKEY一次性密码已在金融、电信、网游等领域被广泛应用，有效地保护了用户的安全。

将 OTP 传递给用户的方式多种多样 （短信、 邮件等）。 但无论是短信还是邮件， 整个 OTP 流程都是相同的：

1. 生成随机的 n 位数字。
2. 在缓存中保存这组数字以便进行后续验证。
3. 准备内容。
4. 发送通知。
5. 发布。

后续引入的任何新 OTP 类型都很有可能需要进行相同的上述步骤。

因此， 我们会有这样的一个场景， 其中某个特定操作的步骤是相同的， 但实现方式却可能有所不同。 这正是适合考虑使用模板方法模式的情况。

首先， 我们定义一个由固定数量的方法组成的基础模板算法。 这就是我们的模板方法。 然后我们将实现每一个步骤方法， 但不会改变模板方法。

