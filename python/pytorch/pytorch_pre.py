import torch

print(torch.__version__) # pytorch版本
print(torch.version.cuda) # cuda版本
print(torch.cuda.is_available()) # 查看cuda是否可用

# Tensor张量是Pytorch里最基本的数据结构。直观上来讲，它是一个多维矩阵，支持GPU加速
# 定义一个2（行） * 3（列）的张量
t = torch.Tensor(2, 3)
print(t)

# 定义一个一维张量，包含三个元素
t1 = torch.Tensor([1,2,3])
print(t1)

# 生成随机浮点数张量, 大小为2 * 3
t2 = torch.randn(2,3)
print(t2)

# 生成随机整数张量，大小为10（最大值为10且为十个元素）
t3 = torch.randperm(10)
print(t3)

# 生成一个一维张量，三个参数分别为起始位、终止位、步长。该方法将被废弃，推荐使用torch.arange()
t4 = torch.range(1, 10, 2)
print(t4)

# tensor运算
# 张量绝对值
t5 = torch.abs(t2)
print(t5)

# 张量相加
t6 = t2 + t5
print(t6)

# 张量的n次方
t7 = torch.pow(t2,2)
print(t7)

# 张量相乘
t8 = torch.mul(t2,t2)
print(t8)


# 张量转化位numpy类型（看着是转化成了python的通用数组）
n1 = t6.numpy()
print(n1)

# 查看尺寸, 返回类型就是torch.Size([2, 3])
i1 = t6.size()
print(i1)

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

print(f"Using {device} device")

# 判断某个对象是在什么环境中运行的，默认打印cpu
print(t2.device)

# 张量转移到cuda环境中, 张量默认第一个参数是位置，第二个参数是设备，要是没有第二参数，默认是cpu
print("===============================================")
t2g =t2.to(device)
print(t2)
print(t2g)

# 环境变量设为cpu
t2c = t2.cpu()
print(t2c)

t2g =t2.to(device)
# 也是将变量转移到cpu环境中
t2c1= t2g.to("cpu")
print(t2c1)

# 若一个没有环境的对象与另外一个有环境a对象进行交流,则环境全变成环境a? 
# RuntimeError: Expected all tensors to be on the same device, but found at least two devices, cuda:0 and cpu!
# t25 = t2g + t2
# print(t25)

# cuda环境下tensor不能直接转化为numpy类型,必须要先转化到cpu环境中
t2gn = t2g.cpu().numpy()
print(t2gn)

# 直接创建cuda类型的张量. torch.tensor([1,2],device)会报错: TypeError: tensor() takes 1 positional argument but 2 were given
td = torch.tensor([1,2],device=device)
print(td)


