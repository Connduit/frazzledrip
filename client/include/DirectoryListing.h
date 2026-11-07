#ifndef DIRECTORY_LISTING_H
#define DIRECTORY_LISTING_H

#include <filesystem>
#include <sstream>
#include <string>
#include <system_error>
#include <stdexcept>

namespace directory_listing {
inline std::string jsonString(const std::string& value) {
    std::ostringstream out;
    out << '"';
    static const char hex[] = "0123456789abcdef";
    for (unsigned char c : value) {
        if (c == '"' || c == '\\') out << '\\' << static_cast<char>(c);
        else if (c < 0x20) out << "\\u00" << hex[c >> 4] << hex[c & 15];
        else out << static_cast<char>(c);
    }
    out << '"';
    return out.str();
}

// Nonrecursive listing: no file contents are opened and no filesystem changes are made.
inline std::string list(const std::string& requested) {
    namespace fs = std::filesystem;
    if (requested.size() > 2048 || requested.find('\0') != std::string::npos)
        throw std::runtime_error("Invalid directory path");
    const auto path = fs::absolute(requested.empty() ? fs::current_path() : fs::u8path(requested)).lexically_normal();
    if (!fs::is_directory(path)) throw std::runtime_error("The path is not a directory");
    // Constructing the iterator throws on access denial instead of reporting an empty folder.
    fs::directory_iterator iterator(path);
    std::ostringstream entries;
    std::size_t count = 0;
    bool truncated = false;
    for (const auto& entry : iterator) {
        if (count == 2000 || entries.tellp() >= std::streampos(8 * 1024 * 1024)) { truncated = true; break; }
        if (count++) entries << ',';
        std::error_code error;
        const bool directory = entry.is_directory(error);
        const bool link = entry.is_symlink(error);
        entries << "{\"name\":" << jsonString(entry.path().filename().u8string())
            << ",\"path\":" << jsonString(entry.path().u8string())
            << ",\"is_directory\":" << (directory ? "true" : "false")
            << ",\"is_link\":" << (link ? "true" : "false");
        if (!directory && entry.is_regular_file(error)) {
            const auto size = entry.file_size(error);
            if (!error) entries << ",\"size\":" << size;
        }
        entries << '}';
    }
    return "{\"path\":" + jsonString(path.u8string()) + ",\"parent\":" + jsonString(path.parent_path().u8string())
        + ",\"truncated\":" + (truncated ? "true" : "false") + ",\"entries\":[" + entries.str() + "]}";
}
}
#endif
