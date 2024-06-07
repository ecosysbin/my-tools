#include <stdio.h>
#include <string.h>
#include <stdlib.h>

// 定义结构体别名
typedef struct student student_def;

struct student
{
    int name;
    int age;
};

int main()
{
    // 使用原生结构体定义变量
    struct student stu1 = {20, 30};
    printf("%d %d\n", stu1.name, stu1.age);

    // 使用结构体别名定义变量
    student_def stu;
    stu.name = 10;
    stu.age = 20;
    printf("%d %d\n", stu.name, stu.age);

    // 取结构体指针变量
    student_def *pstu = &stu;
    (*pstu).name = 40;
    (*pstu).age = 50;
    printf("%d %d\n", stu.name, stu.age);

    // 指针变量简单赋值pstu->属性 = 值，同理(*pstu).属性 = 值
    pstu->name = 60;
    pstu->age = 70;
    printf("%d %d\n", stu.name, stu.age);
};