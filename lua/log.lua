local luabox = require("luabox")

local function debug(message, context_table)
    luabox.log(1, message, context_table)
end

local function info(message, context_table)
    luabox.log(2, message, context_table)
end


local function warn(message, context_table)
    luabox.log(3, message, context_table)
end


local function error(message, context_table)
    luabox.log(4, message, context_table)
end


local function fatal(message, context_table)
    luabox.log(5, message, context_table)
end


return {
    debug=debug,
    info=info,
    warn=warn,
    error=error,
    fatal=fatal,
}