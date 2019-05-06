local key_meta = KEYS[1]
local key_data = KEYS[2]

local credit = tonumber(ARGV[1])
local windowLength = tonumber(ARGV[2])
local bucketLength = tonumber(ARGV[3])
local bestEffort = tonumber(ARGV[4])
local token = tonumber(ARGV[5])
-- local timestamp = tonumber(ARGV[6])
local deduplicationid = ARGV[7]

-- Redis version ≥ 3.2
redis.replicate_commands()
local now = redis.call('TIME')
local utime = tonumber(now[1])
local utimeMicro = tonumber(now[2])
local timestamp = (utime * 1e6 + utimeMicro) * 1e3
local udur = math.floor(windowLength / 1e9)
local slot = math.floor(utime / udur)

-- lookup previous response for the deduplicationid and returns if it is still valid
if (deduplicationid or '') ~= '' then
    local previous_token = tonumber(redis.call("HGET", deduplicationid .. "-" .. key_meta, "token"))
    local previous_expire = tonumber(redis.call("HGET", deduplicationid .. "-" .. key_meta, "expire"))

    if previous_token and previous_expire then
        if timestamp < previous_expire then
            return { previous_token, previous_expire - timestamp }
        end
    end
end

local count = tonumber(redis.call("INCRBY", key_meta .. "-" .. slot, token))
redis.call("PEXPIRE", key_meta .. "-" .. slot, windowLength / 1e6)

-- TODO 过期时间固定为windowLength
local expire = timestamp + windowLength

if count <= credit then
    -- save current request and set expiration time for auto cleanup
    if (deduplicationid or '') ~= '' then
        redis.call("HMSET", deduplicationid .. "-" .. key_meta, "token", token, "expire", expire)
        redis.call("PEXPIRE", deduplicationid .. "-" .. key_meta, math.floor((expire - timestamp) / 1e6))
    end

    return { token, expire - timestamp }
else
    if bestEffort == 1 then
        if count - token < credit then
            -- return maximum available allocated token
            local remaining = token - (count - credit)
            -- save current request and set expiration time for auto cleanup
            if (deduplicationid or '') ~= '' then
                redis.call("HMSET", deduplicationid .. "-" .. key_meta, "token", remaining, "expire", expire)
                redis.call("PEXPIRE", deduplicationid .. "-" .. key_meta, math.floor((expire - timestamp) / 1e6))
            end

            return { remaining, expire - timestamp }
        else
            -- not enough available credit
            return { 0, 0 }
        end
    else
        -- not enough available credit
        return { 0, 0 }
    end
end