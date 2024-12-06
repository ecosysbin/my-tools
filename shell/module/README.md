2. 代码解释
头文件引入：
#include <linux/module.h> ：这个头文件包含了内核模块相关的基本定义、宏以及函数声明等，是编写内核模块必不可少的头文件，像 module_init 和 module_exit 这些用于指定模块初始化和退出函数的宏就在这个头文件中定义。
#include <linux/kernel.h> ：引入内核的一些基础功能相关的定义，例如 printk 函数的相关定义就在这里面， printk 函数用于在内核空间打印消息，和用户空间使用的 printf 函数类似，但使用场景不同。
模块初始化函数 hello_module_init：
函数定义为 static int __init hello_module_init(void)，__init 是一个宏标记，用于告诉编译器这个函数在模块初始化完成后可以被丢弃，以节省内核内存空间。它返回一个 int 类型的值，返回 0 表示初始化成功，其他非零值表示初始化出现问题。
在函数内部，使用 printk 函数打印了一条消息 KERN_INFO "Hello, Linux kernel module is loaded!\n"，KERN_INFO 是一个日志级别宏，用于指定这条消息的重要性级别，这里表示信息级别消息，在内核日志中会按照相应的级别显示和记录。
模块清理函数 hello_module_exit：
定义为 static void __exit hello_module_exit(void)，同样 __exit 宏标记表示这个函数只有在模块被卸载时才会被调用，并且在模块没有被卸载的正常运行期间，编译器可以优化掉这个函数的代码，节省内存。
函数内部使用 printk 函数打印 KERN_INFO "Goodbye, Linux kernel module is unloaded!\n"，用于在模块卸载时输出一条提示消息。
模块入口定义：
module_init(hello_module_init); ：通过这个宏，将 hello_module_init 函数指定为模块加载时的初始化函数，当使用 insmod 或 modprobe 命令加载模块时，内核会自动调用这个函数。
module_exit(hello_module_exit); ：类似地，这个宏将 hello_module_exit 函数指定为模块卸载时的清理函数，当使用 rmmod 命令卸载模块时，内核会调用这个函数。
模块许可证及描述信息：
MODULE_LICENSE("GPL"); ：声明模块所遵循的开源许可证，这里选择了 GPL（通用公共许可证），这是 Linux 内核常用的开源协议声明方式，不同的许可证有不同的使用和分发规则。
MODULE_DESCRIPTION("A simple demo Linux kernel module"); ：简单描述模块的功能或性质，方便后续查看模块信息时了解其大致用途。
MODULE_AUTHOR("Your Name"); ：声明模块的作者信息，也是便于识别模块来源等情况。

3. 编译及测试步骤
编译环境准备：
要编译这个内核模块，需要安装 Linux 内核的编译工具链，在 Ubuntu 等常见的 Linux 发行版中，可以通过 sudo apt-get install build-essential linux-headers-$(uname -r) 命令来安装必要的编译工具（build-essential 包含了如 gcc、make 等常用编译工具）以及与当前内核版本匹配的头文件（linux-headers-$(uname -r) ）。

obj-m += hello_module.o 表示要编译生成的目标模块是 hello_module.o 对应的内核模块（最终生成 hello_module.ko 文件）。
make -C /lib/modules/$(uname -r)/build M=$(PWD) modules 这一行中，-C 选项指定了进入到内核源码目录（ /lib/modules/$(uname -r)/build 是指向当前内核编译相关目录的路径，$(uname -r) 会获取当前系统内核的版本号）去执行 make 命令，M=$(PWD) 表示模块的源代码所在的当前目录， modules 表示执行编译模块的操作。
clean 规则用于清理编译生成的中间文件和目标模块文件，同样是进入内核源码目录执行相应的清理操作。
编译模块：
将上述代码保存为 hello_module.c 文件，Makefile 文件也放在同一目录下，然后在该目录的终端中执行 make 命令，就会编译生成 hello_module.ko 文件，这个就是我们编写的内核模块文件。
测试模块：
使用 sudo insmod hello_module.ko 命令加载模块，然后可以通过 dmesg 命令查看内核日志（因为 printk 函数打印的消息会记录在内核日志中），应该能看到 Hello, Linux kernel module is loaded! 这条消息，表示模块加载成功。