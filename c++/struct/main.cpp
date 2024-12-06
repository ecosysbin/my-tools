#include <iostream>
#include <string>


using namespace std;
// 和c的结构体定义类似
struct Person {
    std::string name;
    int age;
};

// 结构体优点：

// 简单数据封装：适合封装多种类型的简单数据，通常用于数据的存储。
// 轻量级：相比 class，结构体语法更简洁，适合小型数据对象。
// 面向对象支持：支持构造函数、成员函数和访问权限控制，可以实现面向对象的设计。
// 访问权限：与 class 类似，你可以在 struct 中使用 public、private 和 protected 来定义成员的访问权限。在 struct 中，默认所有成员都是 public，而 class 中默认是 private。

// 声明一个结构体类型 Books 
struct Books
{
    string title;
    string author;
    string subject;
    int book_id;
 
    // 构造函数
    Books(string t, string a, string s, int id)
        : title(t), author(a), subject(s), book_id(id) {}
};
 
// 打印书籍信息的函数
void printBookInfo(const Books& book) {
    cout << "书籍标题: " << book.title << endl;
    cout << "书籍作者: " << book.author << endl;
    cout << "书籍类目: " << book.subject << endl;
    cout << "书籍 ID: " << book.book_id << endl;
}

int main() {
    // 结构体是一种用户自定义的数据类型，用于将不同类型的数据组合在一起。与类（class）类似，结构体允许你定义成员变量和成员函数。
    Person p1;
    p1.name = "Alice";
    p1.age = 25;

    std::cout << "Name: " << p1.name << std::endl;
    std::cout << "Age: " << p1.age << std::endl;


    // 创建两本书的对象
    Books Book1("C++ 教程", "Runoob", "编程语言", 12345);
    Books Book2("CSS 教程", "Runoob", "前端技术", 12346);
 
    // 输出书籍信息
    printBookInfo(Book1);
    printBookInfo(Book2);
    return 0;
}

// severless比较重要的三点：
// 1. 提供一个开发人员无需自己搭建的开发部署环境。 2. 服务治理能力，包括服务发现、服务路由、服务熔断、服务降级、服务限流等。 3. 弹性伸缩能力，包括按需扩缩容、弹性伸缩策略、弹性伸缩监控等。（可以缩为0）
// 云原生对serverless的定义：Fuss+Bass, Fuss是指事件驱动架构，Bass是指无状态计算（中间键的能力）。