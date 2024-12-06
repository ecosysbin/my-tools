#include <iostream>

using namespace std;
class MyClass
{
public:
    static int count;
};

int MyClass::count = 10;

class Box
{
public:
    double length;
    double width;
    double height;
    // 成员函数声明
    double get(void);
    void set(double l, double w, double h);
    // 友元函数声明（即函数在类的内部实现--和java一样，也可以在外部实现--和go一样）
    // 把一个类定义为另一个类的友元类，会暴露实现细节，从而降低了封装性。理想的做法是尽可能地对外隐藏每个类的实现细节。
    // 玩的比较花先不考虑
    double getL(void)
    {
        return length;
    }
};

// 成员函数定义
double Box::get(void)
{
    return length * width * height;
}

void Box::set(double l, double w, double h)
{
    length = l;
    width = w;
    height = h;
}

int main()
{
    std::cout << "hello world" << std::endl;
    // 类的静态成员变量比较特殊，在类没有实例化之前，就可以赋值，并访问。
    std::cout << "MyClass count: " << MyClass::count << std::endl;

    Box box1;
    box1.set(10, 20, 30);
    cout << "Box1 volume: " << box1.get() << endl;

    cout << "length:" << box1.length << endl;
    // public成员变量可以直接访问
    // 需要注意的是，私有的成员和受保护的成员不能使用直接成员访问运算符 (.) 来直接访问。我们将在后续的教程中学习如何访问私有成员和受保护的成员。
    // cout << "length:" << box1::length << endl; // 不能用::访问实例的成员，类名访问静态成员可以用::
    return 0;
}

// c++可以在类内通过operator+来重载+运算符，使得类可以像内置类型一样进行运算。
// 如：
// Box box2 = box1 + box3; // 相加运算
// 重载+运算符的类必须提供一个public成员函数，该函数的返回值类型必须与类相同，且至少有一个参数，参数类型可以不同。
// 如：
// Box operator+(Box box1, Box box2)
// {
//     Box result;
//     result.length = box1.length + box2.length;
//     result.width = box1.width + box2.width;
//     result.height = box1.height + box2.height;
//     return result;
// }

// c++没有接口，可以通过抽象类（类的全部成员都是虚函数）来实现接口。