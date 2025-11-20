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

    static void log(LogLevel level, const std::string& msg,
                    const char* file, int line)
    {
        instance().write(level, msg, file, line);
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
            std::cerr << "[Logger] Failed to open log file: " << filename << "\n";
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
        std::strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S",
                      std::localtime(&now));
        return buf;
    }

    void write(LogLevel level,
               const std::string& msg,
               const char* file, int line)
    {
        std::string out =
            "[" + timestamp() + "]"
            "[" + levelToString(level) + "] "
            + msg +
            "  (" + file + ":" + std::to_string(line) + ")";

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
// LOG MACROS — now include file + line
// ------------------------------------------------------------
#define LOG_DEBUG(msg)   Logger::log(LogLevel::Debug,   msg, __FILE__, __LINE__)
#define LOG_INFO(msg)    Logger::log(LogLevel::Info,    msg, __FILE__, __LINE__)
#define LOG_WARN(msg)    Logger::log(LogLevel::Warning, msg, __FILE__, __LINE__)
#define LOG_ERROR(msg)   Logger::log(LogLevel::Error,   msg, __FILE__, __LINE__)

#else
// ------------------------------------------------------------
// Release mode — everything compiles to nothing
// ------------------------------------------------------------
enum class LogLevel { Debug, Info, Warning, Error };

class Logger
{
public:
    static void setLogFile(const std::string&) {}
    static void log(LogLevel, const std::string&, const char*, int) {}
};

#define LOG_DEBUG(msg)
#define LOG_INFO(msg)
#define LOG_WARN(msg)
#define LOG_ERROR(msg)

#endif // _DEBUG
