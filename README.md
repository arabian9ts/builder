# builder

## What is builder?
Code generating tool for initialize entity keep while keeping encapsulation.
In order to lock our business logic into models, we need to use private fields, which makes it difficult for other packages to initialize models.
So, builder generates model builder structs and funcs.

## Installation
```sh
go get -u github.com/arabian9ts/builder
```

## Example
This struct express user entity and is in `./entity` package.
```go
type User struct {
	id        string
	name      string
	digest    string
	timestamp int64
}
```

To generate user builders, specify the target package(s).
```sh
builder entity
```

Then, user builder is generated as following.
```go
func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (userBuilder *UserBuilder) Id(id string) *UserBuilder {
	userBuilder.id = id
	return userBuilder
}

func (userBuilder *UserBuilder) Name(name string) *UserBuilder {
	userBuilder.name = name
	return userBuilder
}

func (userBuilder *UserBuilder) Digest(digest string) *UserBuilder {
	userBuilder.digest = digest
	return userBuilder
}

func (userBuilder *UserBuilder) Timestamp(timestamp int64) *UserBuilder {
	userBuilder.timestamp = timestamp
	return userBuilder
}

func (userBuilder UserBuilder) Build() *User {
	return &User{
		digest:    userBuilder.digest,
		id:        userBuilder.id,
		name:      userBuilder.name,
		timestamp: userBuilder.timestamp,
	}
}
``` 

The Usage of these builder code is ...
```go
user := NewUserBuilder().
    Id("id").
    Name("name").
    Timestamp(time.Now().Unix()).
    Build()
```

## ToDo
- [ ] skip struct tag for ignore generating builder func.
- [ ] getter or setter func with struct tag
