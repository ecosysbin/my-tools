#include <iostream>
#include <vector>
#include <stack>
#include <queue>
#include <unordered_map>
#include <map>
#include <set>

using namespace std;
int main()
{
    // vector 是c++中的一种容器，可以动态地分配内存，可以存放不同类型的数据
    // 一个空的vector
    vector<int> v;
    // 一个有长度的vector
    vector<int> v2(5);
    // 一个有初始值的vector, 这里初始化了一个长度为1的vector，里面存放的是3
    vector<int> v3(1, 3);
    vector<int> v4 = {1,
                      2,
                      3,
                      4};
    // 添加元素
    v4.push_back(5);
    // 访问元素
    cout << v4[0] << endl;
    // 遍历元素
    // 可以使用下标操作符[] 或 at() 方法访问 vector 中的元素：
    for (int i = 0; i < v4.size(); i++)
    {
        cout << v4[i] << " ";
    }
    cout << endl;
    // 也可以使用迭代器遍历 vector 中的元素：
    for (vector<int>::iterator it = v4.begin(); it != v4.end(); it++)
    {
        cout << *it << " ";
    }
    cout << endl;
    // 也可以使用范围for循环遍历 vector 中的元素：
    for (int x : v4)
    {
        cout << x << " ";
    }
    cout << endl;

    // 不同于go, c++的vector可以删除和清理
    // 删除元素
    v4.erase(v4.begin() + 2);
    // vector没有string方法，无法直接全部打印出来
    cout << v4.size() << endl;
    // 清空元素
    v4.clear();
    cout << v4.size() << endl;

    // stack
    stack<int> s;
    s.push(1);
    s.push(2);
    s.push(3);
    int s1 = s.top();
    s.pop();
    int s2 = s.top();
    s.pop();
    cout << s1 << " " << s2 << endl;

    // queue
    queue<int> q;
    q.push(4);
    q.push(5);
    q.push(6);
    int q1 = q.front();
    q.pop();
    int q2 = q.front();
    q.pop();
    cout << q1 << " " << q2 << endl;

    // hash_map
    unordered_map<string, int> m;
    m["apple"] = 10;
    m["banana"] = 20;
    m["orange"] = 30;
    cout << m["apple"] << endl;
    cout << m["banana"] << endl;
    cout << m["orange"] << endl;

    // 有序Map (底层红黑树实现的)
    map<string, int> m2;
    m2["apple"] = 40;
    m2["banana"] = 50;
    m2["orange"] = 60;
    cout << m2["apple"] << endl;
    cout << m2["banana"] << endl;

    // Set 去重后的List
    set<int> s22;
    s22.insert(1);
    s22.insert(2);
    s22.insert(3);
    cout << s22.size() << endl; // 3
    for (int x : s22)
    {
        cout << x << " ";
    }
    cout << endl;
    return 0;
}