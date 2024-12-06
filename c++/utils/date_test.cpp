#include <iostream>
#include <ctime>
#include <chrono>

using namespace std;

string getCurrentTime();
int main()
{
    // time_t now = getCurrentTime(); 不可用, 后续考虑跨文件的问题，以及测试用例的问题
    string now = getCurrentTime();
    cout << "The current time is: " << now << endl;
    // auto now_us = chrono::duration_cast<chrono::microseconds>(now_ms);
    return 0;
}