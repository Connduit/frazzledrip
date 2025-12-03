#ifndef HASHER_H
#define HASHER_H
/* Hasher.h */
/* hashing is used to store passwords securely and verify messages haven't changed */

// TODO: hashes should be a part of a namespace? maybe api namespace? or "precompiled" namespace
#include <cstdint>
#include <string_view>

struct FNV1a32Tag {};
template<typename Tag>
struct Hasher;

template<>
struct Hasher<FNV1a32Tag>
{
    static constexpr uint32_t seed = 0x811C9DC5u;
    static constexpr uint32_t prime = 0x01000193u;

    // Usage: constexpr uint32_t socket_hash = Hasher<FNV1a32Tag>::compute("socket");
    static constexpr uint32_t compute(std::string_view s)
    {
        uint32_t hash = seed;
        for (char c : s)
        {
            hash ^= static_cast<uint8_t>(c);
            hash *= prime;
        }
        return hash;
        
    }
};

#endif
