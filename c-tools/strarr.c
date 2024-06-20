#include <stdio.h>

char ***groupAnagrams();
int main(void)
{
    // 申请字符串数组（从定义看等同于char引用类型的数组）
    char *strs[] = {"eat", "tea", "tan", "ate", "nat", "bat"};
    groupAnagrams(strs, 6, NULL, NULL);
    return 0;
}

char ***groupAnagrams(char **strs, int strsSize, int *returnSize, int **returnColumnSizes)
{
    for (int i = 0; i < strsSize; i++)
    {
        printf("%s\n", strs[i]);
        for (int j = 0; j < 3; j++)
        {
            // 通过char引用的角标可以获取字符串的字符
            printf("%c\n", strs[i][j]);
        }
    }
    return NULL;
}