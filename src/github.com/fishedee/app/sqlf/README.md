# 设计

比起实现代码，我觉得更困难的是，如何设计sqlf的接口？

sqlf的目标是

* 容易上学代码可读
* 开发效率高

设计要点：

* 像xorm和gorm其实上手不太容易，有学习成本，而且做了很多隐式工作，刚开始接手代码可能会崩溃，所以sqlf还是沿用原有query+args的方式来操作数据库。
* 我们需要支持一个struct同时能add，query和mod同一个表的数据，这样开发效率才能高。但是，struct的auto_increment的key在insert和update时是不能传入的，仅在select时才传入。另外，createTime和modifyTime不能依赖数据库的实现，因为，1.不是所有数据库都支持on update CURRENT_TIMESTAMP，2.mysql的on update CURRENT_TIMESTAMP仅在与原数据不同时才能更新timestamp，但是我们的要求时调用update就更新timestamp。因此，sqlf要处理好对这些tag的字段。