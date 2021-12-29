gotag is a command tool that can automatically generate tags for golang struct.

## Quick start
Install gotag
```
go install github.com/sycki/gotag@latest
```

Process your go file
```
gotag -a json:camel,yaml:camel,gorm:snake model/*.go

or 

gotag -a json:aA,yaml:aA,gorm:a_a model/*.go
```

File content
```
type User struct {
    Name               string `json:"name" yaml:"name" gorm:"name"`
    IdCard             string `json:"idCard" yaml:"idCard" gorm:"id_card"`
    ResidentialAddress string `json:"residentialAddress" yaml:"residentialAddress" gorm:"residential_address"`
    CompanyName        string `json:"companyName" yaml:"companyName" gorm:"company_name"`
}
```

## Usage
```
gotag is a command tool that can automatically generate tags for golang struct.

Usage: gotag [OPTIONS] filename ...

Options:
    -a, --add <TagName:Style>   Add tags to struct
                                TagName:     can be any string
                                Style:       "camel", "snake", "go", "upper", "lower"
                                Style alias: "aA", "a_a", "Aa", "AA", "aa"
    -i, --index <Number>        Specify the position for the --add option
                                Negative number: means cover all
                                0 ~ max integer: insert to the specified position
    -r, --remove <TagName>      Remove tags from struct

Examples:
    gotag -a json:camel -a gorm:snake model/*.go

    gotag -a json:aA,gorm:a_a -i -1 model/*.go

    gotag -r json model/*.go

```
