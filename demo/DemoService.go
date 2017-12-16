package demo

type DemoService struct{
	ServiceMeta string `ServiceName:"MyApollo" ServicePort:"8082"`
}

func (this DemoService) Hello(a string) string {
	return "hello word form golang"
}
