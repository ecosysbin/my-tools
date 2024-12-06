import torch
import numpy as np

print(torch.__version__) # pytorch版本
print(torch.version.cuda) # cuda版本
print(torch.cuda.is_available()) # 查看cuda是否可用


# 标量Tensor求导
# 求 f(x) = a*x**2 + b*x + c 的导数
x = torch.tensor(-2.0, requires_grad=True)
a = torch.tensor(1.0)
b = torch.tensor(2.0)
c = torch.tensor(3.0)
y = a*torch.pow(x,2)+b*x+c
y.backward() # backward求得的梯度会存储在自变量x的grad属性中
dy_dx =x.grad
# 会先求f（x）的导函数，然后计算x=-2.0时的值就是dy_dx,使用x**2 更好理解
print(dy_dx)


# 非标量Tensor求导
# 求 f(x) = x**2 的导数，f(x)的导函数为2*x
x = torch.tensor([[-2.0,-1.0],[0.0,1.0]], requires_grad=True)
a = torch.tensor(1.0)
b = torch.tensor(2.0)
c = torch.tensor(3.0)
gradient=torch.tensor([[1.0,1.0],[1.0,1.0]])
y = torch.pow(x,2)
y.backward(gradient=gradient) # backward就是求导数
dy_dx =x.grad
print(dy_dx)
# 打印tensor([[-4., -2.], [ 0.,  2.]])

# 一阶导数的定义：一阶导数是函数在某一点的导数，描述了函数在该点附近的变化率。它表示函数值随自变量变化的速率。‌
# 二阶导数的定义：二阶导数是一阶导数的导数，即函数的二阶导数表示一阶导数的变化率。从原理上看，二阶导数表示一阶导数的变化率；从图形上看，它反映的是函数图像的凹凸性
 
#单个自变量求导
# 求 f(x) = x**4 的导数
x = torch.tensor(1.0, requires_grad=True)
a = torch.tensor(1.0)
b = torch.tensor(2.0)
c = torch.tensor(3.0)
y = torch.pow(x, 4)
#create_graph设置为True,允许创建更高阶级的导数
#求一阶导 4x^3
dy_dx = torch.autograd.grad(y, x, create_graph=True)[0]
#求二阶导 12x^2
dy2_dx2 = torch.autograd.grad(dy_dx, x, create_graph=True)[0]
#求三阶导 24x
dy3_dx3 = torch.autograd.grad(dy2_dx2, x)[0]
print(dy_dx.data, dy2_dx2.data, dy3_dx3)
# 打印tensor(4.), tensor(12.), tensor(24.)
 
# 多个自变量求偏导
# 偏导数表示在其他变量保持不变的情况下，函数对某一变量的变化率
x1 = torch.tensor(1.0, requires_grad=True)
x2 = torch.tensor(2.0, requires_grad=True)
y1 = x1 * x2
y2 = x1 + x2
#只有一个因变量,正常求偏导
dy1_dx1, dy1_dx2 = torch.autograd.grad(outputs=y1, inputs=[x1, x2], retain_graph=True)
print(dy1_dx1, dy1_dx2)
# 若有多个因变量，则对于每个因变量,会将求偏导的结果加起来
dy1_dx, dy2_dx = torch.autograd.grad(outputs=[y1, y2], inputs=[x1, x2])
# dy1_dx, dy2_dx
print(dy1_dx, dy2_dx)

# 求最小值
# 使用自动微分机制配套使用SGD优化器随机梯度下降来求最小值
# 1e-3 是一个科学计数法表示的数值，等价于 (1 \times 10^{-3})，即 0.001。
# 因此，lr=1e-3 表示学习率为 0.001。
#例2-1-3 利用自动微分和优化器求最小值

# f(x) = a*x**2 + b*x + c的最小值
x = torch.tensor(0.0, requires_grad=True)  # x需要被求导
a = torch.tensor(1.0)
b = torch.tensor(-2.0)
c = torch.tensor(1.0)
optimizer = torch.optim.SGD(params=[x], lr=0.01)  #SGD为随机梯度下降
print(optimizer)

def f(x):
    result = a * torch.pow(x, 2) + b * x + c
    return (result)
 
for i in range(500):
    optimizer.zero_grad()  #将模型的参数初始化为0
    y = f(x)
    y.backward()  #反向传播计算梯度
    optimizer.step()  #更新所有的参数
print("y=", y.data, ";", "x=", x.data)