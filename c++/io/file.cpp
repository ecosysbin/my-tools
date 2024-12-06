#include <iostream>
#include <fstream>

using namespace std;

int main()
{
    char data[100];

    // 写模式打开文件
    ofstream outfile;
    // 不存在就会创建
    outfile.open("./hello.txt");
    cout << "enter your name: ";
    cin.getline(data, 100);

    // 向文件写入用户输入的数据
    outfile << data << endl;

    // 关闭文件
    outfile.close();

    // 读模式打开文件
    ifstream infile;
    infile.open("./1hello1.txt");
    cout << "reading from file:" << endl;
    infile >> data;
    cout << data << endl;

    // 再次从文件读取数据，并显示它
    //    infile >> data;
    //    cout << data << endl;

    // 关闭文件
    infile.close();
    return 0;
}