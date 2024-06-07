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