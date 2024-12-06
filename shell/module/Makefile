obj-m += hello_module.o

all:
	# 注意使用tab键缩进, 而不是空格, 否则make命令会报错. $(uname -r)获取不到时，执行命令生成版本写进去。
    make -C /lib/modules/$(uname -r)/build M=$(PWD) modules

clean:
    make -C /lib/modules/$(uname -r)/build M=$(PWD) clean