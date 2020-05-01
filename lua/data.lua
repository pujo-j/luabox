local luabox = require("luabox")

local function toYaml(data)
    return luabox.yamlRepr(data)
end

local function parseYaml(string)
    return luabox.yamlParse(string)
end

local function toJson(data)
    return luabox.jsonRepr(data)
end

local function parseJson(string)
    return luabox.jsonParse(string)
end

return {
    toYaml = toYaml,
    fromYaml = parseYaml,
    toJson = toJson,
    fromJson = parseJson,
}

