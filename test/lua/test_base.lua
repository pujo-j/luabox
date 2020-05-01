print("toto")

local m = require("mod")
m.test()
test2()

local log = require('log')
log.info("Test")

local data = require('data')

print(data.toYaml({
    toto = 'tata'
}))

local t = {}
t[1] = "toto"
t[2] = "titi"
t[3] = "tata"
t[4] = "tutu"
t[5] = "foo"

print(#t)

print(data.toJson(t))
print(data.toYaml(t))