#include <iostream>

using namespace std;
#define PI 3.14159
int main()
{
    // 变量、语句之间为什么要加空格？这样才能让编译器识别语句中的某个元素从哪里结束，下个元素从哪里开始。
    // cout是ostream类的一个对象，是extern声明的全局对象，cout代表标准输出设备，通常是控制台。
    // <<运算符被重载用于ostream类，使得它可以将各种数据类型（如整数、浮点数、字符、字符串等）发送到输出流中。
    // C++ 的设计理念强调数据抽象和面向对象编程，因此cout和<<运算符的使用频率很高。
    // 类型安全性：
    // 例如，当你使用cout << 5;时，C++ 编译器会根据5的类型（这里是int）来调用ostream类中对应的operator<<重载函数。这种基于类型的重载机制提供了更好的类型安全性。与printf不同，printf需要在格式字符串中准确地指定输出类型，如果格式字符串与实际参数类型不匹配，可能会导致错误（如输出乱码或程序崩溃）。而cout <<会根据参数类型自动选择正确的输出方式。
    // 易用性和可读性：
    // cout <<的语法更加直观和易于理解。它采用了一种类似于 “流” 的方式，你可以将多个输出操作连接在一起。例如，cout << "The value of x is: " << x << endl;可以很方便地将一个字符串和一个变量的值依次输出到控制台。这种链式调用的方式使得代码更加简洁和易读，而不需要像printf那样在一个格式字符串中组合各种输出内容。
    // 可扩展性：
    // C++ 程序员可以为自定义的数据类型重载operator<<运算符。例如，如果你定义了一个名为Complex的复数类，你可以通过重载operator<<来定义如何将Complex对象输出到控制台。这样，当你使用cout <<来输出Complex对象时，就可以按照你自定义的格式进行输出。
    // 以下是一个简单的示例：
    // #include <iostream>
    // class Complex {
    // public:
    //     double real;
    //     double imag;
    //     Complex(double r = 0, double i = 0) : real(r), imag(i) {}
    // };
    // std::ostream& operator<<(std::ostream& os, const Complex& c) {
    //     os << c.real << " + " << c.imag << "i";
    //     return os;
    // }
    // int main() {
    //     Complex c(3.0, 4.0);
    //     std::cout << c << std::endl;
    //     return 0;
    // }
    // 在这个例子中，通过重载operator<<，可以很方便地将Complex对象以a + bi的形式输出到控制台。这展示了cout <<在处理自定义数据类型时的强大可扩展性。
    cout << "hello world" << endl;
    // 使用了标准命名空间 std 对 cout 进行了限定。在 C++ 中，cout 是定义在 std 命名空间中的标准输出流对象。当使用 std::cout 时，明确地指定了 cout 所属的命名空间，这种方式更加严谨，适合在复杂的程序中避免命名冲突。
    // 没有显式地使用命名空间限定。这种写法在实际使用中更为简洁，但前提是在当前的代码环境中已经通过 using namespace std; 等方式将 std 命名空间引入到了当前作用域。如果没有引入 std 命名空间而直接使用 cout，编译器会报错，提示找不到 cout 的定义。
    // 命名空间的主要目的是避免命名冲突，将相关的代码元素组织在一起，使得不同命名空间中的同名标识符可以共存.
    // 命名空间和类都可以用::访问成员
    std::cout << "hello world" << endl;

    // #打头即预处理指令
    // #define 预处理指令用于创建符号常量。该符号常量通常称为宏，指令的一般形式是：
    // #define macro-name replacement-text
    // 当这一行代码出现在一个文件中时，在该文件中后续出现的所有宏都将会在程序编译之前被替换为 replacement-text。例如：
    // gcc -E test.cpp > test.p会看到预处理之后的源代码。
    cout << "Value of PI :" << PI << endl;

    // 系统预定义的宏
    cout << "Value of __LINE__ : " << __LINE__ << endl; // 当前行号
    cout << "Value of __FILE__ : " << __FILE__ << endl; // 当前文件名
    cout << "Value of __DATE__ : " << __DATE__ << endl; // 编译日期
    cout << "Value of __TIME__ : " << __TIME__ << endl; // 编译时间
    return 0;
}