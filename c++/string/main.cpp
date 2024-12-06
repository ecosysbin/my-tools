# include <iostream>
# include <cstring>

using namespace std;
int main() {
   // c++ 支持c风格的string. 如下同char site[] = "RUNOOB";
   char site[7] = {'R', 'U', 'N', 'O', 'O', 'B', '\0'};
 
   cout << "菜鸟教程: ";
   cout << site << endl;  // 输出 "菜鸟教程: RUNOOB"

   // c++并且提供了大量对c字符串的操作函数。
   cout << strlen(site) << endl; 

   // 同时c++标准库提供了string类, 用来管理字符串.
   string str1 = "RUNOOB";
   string str2 = "Googlea";
   string str3;
   int len;
   // 复制str1到str3
   str3 = str1;
   cout << "str3 = " << str3 << endl;
   // 连接str1和str2
   str3 = str1 + str2;
   cout << "str3 = " << str3 << endl;
   // 获取str3的长度
   len = str3.size();
   cout << "str3 length = " << len << endl;
   return 0;
}