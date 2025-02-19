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

## git add 命令将文件添加到暂存区后，如果想要撤销这个操作
```
git reset HEAD <file_name>
```

## 合并最近的几次commit，打开交互式编辑器后，按i进入编辑模式，将pick改为s或squash，保存退出，然后输入提交信息，保存退出，最后输入`git push -f`强制推送。
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

## 将远端的pr合并到本地分支
(前期通过git pull volcano pull/3874/head:network-topology将远端未提交的pr拉倒了本地，在远端pr合入后需要在本地进行rebase合并)
```
1. 先更新本地仓库信息，确保最新
git fetch volcano
2. 进行rebase
git rebase volcano/network-topology
3. 解决冲突
4. 提交冲突后，继续rebase
git add *
git rebase --continue
5. 推送到远端仓库
git push -f origin network-topology
```

## 获取当前分支的commit哈希值
git rev - parse --verify HEAD是一个 Git 命令，主要用于获取当前分支（HEAD）所指向的提交（commit）的哈希值，并验证其是否存在。

## 将远端未合入的PR和并到本地，但是git pull的时候会有冲突。 强行git pull -f又会覆盖本地的所有修改。
1. 先在本地创建新的分支，并切换到该分支
```
git checkout -b new_branch
```
2. 然后执行如下命令，将远端未合入的PR和并到本地（有冲突时可以直接-f）
```
git pull -f volcano pull/3965/head:network-topology
```
2.1 git pll 可能仍然会冲突，这时简单起见，可以先git checkout到pull PR之前的commit，然后git pull
```
git checkout <commit-id>
```

3. 将原来的分支已提交的commit Cherry pick到新的分支
```
git cherry-pick <commit-id>"
```
4. 解决冲突
4.1 解决冲突后执行如下命令，将解决冲突后的文件添加到暂存区
```
git add <file>
```
4.2 执行如下命令，继续解决冲突
```
git cherry-pick --continue
```
4.3 执行git add xxx 命令将修改提交到缓冲区
4.4 执行git checkout -b temp 将branch从* (HEAD detached from 04ce0fc20)切换到新分支

5. 将新分支覆盖原分支(可以确保已创建的PR不变化，不然可以使用新分支重新创建PR)
5.1 切换到原来分支
git checkout network-topology  # 切换到 network-topology 分支
5.2 将新分支强行覆盖原来分支
git reset --hard temp  # 强制将 network-topology 分支重置为 temp 分支的状态

6. commit并push到远端仓库

7. github 将远端仓库的clone到本地仓库，假设最终要向远程xxx分支提交代码，本地不建议直接在xxx分支提交代码。这样假如有多个特性开发则都可以从xxx分支checkout进行开发。

8. github提交pending的评论
