1. 交互式debug
在代码中加入breakpoint()
 68                 parameters = model.parameters()
 69                 breakpoint()
 70  ->             loss.backward()
2. 运行代码即会停在breakpoint的位置

3. l和ll可以展示breakpoint的上下文代码片段

4. n可以跳到下一行

5. q可以退出交互式debug