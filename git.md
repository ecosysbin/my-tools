# git常用命令

## 查到最近的提交记录(会返回commit-id)
```
git log -n 1
```

## 回到某个提交版本
```
git checkout <commit-id>
```

## 清理未跟踪文件
```
git clean -fd
```

## 切换到远程分支
```
git checkout origin/dev
```

## 创建并切换到本地新分支
```
git checkout -b dev
```

## 删除本地分支
```
git branch -d dev
```

## 删除远程分支
```
git push origin --delete dev
```




