#include <stdio.h>
#include <dlfcn.h>

int (*module_pr)(void);

int main(void)
{
    void *h;

    printf("------------perform------------\n");
    h = dlopen("./libmodule.so", RTLD_NOW | RTLD_GLOBAL);
    if (!h)
    {
        printf("dlopen error:%s\n", dlerror());
        return -1;
    }

    module_pr = dlsym(h, "m_pr");
    if (dlerror() != NULL)
    {
        printf("dlsym m_pr error:%s\n", dlerror());
        return -1;
    }

    module_pr();
    dlclose(h);
    dlerror();

    printf("------------------------\n");

    return 0;
}
