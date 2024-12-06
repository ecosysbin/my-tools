# 1. 安装gcc
wget https://github.com/skeeto/w64devkit/releases/download/v2.0.0/w64devkit-x64-2.0.0.exe

## 1.1 安装
双击安装在指定目录

## 1.2 配置环境变量

# 2. 安装vscode c/c++开发插件
创建.c文件根据提示安装即可（须有外网，内网环境参考go开发环境插件包导入）

# 3. 编译
gcc -o main.exe main.c
# 4. 编译选项，编译过程中会输出详细信息
gcc --verbose ./main.c -static 
# 5. -g 参数是为了在生成的可执行文件中添加调试信息，方便后续调试
gcc -g main.c -o main  