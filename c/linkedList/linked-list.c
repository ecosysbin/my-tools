#include <stdio.h>
#include <stdlib.h>

typedef struct Node
{
    int data;
    struct Node *next;
} MyNode;

void print_Linded_List(MyNode *head)
{
    MyNode *tmp = head;
    while (tmp != NULL)
    {
        printf("data: %d\n", tmp->data);
        tmp = tmp->next;
    }
}
MyNode *revert_Linked_List(MyNode *head)
{
    MyNode *prev = NULL;
    MyNode *current = head;
    MyNode *next = NULL;
    while (current != NULL)
    {
        printf("current:%d, prev:%d, next:%d\n", current->data, prev == NULL ? -1 : prev->data, next == NULL ? -1 : next->data);
        /* code */
        next = current->next;
        current->next = prev;
        prev = current;
        current = next;
    }
    // print_Linded_List(prev);
    // head = prev; // 这里不能直接赋值，因为这里的prev是反转后的链表的头结点，而head是原链表的头结点
    return prev;
}

MyNode *create_Node(int data)
{
    MyNode *newNode = (MyNode *)malloc(sizeof(MyNode));
    newNode->data = data;
    newNode->next = NULL;
    return newNode;
}

void insert_at_tail(MyNode *head, int data)
{
    MyNode *newNode = create_Node(data);
    if (head == NULL)
    {
        head = newNode;
        return;
    }
    else
    {
        MyNode *tmp = head;
        while (tmp->next != NULL)
        {
            tmp = tmp->next;
        }
        tmp->next = newNode;
    }
}

// int main(void)
// {
//     MyNode head = {
//         .data = 0,
//     };
//     insert_at_tail(&head, 1);
//     insert_at_tail(&head, 2);
//     insert_at_tail(&head, 3);
//     insert_at_tail(&head, 4);
//     print_Linded_List(&head);
//     // 反转链表
//     MyNode *newHead = revert_Linked_List(&head);
//     print_Linded_List(newHead);
//     return 0;
// }
