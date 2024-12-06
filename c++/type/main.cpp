# include <iostream>

using namespace std;

// 变量声明
extern int a, b;
extern int c;
extern float f;


const string name = "jerry";
#define age 20;

void example_function(); // 同c，函数声明

int main() {
    cout << "hello world" << endl;
    int i = 10;
    // static cast
    // 静态转换通常用于比较类型相似的对象之间的转换，例如将 int 类型转换为 float 类型。
    // 静态转换不进行任何运行时类型检查，因此可能会导致运行时错误。
    float j = static_cast<float>(i);
    cout << "j = " << j << endl;
    // Dynamic cast
    // 动态转换用于将基类指针或引用转换为派生类指针或引用。
    // 态转换在运行时进行类型检查，如果不能进行转换则返回空指针或引发异常。
    class Base {
        public:
        virtual void func(){}
    };
    class Derived : public Base {};
    Base* ptr_base = new Derived;
    Derived* ptr_derived = dynamic_cast<Derived*>(ptr_base);
    // the operand of a runtime dynamic_cast must have a polymorphic class type
    // 运行时动态转换的操作数必须是多态类类型。
    // 一个类是多态的，当且仅当它至少包含一个虚函数。
    // 注意，即使Parent成为多态类，dynamic_cast在运行时如果发现转换不合法（如p实际上不是Child类型），会返回nullptr（对于指针类型的转换）或者抛出std::bad_cast异常（对于引用类型的转换）
    
    // Const Cast
    // 常量转换用于将 const 类型的对象转换为非 const 类型的对象。
    // 常量转换只能用于转换掉 const 属性，不能改变对象的类型。
    const int k = 10;
    int& l = const_cast<int&>(k);
    // the type in a const_cast must be a pointer, reference, or pointer to member to an object type
    l = 20; // l 虽然能被修改值，但是打印时还是10
    cout << "l = " << i << endl; 
    cout << "k = " << k << endl; 
    // Reinterpret Cast
    // 重新解释转换将一个数据类型的值重新解释为另一个数据类型的值，通常用于在不同的数据类型之间进行转换。
    // 重新解释转换不会对数据进行任何检查，因此可能会导致运行时错误。
    // int i = 10;
    // float f = reinterpret_cast<float&>(i); // 重新解释将int类型转换为float类型

    // 变量定义
    // 用 extern 关键字在任何地方声明一个变量。虽然您可以在 C++ 程序中多次声明一个变量，但变量只能在某个文件、函数或代码块中被定义一次
    int a, b;
    int c;
    float f;
    // int a; 编译时会报错 redeclared
    // 实际初始化
    a = 10;
    b = 20;
    c = a + b;

    // 左值（lvalue）：指向内存位置的表达式被称为左值（lvalue）表达式。左值可以出现在赋值号的左边或右边。
    // 右值（rvalue）：术语右值（rvalue）指的是存储在内存中某些地址的数值。右值是不能对其进行赋值的表达式，也就是说，右值可以出现在赋值号的右边，但不能出现在赋值号的左边。
    // 变量是左值，因此可以出现在赋值号的左边。数值型的字面值是右值，因此不能被赋值，不能出现在赋值号的左边。下面是一个有效的语句：

    // const和define都可以定义全局变量，但是define定义的全局变量在打印时无法被追加，而const定义的全局变量可以被追加。
    // 建议使用const定义全局变量，因为const定义的全局变量可以被优化，而define定义的全局变量不能被优化。
    // cout << age << endl; // 打印age的值
    cout << age; 

    short int q;           // 有符号短整数
    short unsigned int w;  // 无符号短整数
    w = 50000;
 
    q = w;
    cout << q << " " << w; // 会输出-15536 50000，即使用中不要将无符号短整数赋值给有符号短整数，否则会出现数据溢出。

    for (i = 0; i < 10; i++){
        example_function();
    }
    return 0;
}

class Example {
public:
    int get_value() const {
        return value_; // const 关键字表示该成员函数是常量成员函数，即通过const声明的实例（const Example obj; obj.get_value()）才能使用。则不会修改对象中的数据成员
    }
    void set_value(int value) const {
        value_ = value; 
    }
private:
    mutable int value_; // 但是mutable 关键字允许在 const 成员函数中修改成员变量，又打破了const常量成员函数的限制
};

void example_function() {
    static int count = 0; // static 关键字使变量 count 存储在程序生命周期内都存在
    count++;
    cout << "Count: " << count << endl;
}