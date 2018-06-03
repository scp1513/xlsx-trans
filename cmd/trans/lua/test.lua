--[=[
[double] id: ID
[double] val: 值
[double] str: 字符串
[double] luaTable: lua测试
[double] luaArr: lua测试2
[double] comb.a: 合并1
[double] comb.b: 合并1
[double] comb.c.a: 深合并
[double] comb.d.a: 深合并2
[double] comb.d.d: 深合并3
[double] comb.d.c: json测试
[double] json2: json测试2
[server] srv: 服务器
[client] cli: 客户端
[double] skip: 跳过左边列
]=]

local data = {
	[1] = {
		id = 1,
		val = 1,
		str = "a",
		luaTable = {
			a = 1,
			b = 2,
		},
		luaArr = {
			1,
			2,
			3,
		},
		json2 = {
			2,
			3,
			4,
		},
		cli = 4,
		skip = 1,
	},
	[2] = {
		id = 2,
		val = 2,
		str = "b",
		luaTable = {
			a = 3,
			b = 4,
		},
		luaArr = {
			1,
			2,
			3,
		},
		json2 = {
			2,
			3,
			4,
		},
		cli = 3,
		skip = 2,
	},
	[3] = {
		id = 3,
		val = -3,
		str = "c",
		luaTable = {
			c = 3,
			d = 4,
		},
		luaArr = {
			[1] = 1,
			[2] = 2,
			[3] = 3,
			a = 3,
		},
		json2 = {
			2,
			3,
			"a",
		},
		cli = 2,
		skip = 3,
	},
	[4] = {
		id = 5,
		val = 0,
		str = "",
		luaTable = nil,
		luaArr = nil,
		json2 = nil,
		cli = 0,
		skip = 0,
	},
}

return data