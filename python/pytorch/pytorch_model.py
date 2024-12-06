import torch
from torch.utils.data import DataLoader
from torch.utils.data.dataset import TensorDataset
 
# 自构建数据集
dataset = TensorDataset(torch.arange(1, 40))
dl = DataLoader(dataset,
                batch_size=10,
                shuffle=True,
                num_workers=1,
                drop_last=True)
# 数据输出
for batch in dl:
    print(batch)