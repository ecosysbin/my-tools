#include <linux/module.h> // 引入内核模块相关头文件
#include <linux/kernel.h> // 引入内核相关的基本头文件

// 模块初始化函数，模块加载时执行
static int __init hello_module_init(void)
{
    printk(KERN_INFO "Hello, Linux kernel module is loaded!\n");
    return 0;
}

// 模块清理函数，模块卸载时执行
static void __exit hello_module_exit(void)
{
    printk(KERN_INFO "Goodbye, Linux kernel module is unloaded!\n");
}

// 定义模块初始化和清理函数的入口
module_init(hello_module_init);
module_exit(hello_module_exit);

// 模块许可证声明，开源协议相关，这里选择GPL
MODULE_LICENSE("GPL");

// 模块的简单描述信息
MODULE_DESCRIPTION("A simple demo Linux kernel module");

// 模块作者信息
MODULE_AUTHOR("Your Name");