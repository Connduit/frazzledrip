#include "Logger.h"

#ifdef _DEBUG

Logger::Logger()
    : output_(&std::cout)   // default output is console
{
}

Logger& Logger::instance()
{
    static Logger instance;
    return instance;
}

void Logger::setOutput(std::ostream* out)
{
    if (out)
        output_ = out;
}

const char* Logger::levelToString(Level level) const
{
    switch (level) {
        case Level::Trace:    return "TRACE";
        case Level::Debug:    return "DEBUG";
        case Level::Info:     return "INFO";
        case Level::Warn:     return "WARN";
        case Level::Error:    return "ERROR";
        case Level::Critical: return "CRITICAL";
        default:              return "UNKNOWN";
    }
}

std::string Logger::currentTimestamp() const
{
    auto t = std::time(nullptr);
    std::tm tm{};
#ifdef _WIN32
    localtime_s(&tm, &t);
#else
    localtime_r(&t, &tm);
#endif

    std::ostringstream oss;
    oss << std::put_time(&tm, "%Y-%m-%d %H:%M:%S");
    return oss.str();
}

void Logger::log(Level level,
                 const std::string& msg,
                 const char* file,
                 int line)
{
    (*output_) << "[" << currentTimestamp() << "] "
               << "[" << levelToString(level) << "] "
               << file << ":" << line << " - "
               << msg << std::endl;
}

#endif // _DEBUG
