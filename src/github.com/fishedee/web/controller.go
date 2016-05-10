package web

type ControllerInterface interface {
	AutoRender(interface{}, string)
}

type Controller struct {
	*Basic
}

func (this *Controller) AutoRender(result interface{}, view string) {

}

func (this *Controller) Check(requireStruct interface{}) {
	this.Ctx.GetParamToStruct(requireStruct)
}

func (this *Controller) CheckGet(requireStruct interface{}) {
	this.Ctx.GetUrlParamToStruct(requireStruct)
}

func (this *Controller) CheckPost(requireStruct interface{}) {
	this.Ctx.GetFormParamToStruct(requireStruct)
}

func (this *Controller) Write(data []byte) {
	this.Ctx.Write(data)
}

func (this *Controller) WriteHeader(key string, value string) {
	this.Ctx.WriteHeader(key, value)
}

func (this *Controller) WriteMimeHeader(mime string, title string) {
	this.Ctx.WriteMimeHeader(mime, title)
}

func InitController(controller interface{}) {
	err := addIocTarget(controller)
	if err != nil {
		panic(err)
	}
}
