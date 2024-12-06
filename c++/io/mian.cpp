#include <iostream>

using namespace std;

int main()
{
    string name;
    cout << "Hello, world!" << endl;
    // 输出流 cout 用于输出到屏幕或文件
    // 输入流cin用于从键盘输入数据
    cout << "please input your name: ";
    cin >> name;
    cout << "your name is " << name << endl;

    // 预定义的对象 cerr 是 iostream 类的一个实例。cerr 对象附属到标准输出设备，通常也是显示屏，但是 cerr 对象是非缓冲的，且每个流插入到 cerr 都会立即输出。
    // cerr 也是与流插入运算符 << 结合使用的
    string msg = "Error: invalid input";
    cerr << "error msg:" << msg << endl;  // 从打印看好像也和cont一样。

    // 预定义的对象 clog 是 iostream 类的一个实例。clog 对象附属到标准输出设备，通常也是显示屏，但是 clog 对象是缓冲的。这意味着每个流插入到 clog 都会先存储在缓冲区，直到缓冲填满或者缓冲区刷新时才会输出。
    // clog 也是与流插入运算符 << 结合使用的
    string log = "Warning: file not found";
    clog << "log msg:" << log << endl;  // 从打印看好像也和cerr一样。

    // 但在编写和执行大型程序时，它们之间的差异就变得非常明显。所以良好的编程实践告诉我们，使用 cerr 流来显示错误消息，而其他的日志消息则使用 clog 流来输出
    return 0;
}