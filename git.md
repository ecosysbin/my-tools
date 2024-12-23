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

## 删除本地分支(强制删除)（看起来远端分支也会被一起删除）
```
git branch -D dev
```


## 删除远程分支
```
git push origin --delete dev
```

## 保留本地工作区，更改commit到某次提交
```
git reset --soft <commit-id>
```

## 合并最近的几次commit
```
git rebase -i HEAD~3
```

## 强制推送（git rebase -i后会提示落后一个分支，pull后提交结果之前的commit还在，还会出现空的merge commit, 这时其实可以强制推送）
```
git push -f
```

## 提交时增加签名
```
git commit -s -m "commit message"
```

## 添加一个远端仓库到本地remote
```
git remote add volcano https://github.com/volcano-sh/volcano.git
```

## 将远端仓库的pr合并到本地分支
```
git pull volcano pull/3874/head:network-topology
```

