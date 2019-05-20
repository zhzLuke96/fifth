# fif-code
> 本身作为底层语言，forth已经足够了，fif在为基本的forth提供了更多的函数支持和抽象

### Usage
```
Hit CTRL+C or type "exit" or "quit" to quit.
>>> "hello world"
>>> print
hello world
>>>
```

```
>>> "square" func dup mul ret endfunc store
>>> 12 square print
144
```

### Closure
```
"main" func 
	"count1"
	"counter" call
	store
	"count2"
	"counter" call
	store

	"count1" call println
	"count1" call println
	"count2" call println
endfunc store

"counter" func
	"count" 0 store
	func "count" load 1 add dup "count" swap storev ret endfunc
	ret
endfunc store

"main" call
```

### Immediately Invoked Function Expression
```
"main" func
	&end jmp

	func
		&end jmp
		"foo" println
		end: 
		"bar" println
		ret
	endfunc call

	
	"sub" func 
		"hello world" println ret
	endfunc store
	&& jmp
	"sub" call
	ret
endfunc store

"main" call
```

### oop
```
"newMap" func 
	func 
		"val" arg "key" arg "opt" arg
		"opt" load "get" strEqul 0 &set equljmp
		"key" load load ret
		set:
		"opt" load "set" strEqul 0 &end equljmp
		"key" load "val" load stores
		end: ret
	endfunc
	ret 
endfunc store

"main" func
	"map1" newMap store
	"set" "name" "alice" "map1" load call
	"set" "age" 20 "map1" call
	"get" "name" nop "map1" call println
	"get" "age" nop "map1" call println
endfunc store

"main" call
```