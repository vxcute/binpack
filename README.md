# binpack

a golang package for packing and unpacking structs. usefull for parsing network protocols and file formats. 

Example: 

    ```go
    package main 

    import "github.com/vxcute/binpack" 

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

        if err := binpack.Unpack(&usr); err != nil {
            log.Fatal(err)
        }

        fmt.Println("User: ", usr)
    }
    ```