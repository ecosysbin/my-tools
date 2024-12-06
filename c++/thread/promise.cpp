#include <thread>
#include <future>
#include <utility>
#include <chrono>
#include <iostream>

void compute(int a, int b, std::promise<int> promise_)
{
    int sum = a + b;
    // 模拟耗时
    std::this_thread::sleep_for(std::chrono::seconds(5));
    promise_.set_value(sum);
}

int main()
{
    std::promise<int> promise_;
    std::future<int> future_ = promise_.get_future();
    std::thread t(compute, 3, 4, std::move(promise_));
    std::cout << "waiting for compute result..." << std::endl;
    // 在主线程中，future_ 对象调用 get() 可返回被赋值的计算结果，如果计算结果未就绪，get() 会等待并阻塞当前线程。
    std::cout << "get compute result:"
              << future_.get()
              << std::endl;
    std::cout << "compute done." << std::endl;
    t.join();
    return 0;
}