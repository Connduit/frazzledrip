#ifndef LOGGER_H
#define LOGGER_H

#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <ctime>
#include <iomanip>

#ifdef _DEBUG   // Entire logger only exists in debug mode

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

    // Singleton accessor
    static Logger& instance();

    // Set output stream (default = std::cout)
    void setOutput(std::ostream* out);

    // Log function
    void log(Level level,
             const std::string& msg,
             const char* file,
             int line);

private:
    Logger();
    std::ostream* output_;
    std::ofstream fileStream_;

    const char* levelToString(Level level) const;
    std::string currentTimestamp() const;
};

// Logging macros
#define LOG_TRACE(msg) Logger::instance().log(Logger::Level::Trace, msg, __FILE__, __LINE__)
#define LOG_DEBUG(msg) Logger::instance().log(Logger::Level::Debug, msg, __FILE__, __LINE__)
#define LOG_INFO(msg)  Logger::instance().log(Logger::Level::Info,  msg, __FILE__, __LINE__)
#define LOG_WARN(msg)  Logger::instance().log(Logger::Level::Warn,  msg, __FILE__, __LINE__)
#define LOG_ERROR(msg) Logger::instance().log(Logger::Level::Error, msg, __FILE__, __LINE__)
#define LOG_CRITICAL(msg) Logger::instance().log(Logger::Level::Critical, msg, __FILE__, __LINE__)

#else   // If NOT _DEBUG — remove all logging

class Logger
{
public:
    enum class Level { Trace, Debug, Info, Warn, Error, Critical };

    static inline Logger& instance() {
        static Logger dummy;
        return dummy;
    }

    void setOutput(std::ostream*) {}
    void log(Level, const std::string&, const char*, int) {}

private:
    Logger() = default;
};

// Define macros to nothing
#define LOG_TRACE(msg)
#define LOG_DEBUG(msg)
#define LOG_INFO(msg)
#define LOG_WARN(msg)
#define LOG_ERROR(msg)
#define LOG_CRITICAL(msg)

#endif // _DEBUG

#endif // LOGGER_H
