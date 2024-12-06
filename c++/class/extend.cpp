#include <iostream>
using namespace std;

class Animal
{
public:
    // const 成员函数可以被常量对象调用. 如：
    // const Animal myAnimal;
    // myAnimal.sound(); // 正确
    // 虚函数派生类可以不重写，要是纯虚函数（只有声明，没有实现）则派生类必须实现。
    // virtual int area() = 0;  用=0表示纯虚函数。
    // java是不是都是纯虚函数？是的，但是名称叫抽象方法。通过方法重载可以达到c++虚函数的效果，this调用的自己，super调用的父类。
    virtual void sound() const
    {
        cout << "Animal sound" << endl;
    }
    // 析构函数在实例对象被销毁时被调用.主要用来释放资源.
    virtual ~Animal()
    {
        cout << "Animal destructor" << endl;
    }
};

class Dog : public Animal
{
public:
    // 要是没有重写虚函数，调用时会调用父类的虚函数。
    void sound() const
    {
        cout << "Dog sount" << endl;
    }
    // 要是没有实现析构函数，用父类引用调用时，会只会调用父类的析构函数。实现后先调用子类的析构函数，再调用父类的析构函数。
    ~Dog()
    {
        cout << "Dog destructor" << endl;
    }
    // 也可以有构造函数，通过参数变化区分不同的构造函数，默认构造函数参数为空，使用方式同java。
};

class Cat : public Animal
{
public:
    void sound() const
    {
        cout << "Cat sound" << endl;
    }
    ~Cat()
    {
        cout << "Cat destructor" << endl;
    }
};

int main()
{
    Animal *pAnimal; // 基类指针

    // 创建Dog对象，指向它的基类指针
    pAnimal = new Dog;
    pAnimal->sound(); // 调用Dog的sound()函数
    delete pAnimal;   // 释放Dog对象

    // 创建Cat对象，指向它的基类指针
    pAnimal = new Cat;
    pAnimal->sound(); // 调用Cat的sound()函数
    delete pAnimal;   // 释放Cat对象

    Dog myDog;
    // 通过对象字面量,通过.调用属性和函数，new创建的对象返回的指针通过->调用属性和函数
    myDog.sound(); // 调用Dog的sound()函数
    return 0;
}