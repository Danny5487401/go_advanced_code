#高级加密标准（Advanced Encryption Standard，AES）
又称Rijndael加密法，一个对称分组密码算法，是美国联邦政府采用的一种区块加密标准。这个标准用来替代原先的DES（Data Encryption Standard），已经被多方分析且广为全世界所使用。
AES中常见的有三种解决方案，分别为AES-128、AES-192和AES-256。

AES加密过程涉及到4种操作：字节替代（SubBytes）、行移位（ShiftRows）、列混淆（MixColumns）和轮密钥加（AddRoundKey）。
解密过程分别为对应的逆操作。由于每一步操作都是可逆的，按照相反的顺序进行解密即可恢复明文。
加解密中每轮的密钥分别由初始密钥扩展得到。算法中16字节的明文、密文和轮密钥都以一个4x4的矩阵表示。

如果采用真正的128位加密技术甚至256位加密技术，蛮力攻击要取得成功需要耗费相当长的时间。
1. 电码本模式（Electronic Codebook Book (ECB)）、

2. 密码分组链接模式（Cipher Block Chaining (CBC)）、

3. 计算器模式（Counter (CTR)）、

4. 密码反馈模式（Cipher FeedBack (CFB)）

5. 输出反馈模式（Output FeedBack (OFB)）

Note: 出于安全考虑，golang默认并不支持ECB模式。