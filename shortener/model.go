package shortener

type Redirect struct {
	Code     string `json: "code" bson: "code" mspack: "code"`
	URL      string `json: "url" bson: "url" mspack: "url" validate: "empty=false & format=url"`
	CreateAt int64  `json: "create_at" bson: "create_at" mspack: "create_at"`
}
