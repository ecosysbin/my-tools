# 编译动态链接库
gcc -shared -fPIC -o libmodule.so module.c

# 隐式动态链接库使用
## 引入动态链接库头文件
#include "module.h"
## 调用动态链接库头文件中的函数
如：test.c
## 编译并运行出错（需要将动态链接库加入到/usr/lib目录或者加入到LD_LIBRARY_PATH环境变量）
gcc test.c
./a.out

# 显式动态链接库使用
## 无需引入动态链接库头文件（如：test-perform.c）,需要引入dlfcn.h，显式的从动态链接库中解析出函数地址
## 编译并运行成功
gcc test-perform.c
./a.out

## 参考：https://www.jianshu.com/p/5594087fcdf7