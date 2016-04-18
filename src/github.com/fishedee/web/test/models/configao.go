package test

type ConfigAoModel struct {
	BaseModel
	data map[string]string
}

func (this *ConfigAoModel) Set(key string, value string) {
	if this.data == nil {
		this.data = map[string]string{}
	}
	this.data[key] = value
}

func (this *ConfigAoModel) Get(key string) string {
	return this.data[key]
}
