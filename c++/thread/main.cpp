#include <iostream>
#include <thread>

// c不直接支持多线程，则c++的出现有了更大意思，加上面向对象编程的概念，使得c++可以更好地支持多线程编程，更好的支持复杂的业务逻辑，高并发的程序。
// 再加上c++天生的支持string类型等方便的操作，相比于其他高级语句，又有更好的性能，则c++永远都会有江湖地位。

// 定义一个thread_local类型的变量，具有线程局部性，各线程对其进行修改互不干扰
thread_local int threadSpecificVar = 0;

void threadFunction(int ids, int name, int age)
{
    threadSpecificVar = ids;
    // std::cout << "Thread " << id << ": threadSpecificVar = " << threadSpecificVar << std::endl;
    for (int i = 0; i < 10; i++)
    {
        std::cout << "Thread " << ids << std::endl;
    }
}
class PrintTask
{
public:
    // operator 并不是一个普通意义上的函数名，它是 C++ 中的关键字，用于操作符重载。其使用有严格的规则限制，不能随意将其当作普通函数名来使用,也不能拼在函数名称中。
    // 这里的 operator() 就是一个重载了 () 操作符的函数，它是一个函数对象，可以像函数一样被调用。
    void operator()(int count) const
    {
        for (int i = 0; i < count; ++i)
        {
            std::cout << "Hello from thread (function object)!\n";
        }
    }
};

void increment(int &x)
{
    ++x;
}

int main()
{
    // 多线程执行函数可以自定义参数，创建线程实例时，传入参数。没有参数的函数可以直接创建线程实例。参数名程也可以自定义。
    // 一致性和简洁性：C++ 标准库中的许多类型都支持直接在声明时通过构造函数进行初始化，这种方式提供了一种统一和简洁的对象创建和初始化语法。
    // 对于 std::thread 对象，使用 std::thread t1(threadFunction, 1); 这样的构造函数初始化形式与其他标准库类型的使用方式保持一致，增强了代码的一致性和可读性，使开发者能够更自然地编写和理解多线程相关的代码。
    // 1. 通过函数指针
    std::thread t1(threadFunction, 1, 1, 2);
    std::thread t2(threadFunction, 2, 2, 2);

    // t1.join();
    // t2.join();
    t1.detach();
    t2.detach();

    // 线程是程序执行中的单一顺序控制流，多个线程可以在同一个进程中独立运行。
    // 线程共享进程的地址空间、文件描述符、堆和全局变量等资源，但每个线程有自己的栈、寄存器和程序计数器

    // 并发：多个任务在时间片段内交替执行，表现出同时进行的效果。
    // 并行：多个任务在多个处理器或处理器核上同时执行。

    // 2. 通过函数对象
    std::thread t3(PrintTask(), 5); // 创建线程，传递函数对象和参数，使用了operator()重载的函数对象
    t3.join();                      // 等待线程完成

    // 3. 通过lambda表达式, 实际同1. 就是匿名函数。
    // join() 用于等待线程完成执行， detach() 将线程与主线程分离，线程在后台独立运行，主线程不再等待它。

    int num = 0;
    // c++明确的区分了值传递和引用传递，使用ref关键字可以传递引用。
    std::thread t(increment, std::ref(num)); // 使用 std::ref 传递引用
    t.join();
    std::cout << "Value after increment: " << num << std::endl;
    return 0;
}

// 对spring cloud的思考，spring cloud做的比较好的是在业务庞大的时候对微服务拆分后的管理（关键是服务注册与发现）和监控。对各微服务来讲，
// 只需要关注自己的业务逻辑，而不用关心服务注册与发现，监控等问题。
// 这在单一的网络、存储架构下是没问题的，甚至容器化，迁移上kubernetes也是有意义的，毕竟kubernetes要做到服务治理上istio对服务改造还是有一些代价。
// 但是要是做多网络架构，如：业务的管理网和业务网络和存储网络分开，则需要考虑服务注册与发现，监控等问题。另外存储卷的管理也是个问题，传统的玩法是将存储卷挂在宿主机上，这样数据的安全

// go和c++相比，从软件工程化角度，有包的概念，这样有利于多文件项目的管理，整个工程代码看起来更加的结构化。