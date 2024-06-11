## 将库文件编译成.c文件
gcc -c ./sub.c
gcc -c ./add.c


## 生成动态链接库文件
ar -crv libmymath.a sub.o add.o


## 编译main.c生成可执行文件
-L后面的.表示从当前路径开始检索库文件，-l后面的mymath表示链接的库文件名（去掉lib前缀和.a后缀）
gcc -o main main.c -L. –lmymath 