#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void check_malloc();
void check_calloc();
void check_memcpy();

int main(void)
{
    printf("Hello World!\n");
    // check_malloc();
    // check_calloc();
    check_memcpy();
    return 0;
}

void check_malloc()
{
    int *list = malloc(100 * sizeof(int));
    for (int i = 0; i < 100; i++)
    {
        printf("index: %d, value:%d \n", i, list[i]);
    }
    free(list);
}

void check_calloc()
{
    // 申请内存，并初始化为0
    int *list = calloc(100, sizeof(int));
    // 没有角标越界，需要代码自己控制，不然拿到的就是错误的值，可能把其他内存也给覆盖了
    // 申请了对应的长度内存，就用对应长度的内存，c内存管理器保证其他程序申请内存时不会覆盖已经申请的内存。
    for (int i = 0; i < 106; i++)
    {
        printf("index: %d, value: %d \n", i, list[i]);
    }
    for (int i = 0; i < 100; i++)
    {
        list[i] = i;
    }
    // 在原来的基础上增加申请内存，增加申请的内存不为0，如何确认是增加的内存还是盗用别人的内存?
    int *new_list = realloc(list, 100 * sizeof(int));

    for (int i = 0; i < 200; i++)
    {
        printf("new_list index: %d, value: %d \n", i, new_list[i]);
    }
    printf("new_list size: %d \n", sizeof(new_list));
    free(list);
}

void check_memcpy()
{
    char *src = "hello world";
    char *dest = malloc(11);
    // 将一块内存复制到另一块内存，不会修改原来内存的内容，只是将内容复制过去
    memcpy(dest, src, 5);
    printf("dest: %s \n", dest);
}