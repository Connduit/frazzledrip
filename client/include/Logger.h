#pragma once

#ifdef _DEBUG

#include <iostream>
#include <fstream>
#include <string>
#include <ctime>

enum class LogLevel
{
    Debug,
    Info,
    Warning,
    Error
};

class Logger
{
public:
    static void setLogFile(const std::string& filename)
    {
        instance().openFile(filename);
    }

    static void log(LogLevel level, const std::string& msg)
    {
        instance().write(level, msg);
    }

private:
    std::ofstream file_;
    bool fileEnabled_ = false;

    static Logger& instance()
    {
        static Logger inst;
        return inst;
    }

    void openFile(const std::string& filename)
    {
        file_.open(filename, std::ios::out | std::ios::app);
        fileEnabled_ = file_.is_open();

        if (!fileEnabled_)
            std::cerr << "[Logger] Failed to open file: " << filename << "\n";
    }

    const char* levelToString(LogLevel lvl)
    {
        switch (lvl)
        {
            case LogLevel::Debug:   return "DEBUG";
            case LogLevel::Info:    return "INFO";
            case LogLevel::Warning: return "WARN";
            case LogLevel::Error:   return "ERROR";
            default:                return "UNKNOWN";
        }
    }

    std::string timestamp()
    {
        std::time_t now = std::time(nullptr);
        char buf[32];
        std::strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S", std::localtime(&now));
        return buf;
    }

    void write(LogLevel level, const std::string& msg)
    {
        std::string out = "[" + timestamp() + "][" + levelToString(level) + "] " + msg;

        if (fileEnabled_)
        {
            file_ << out << "\n";
            file_.flush();
        }
        else
        {
            std::cout << out << std::endl;
        }
    }
};

// ------------------------------------------------------------
// Convenience Macros
// ------------------------------------------------------------
#define LOG_DEBUG(msg)   Logger::log(LogLevel::Debug,   msg)
#define LOG_INFO(msg)    Logger::log(LogLevel::Info,    msg)
#define LOG_WARN(msg)    Logger::log(LogLevel::Warning, msg)
#define LOG_ERROR(msg)   Logger::log(LogLevel::Error,   msg)

#else
// ------------------------------------------------------------
// Release mode: logger does nothing
// ------------------------------------------------------------
enum class LogLevel { Debug, Info, Warning, Error };

class Logger
{
public:
    static void setLogFile(const std::string&) {}
    static void log(LogLevel, const std::string&) {}
};

#define LOG_DEBUG(msg)
#define LOG_INFO(msg)
#define LOG_WARN(msg)
#define LOG_ERROR(msg)

#endif  // _DEBUG
