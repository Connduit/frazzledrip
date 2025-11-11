#ifndef DEBUG_H
#define DEBUG_H

#include <stdio.h>

#ifdef _DEBUG
    #define DEBUG_PRINT(fmt, ...) \
        fprintf(stderr, "[DEBUG] " fmt "\n", ##__VA_ARGS__)
#else
    #define DEBUG_PRINT(fmt, ...) \
        do { } while (0)
#endif

#endif // DEBUG_H
