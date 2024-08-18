package contextx

type Key struct {
	name string
}

func (k Key) String() string {
	return k.name
}

var UserIDKey = Key{
	name: "uid",
}
