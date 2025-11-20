#include "Logger.h"

#ifdef _DEBUG

Logger::Logger()
    : output_(&std::cout),
      minLevel_(Level::Trace)
{
}

Logger& Logger::instance()
{
    static Logger inst;
    return inst;
}

void Logger::setOutput(std::ostream* out)
{
    output_ = out ? out : &std::cout;
}

void Logger::setMinLevel(Level level)
{
    minLevel_ = level;
}

void Logger::logf(Level level,
                  const char* file,
                  int line,
                  const char* fmt,
                  ...)
{
    if (level < minLevel_)
        return;

    char buffer[2048];

    va_list args;
    va_start(args, fmt);
    std::vsnprintf(buffer, sizeof(buffer), fmt, args);
    va_end(args);

    logImpl(level, buffer, file, line);
}

void Logger::logImpl(Level level,
                     const std::string& msg,
                     const char* file,
                     int line)
{
    (*output_) << "[" << currentTimestamp() << "] "
               << "[" << levelToString(level) << "] "
               << file << ":" << line << " - "
               << msg << std::endl;
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

#endif // _DEBUG
