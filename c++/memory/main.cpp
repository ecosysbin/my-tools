#include <iostream>
#include <iomanip>

using namespace std;
int main()
{
	cout << "Hello, world!" << endl;
	// c++中支持malloc，但是不建议使用。new 与 malloc() 函数相比，其主要的优点是，new 不只是分配了内存，它还创建了对象。

	double *pvalue = NULL; // 初始化为 null 的指针
	pvalue = new double;   // 为变量请求内存

	*pvalue = 29494.99;								 // 在分配的地址存储值
	cout << "Value of pvalue : " << *pvalue << endl; // 打印的是29495

	// 计算机中浮点数（像 double 类型）的存储遵循 IEEE 754 标准等相关规范。浮点数并不能总是精确地表示所有的十进制小数数值，在将 29494.99 这个十进制小数转换为二进制浮点数进行存储时，会存在一定的精度损失。
	// 例如，29494.99 这个数值用二进制浮点数表示时，实际存储的值是一个最接近它的、能够用二进制浮点数格式表示的近似值，这个近似值可能在转换回十进制显示时，和原始输入的 29494.99 就出现了细微差别。
	// 默认情况下，cout 输出 double 类型时，对于小数部分会进行合适的舍入等处理，在这里由于存储的近似值以及输出格式的综合作用，就显示成了 29495 这样一个和期望的 29494.99 有差异的整数近似值。
	// 如果想要更精确地控制输出格式，使其按照期望显示出比如保留两位小数的形式，可以通过设置输出流的精度控制来实现，像下面这样修改代码：
	cout << fixed << setprecision(2);				 // 设置输出流的精度控制
	cout << "Value of pvalue : " << *pvalue << endl; // 打印的结果为 29494.99
	delete pvalue;									 // 释放内存

	return 0;
}