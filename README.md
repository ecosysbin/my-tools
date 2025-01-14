readme

昨天看文档，volcano通过动态链接库对插件进行扩展https://github.com/volcano-sh/volcano/blob/master/docs/design/custom-plugin.md。有个问题没弄懂：我们volcano系统的插件函数都是要注册到session功能点的map上的，这样调度流程走到对应功能点上从session获取到插件函数，然后进行处理。custom-plugin如何注册到session功能点的map上呢？往session哪个功能点注册呢？
