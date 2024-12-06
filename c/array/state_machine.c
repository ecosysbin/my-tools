#include <stdio.h>

// 定义状态枚举类型
typedef enum
{
    STOPPED,
    MOVING_UP,
    MOVING_DOWN
} ElevatorState;

// 定义事件枚举类型
typedef enum
{
    BUTTON_PRESSED_UP,   // 向上按下按钮
    BUTTON_PRESSED_DOWN, // 向下按下按钮
    ARRIVED_AT_FLOOR     // 到达目标楼层
} ElevatorEvent;

// 状态机结构体，包含当前状态和目标楼层等信息
typedef struct
{
    ElevatorState currentState;
    int targetFloor;
} Elevator;

// 状态转移函数，根据当前状态和事件决定下一个状态
ElevatorState transition(Elevator *elevator, ElevatorEvent event, int floor)
{
    switch (elevator->currentState)
    {
    case STOPPED:
        if (event == BUTTON_PRESSED_UP)
        {
            elevator->targetFloor = floor;
            return MOVING_UP;
        }
        else if (event == BUTTON_PRESSED_DOWN)
        {
            elevator->targetFloor = floor;
            return MOVING_DOWN;
        }
        break;
    case MOVING_UP:
        if (event == ARRIVED_AT_FLOOR && floor == elevator->targetFloor)
        {
            return STOPPED;
        }
        break;
    case MOVING_DOWN:
        if (event == ARRIVED_AT_FLOOR && floor == elevator->targetFloor)
        {
            return STOPPED;
        }
        break;
    }
    return elevator->currentState;
}

int main()
{
    Elevator elevator;
    elevator.currentState = STOPPED;
    elevator.targetFloor = 1;

    // 模拟一些事件触发状态机的运行
    elevator.currentState = transition(&elevator, BUTTON_PRESSED_UP, 5);
    printf("current elevator state: %s\n", (elevator.currentState == MOVING_UP ? "UP" : (elevator.currentState == MOVING_DOWN ? "DOWN" : "STOPPED")));

    elevator.currentState = transition(&elevator, ARRIVED_AT_FLOOR, 5);
    printf("current elevator state: %s\n", (elevator.currentState == MOVING_UP ? "UP" : (elevator.currentState == MOVING_DOWN ? "DOWN" : "STOPPED")));

    return 0;
}