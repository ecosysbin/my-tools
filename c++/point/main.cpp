#include <iostream>

using namespace std;
// 下面的逻辑和go语言中的指针类似
void test(int input);
void testPtr(int* input);
int main() {
    // 定义一个变量
    int var = 20;
    // 定义一个指针变量
    int *ptr;
    // 变量的地址给指针
    ptr = &var;
    // 输出指针变量的值（变量的地址）
    cout << ptr << endl;
    // 输出指针变量指向的变量的值
    cout << *ptr << endl;

    // 指针和引用的区别
    //     引用变量是一个别名，也就是说，它是某个已存在变量的另一个名字。一旦把引用初始化为某个变量，就可以使用该引用名称或变量名称来指向变量。

    // C++ 引用 vs 指针
    // 引用很容易与指针混淆，它们之间有三个主要的不同：

    // 不存在空引用。引用必须连接到一块合法的内存。
    // 一旦引用被初始化为一个对象，就不能被指向到另一个对象。指针可以在任何时候指向到另一个对象。
    // 引用必须在创建时被初始化。指针可以在任何时间被初始化。
    // C++ 中创建引用
    // 试想变量名称是变量附属在内存位置中的标签，您可以把引用当成是变量附属在内存位置中的第二个标签。

    // 声明一个变量
    int i = 10;
    // 为i声明一个引用变量
    int& ref = i;   // & 读作引用符号,r 是一个初始化为 i 的整型引用
    cout << "Value of i : " << i << endl;  // 两者的打印都是10
    cout << "Value of i reference : " << ref  << endl;
    // 把引用作为参数	C++ 支持把引用作为参数传给函数，这比传一般的参数更安全。
    test(i);
    cout << "Value of i after function call : " << i << endl;
    cout << "Value of i reference after function call : " << ref << endl;
    test(ref);
    cout << "Value of i after function call : " << i << endl;
    cout << "Value of i reference after function call : " << ref << endl;
    // 初步看起来引用并不像指针那样。函数内修改，外部并不会受影响。

    int* ptri = &i;
    testPtr(ptri);
    cout << "Value of i after function call : " << i << endl;  // 11
    cout << "Value of  ptr after function call : " << *ptri << endl;  // 11
}

void testPtr(int* input) {
    *input = *input + 1;
}

void test(int input) {
    input = input + 1;
}