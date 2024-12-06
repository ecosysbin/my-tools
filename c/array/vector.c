#include <stdio.h>
#include <string.h>
#include <stdlib.h>

typedef struct List
{
    int *data;
    int size;
    int capacity;
} List_def;

List_def newList(int size)
{
    // malloc仅负责分配指定大小的内存空间，分配的内存内容是随机的，不会自动初始化为零。而calloc在分配内存后，会自动将分配的内存空间初始化为零，即将所有字节设置为0
    // 使用场景‌：由于malloc不初始化内存，适用于需要快速分配内存且不关心初始值的情况；而calloc适用于需要初始化内存为零的场景，例如在C语言中初始化数组时非常有用。
    // realloc 用于重新分配内存。它接受两个参数，即一个先前分配的指针和一个新的内存大小，然后尝试重新调整先前分配的内存块的大小。如果调整成功，它将返回一个指向重新分配内存的指针，否则返回一个空指针
    int *list = malloc(sizeof(int) * size);
    List_def list_def = {list, size, 0};
    return list_def;
}

void appendList(List_def *list, int data)
{
    if (list->size == 0)
    {
        list->size = 100;
        list->data = malloc(sizeof(int) * 100);
    }

    if (list->size == list->capacity)
    {
        int newSize = list->capacity * 2;
        int *newList = malloc(sizeof(int) * newSize);
        for (int i = 0; i < list->size; i++)
        {
            newList[i] = list->data[i];
        }
        free(list->data);
        list->data = newList;
        list->size = newSize;
    }
    list->data[list->capacity++] = data;
}

void printList(List_def *list)
{
    for (int i = 0; i < list->capacity; i++)
    {
        printf("%d\n", list->data[i]);
    }
}

// 常用的list map string实现就能写程序
int main()
{
    List_def list = newList(5);
    List_def *pList = &list;
    appendList(pList, 1);
    printList(pList);
    appendList(pList, 2);
    printList(pList);
}