#!/bin/bash

# 安装 expect
apt install expect -y

# 受控主机 ip
SERVERS_STR=$(ansible all --list-hosts)
SERVERS=${SERVERS_STR: 13}
echo "$SERVERS"

# root密码
PASSWD="123456"

function sshCopyId {
	expect -c "
	set timeout -1;  
	spawn ssh-copy-id -o stricthostkeychecking=no $1;
	expect {
		\"yes/no\" { send \"yes\r\" ;exp_continue; } 
		\"password:\" { send \"$PASSWD\r\";exp_continue; }
	};
	expect eof;"
}
 
for server in $SERVERS
do
	sshCopyId "$server"
done
