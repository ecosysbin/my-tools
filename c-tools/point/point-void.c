#include <stdio.h>

void print_int(void **ptr)
{
    // 取ptr内容，并将其转换成int指针类型，赋值给int_ptr
    int *int_ptr = (int *)*ptr;
    // 取int_ptr的内容，是一个int类型
    printf("The integer value is: %d\n", *int_ptr);
}

int main()
{
    // num 是一个int类型变量
    int num = 42;
    // 取num的地址，并将地址赋值给int指针类型的变量num_ptr
    int *num_ptr = &num;
    // 将num_ptr的地址赋值给void指针类型的变量void_ptr, 因为指针不是像int,char或者struct这种基本类型，则使用void*来表示
    void *void_ptr = num_ptr;
    // 将void_ptr的地址赋值给void指针指针类型的变量void_ptr_ptr
    void **void_ptr_ptr = &void_ptr;

    print_int(void_ptr_ptr);

    return 0;
}