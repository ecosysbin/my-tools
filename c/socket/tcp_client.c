#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <errno.h>

#define SERVER_IP "127.0.0.1"
#define PORT 8888
#define BUFFER_SIZE 1024

int main()
{
    int sock;
    struct sockaddr_in server_addr;
    char buffer[BUFFER_SIZE];

    // 创建套接字
    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) == -1)
    {
        perror("socket creation failed");
        exit(EXIT_FAILURE);
    }

    // 初始化服务器地址结构体
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(PORT);
    if (inet_pton(AF_INET, SERVER_IP, &server_addr.sin_addr) <= 0)
    {
        perror("inet_pton failed");
        close(sock);
        exit(EXIT_FAILURE);
    }

    // 连接服务器
    if (connect(sock, (struct sockaddr *)&server_addr, sizeof(server_addr)) == -1)
    {
        perror("connect failed");
        close(sock);
        exit(EXIT_FAILURE);
    }

    printf("已连接到服务器\n");

    // 循环发送和接收数据
    while (1)
    {
        memset(buffer, 0, BUFFER_SIZE);
        // 从用户输入获取数据, 可能存在gcc读不到STDIN的问题
        printf("请输入要发送给服务器的消息: ");
        fgets(buffer, BUFFER_SIZE, STDIN);

        // 去除换行符
        buffer[strcspn(buffer, "\n")] = '\0';

        // 发送数据到服务器
        if (send(sock, buffer, strlen(buffer), 0) == -1)
        {
            perror("send failed");
            break;
        }

        // 接收服务器响应
        memset(buffer, 0, BUFFER_SIZE);
        ssize_t bytes_read = recv(sock, buffer, BUFFER_SIZE - 1, 0);
        if (bytes_read <= 0)
        {
            if (bytes_read == 0)
            {
                printf("服务器已关闭连接\n");
            }
            else
            {
                perror("recv failed");
            }
            break;
        }
        buffer[bytes_read] = '\0';
        printf("从服务器接收到响应: %s\n", buffer);
    }

    // 关闭套接字
    close(sock);

    return 0;
}