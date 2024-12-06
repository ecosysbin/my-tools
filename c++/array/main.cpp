# include <iostream>
# include <stdio.h>
using namespace std;

int main() {
    // 声明一个长度为10的整型数组，或int n[ 10 ]
    int n[10];
    // c++和c一样，声明的数组内的元素不会清理，也会有越界风险，需要程序员自己负责
    for (int i = 0; i < 11; i++) {
        cout << n[i] << endl;
    }
    // c++和c一样，printf函数可以用来输出字符串到控制台。语法可以混用。
    printf("Hello c++, c is comming !\n");
    return 0;
}