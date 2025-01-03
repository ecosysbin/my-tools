# C、C++、Go 和 Java 的内存分配分区存在着一些不同，主要体现在以下几个方面：
## 栈内存
C 和 C++
手动管理：栈内存的分配和释放主要由编译器自动完成，但程序员也可以通过一些特殊的指令或函数进行一定程度的手动控制。
大小固定：栈内存的大小通常在程序启动时就已经确定，一般相对较小，取决于操作系统和编译器的默认设置，如在 Linux 系统上使用 GCC 编译器时，默认栈大小一般是 8MB。
变量存储：局部变量、函数参数等通常存储在栈内存中。当函数被调用时，这些变量会被压入栈中，函数执行完毕后再弹出栈。
Go
自动管理：由 Go 运行时自动管理，程序员一般无需手动干预栈内存的分配和释放。
动态调整：每个 goroutine 的栈大小初始时通常较小，一般为 2KB 左右，并且在运行过程中会根据需要动态扩容和收缩。
变量存储：类似于 C 和 C++，局部变量、函数参数等存放在栈内存中。
Java
自动管理：由 Java 虚拟机（JVM）自动管理栈内存的分配和释放，程序员无需操心。
大小动态：栈内存的大小会随着线程的执行而动态变化，但一般没有固定的初始大小设定，由 JVM 根据线程的需求进行分配。
方法执行：在方法调用时，会在栈中创建一个栈帧，用于存储方法的局部变量、操作数栈、动态链接、方法返回地址等信息。
## 堆内存
C 和 C++
手动分配释放：需要程序员通过malloc、new等函数手动进行内存分配，使用free、delete等函数手动释放内存，若忘记释放可能导致内存泄漏。
内存碎片：频繁的分配和释放可能导致堆内存中出现碎片，影响内存的使用效率和程序性能。
对象存储：通过new创建的对象或通过malloc分配的动态内存空间通常存放在堆内存中。
Go
自动回收：由 Go 的垃圾回收器（GC）自动管理堆内存的回收，程序员不需要显式地释放堆内存中的对象。
并发友好：Go 的 GC 机制在设计上考虑了并发编程的需求，尽量减少对程序并发性能的影响。
动态分配：通过make和new等函数在堆内存中分配对象或数据结构，如创建切片、映射等。
Java
自动回收：依赖 Java 虚拟机的垃圾回收机制（JVM GC）来自动回收堆内存中的对象，不需要程序员手动释放。
分代回收：JVM 的 GC 通常采用分代回收策略，将堆内存分为年轻代、年老代等不同的区域，针对不同代的对象采用不同的回收方式，以提高回收效率。
对象存储：几乎所有的对象都在堆内存中创建和存储，只有一些基本数据类型的局部变量和方法参数等可能存储在栈内存中。
--- 完全面向对象的语言

## 全局 / 静态存储区
C 和 C++
编译时分配：全局变量和静态变量存放在全局 / 静态存储区，在编译时就分配好了内存，其生命周期与程序的生命周期相同。
共享访问：全局变量可以在整个程序的不同函数和文件中共享访问，静态变量则根据其定义的不同，有静态局部变量和静态全局变量之分，分别在局部或文件范围内共享。
Go
类似但有别：Go 语言中没有直接等同于 C 和 C++ 的全局 / 静态存储区概念，但有全局变量和包级静态变量，它们在程序启动时就会被初始化并分配内存，其生命周期与程序相同，在不同的函数和包之间可以共享。
初始化顺序：按照包的导入顺序依次初始化，在同一个包内，按照变量声明的顺序初始化。
Java
类变量：有类变量的概念，相当于 C 和 C++ 的静态变量，它属于类的所有实例共享，在类加载时就会被分配到内存的特定区域，其生命周期与类的生命周期相同。
类加载：在类加载过程中，类的静态代码块、类变量等会被依次初始化，并且按照它们在类中的定义顺序进行。
## 常量存储区
C 和 C++
只读区域：常量数据存放在常量存储区，在程序运行期间是只读的，不允许修改。如字符串常量、数值常量等通常存放在这里。
内存共享：常量存储区的内存通常是共享的，多个相同的常量在内存中可能只占用一份空间，以提高内存的使用效率。
Go
无显式区分：Go 语言中没有明确单独划分出常量存储区，但常量在内存中的存储方式与 C 和 C++ 类似，也是不可变的，且编译器会尽量优化常量的存储，可能存在一定程度的内存共享。
编译时确定：常量的值在编译时就已经确定，不能在程序运行时随意修改。
Java
运行时常量池：Java 有运行时常量池，用于存放编译期生成的各种字面量和符号引用，如字符串常量、类和接口的全限定名、字段和方法的引用等。
动态扩展：运行时常量池在 Java 虚拟机的方法区中，它在程序运行过程中可以动态扩展，以适应不断生成的新常量。


## 代码区
存储内容
存放程序的可执行代码，也就是 CPU 要执行的机器指令。这些指令是由编译器将源代码编译后生成的，在程序运行期间是只读的，不允许修改。
特点与作用
代码区的内存通常是共享的，多个相同程序的进程可以共享同一份代码区的内存，这样可以节省内存空间。同时，由于其只读的特性，保证了程序代码的稳定性和安全性，防止代码在运行过程中被意外篡改。

## 自由存储区（仅 C++）
存储内容
通过malloc、free等函数动态分配和释放的内存空间，它与堆内存类似，但又有所不同。堆内存是由操作系统管理的一块连续的内存区域，而自由存储区是 C++ 中通过特定的库函数来管理的内存区域，其实现可以基于堆，也可以有其他的实现方式。
特点与作用
自由存储区提供了一种更灵活的动态内存管理方式，与new、delete操作符管理的内存区域有所区别，但目的都是为了满足程序在运行过程中动态分配和释放内存的需求。不过在实际应用中，由于new和delete在 C++ 中更为常用和方便，自由存储区的使用相对较少。

## 栈内存的分配和释放
栈内存的分配和释放主要由编译器（所以c、c++不同的编译器程序的栈内存大小可能不同）自动完成，但程序员也可以通过一些特殊的指令或函数进行一定程度的手动控制。

