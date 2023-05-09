# binpack

a golang package for packing and unpacking structs. usefull for parsing network protocols and file formats, its still in development and has some limitations, for examples 
it doesn't support fields of struct types and pointer types but I will add them soon.

Example: 

```go

    package main 

    import ( 
        "github.com/vxcute/binpack" 
        "log" 
        "fmt"
    )

    type User struct {
        Name string 
        Age int32 
        Gender string
    }

    func main() {

        user := &User {
            Name: "Ahmed", 
            Age: 18, 
            Gender: "Male",
        }

        b, err := binpack.Pack(user) 

        if err != nil {
            log.Fatal(err)
        }

        fmt.Println("User bytes: ", b)

        var usr User 

        if err := binpack.Unpack(b, &usr); err != nil {
            log.Fatal(err)
        }
        
        fmt.Println("User: ", user)
        fmt.Println("User: ", usr)
    }

```