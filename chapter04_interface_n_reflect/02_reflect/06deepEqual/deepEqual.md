# reflect.DeepEqual函数：判断两个值是否一致

## 背景

对于array、slice、map、struct等类型，想要比较两个值是否相等，不能使用==，处理起来十分麻烦，在对效率没有太大要求的情况下，reflect包中的DeepEqual函数完美的解决了比较问题