# fifth
🛸Coding for code.

# Program <==> Data
> Program Is Data, Program Operating Data

# build
```
bash build_parser.sh && bash build.sh
```

***requirement***
- [goyacc](https://godoc.org/golang.org/x/tools/cmd/goyacc)

# Usage
```go
func print(text){
    __fifcode__ {
        'text' load print
    }
    return
}

func sum(a,b){
    return a+b 
}

func main(){
    a = 10
    b = -8.5
    res = sum(a,b)
    print(res)
}

main()
```

# generator
```go
gen counter1(num){
    while(1){
        getv = yield num
        
        if getv == null {
            num = num + 1
        } else {
            num = num + getv
        }
    }
}

gen counter2(){
    num = 0
    while true{
        yield num
        num = num + 1
    }
}
```

# Functional
```go
func repeat(fn){
    return func(num){
        return fn(fn(num))
    }
}

func twice(num){
    return num * 2
}

twice2 = repeat(twice)

print(twice2(10))
print(repeat(twice)(10))
```

## Y combinator
```go
var Y,F
Y = func(f){
    (func(x) {
        x(x)
    })(func(x){
        f(func(y){
            x(x)(y)
        })
    })
}


F = func(g){
    func(n){
        if n == 0 {
            return 1
        } else {
            return n * g(n-1)
        }
    }
}

FACT = Y(F)
FACT(5) // =>  120
```

# code gen (macro)
> come soon
```go
macro _if{
    ($test?$texpr:$fexpr) => {
        if $test {
            return literal!($texpr)
        } else {
            return literal!($fexpr)
        }
    }
}
a = 1
b = 0
// codegen call
_if{a>b?print(a):print(b)}
// => FUNCTION [print(a)]

// macro call
_if!{a>b?print(a):print(b)}
// => 1
```

> more macro, not more always better.
> <br>but, more sweet.

# LICENSE
GPL-3.0