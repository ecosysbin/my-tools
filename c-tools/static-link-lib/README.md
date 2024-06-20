## 将库文件编译成.c文件
gcc -c ./sub.c
gcc -c ./add.c


## 生成静态态链接库文件
ar -crv libmymath.a sub.o add.o


## 编译main.c生成可执行文件
-L后面的.表示从当前路径开始检索库文件，-l后面的mymath表示链接的库文件名（去掉lib前缀和.a后缀）
gcc -o main main.c -L. –lmymath 

## 静态库.h文件中的_SUB_H_宏定义如下，即判断宏_SUB_H是否已定义，如果未定义则定义，否则不定义。增加这个逻辑判断，则可以减少编译次数，提供编译
#ifndef _SUB_H_

#define _SUB_H_

int sub(int a, int b);

#endif