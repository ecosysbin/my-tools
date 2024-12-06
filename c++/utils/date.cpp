# include <iostream>
# include <string>
# include <ctime>

using namespace std;

string getCurrentTime() {
    time_t now = time(0);
    char* dt = ctime(&now);
    return string(dt);
}

