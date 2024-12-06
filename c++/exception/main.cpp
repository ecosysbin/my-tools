#include <iostream>

using namespace std;

double division(int a, int b)
{
    if (b == 0)
    {
        throw "Division by zero!";
    }
    return (a / b);
}

int main()
{
    int x = 50;
    int y = 0;
    double z = 0;
    try
    {
        z = division(x, y);
        cout << z << endl;
    }
    // 这里用const char* 或者string* 都可以，因为string继承自const char*
    // 通过...捕获所有异常
    catch (...)
    {
        // cout << msg << endl;
        cout << "Error: Division by zero" << endl;
    }
    return 0;
}