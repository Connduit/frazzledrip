#ifndef LOGGER_H
#define LOGGER_H

#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <ctime>
#include <iomanip>
#include <cstdarg>   // va_list, va_start, va_end
#include <cstdio>    // std::vsnprintf

#ifdef _DEBUG

// TODO: add thread-safety / multithreaded capability
class Logger
{
public:
    enum class Level {
        Trace,
        Debug,
        Info,
        Warn,
        Error,
        Critical
    };

    static Logger& instance();

    void setOutput(std::ostream* out);
    void setMinLevel(Level level);

    // printf-style variadic function
    void logf(Level level, const char* file, int line,
              const char* fmt, ...);

private:
    Logger();

    // internal implementation that receives final formatted string
    void logImpl(Level level,
                 const std::string& formatted,
                 const char* file,
                 int line);

    std::ostream* output_;
    std::ofstream fileStream_;
    Level minLevel_;

    const char* levelToString(Level level) const;
    std::string currentTimestamp() const;
};

// Macros redirect into logf()
#define LOG_TRACE(fmt, ...) Logger::instance().logf(Logger::Level::Trace, __FILE__, __LINE__, fmt, ##__VA_ARGS__)
#define LOG_DEBUG(fmt, ...) Logger::instance().logf(Logger::Level::Debug, __FILE__, __LINE__, fmt, ##__VA_ARGS__)
#define LOG_INFO(fmt, ...)  Logger::instance().logf(Logger::Level::Info,  __FILE__, __LINE__, fmt, ##__VA_ARGS__)
#define LOG_WARN(fmt, ...)  Logger::instance().logf(Logger::Level::Warn,  __FILE__, __LINE__, fmt, ##__VA_ARGS__)
#define LOG_ERROR(fmt, ...) Logger::instance().logf(Logger::Level::Error, __FILE__, __LINE__, fmt, ##__VA_ARGS__)
#define LOG_CRITICAL(fmt, ...) Logger::instance().logf(Logger::Level::Critical, __FILE__, __LINE__, fmt, ##__VA_ARGS__)

#else

// Release build — all no-ops
class Logger
{
public:
    enum class Level { Trace, Debug, Info, Warn, Error, Critical };
    static Logger& instance() { static Logger dummy; return dummy; }
    void setOutput(std::ostream*) {}
    void setMinLevel(Level) {}
    void logf(Level, const char*, int, const char*, ...) {}
};
#define LOG_TRACE(fmt, ...)
#define LOG_DEBUG(fmt, ...)
#define LOG_INFO(fmt, ...)
#define LOG_WARN(fmt, ...)
#define LOG_ERROR(fmt, ...)
#define LOG_CRITICAL(fmt, ...)

#endif // _DEBUG

#endif // LOGGER_H
