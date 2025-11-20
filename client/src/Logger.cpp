#include "Logger.h"
Logger& Logger::Instance() {
    static Logger instance;
    return instance;
}

Logger::Logger() : out_(&std::cout) {}

void Logger::SetOutputFile(const std::string& filename) {
    file_.open(filename);
    if (file_.is_open()) {
        out_ = &file_;
    }
}

void Logger::Log(LogLevel level, const std::string& msg,
                 const char* file, int line)
{
    (*out_) << "[" << file << ":" << line << "] "
            << msg << std::endl;
}
