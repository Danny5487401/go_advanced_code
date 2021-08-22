package _6memory

/*
 类似于TCMalloc的结构
 使用span机制来减少碎片. 每个span至少为一个页(go中的一个page为8KB). 每一种span用于一个范围的内存分配需求. 比
如16-32byte使用分配32byte的span, 112-128使用分配128byte的span.
 一共有67个size范围, 8byte-32KB, 每个size有两种类型(scan和noscan, 表示分配的对象是否会包含指针)
 多层次Cache来减少分配的冲突. per-P无锁的mcache, 全局67*2个对应不同size的span的后备mcentral, 全局1个的mheap.
 mheap中以treap的结构维护空闲连续page. 归还内存到heap时, 连续地址会进行合并.
 stack分配也是多层次和多class的.
 对象由GC进行回收. sysmon会定时把空余的内存归还给操作系统
*/
