package tube

var (
	hub *HubT
)

func Hub() *HubT {
	if hub == nil {
		hub = new(HubT).Init()
	}
	return hub
}

func (self *HubT) Init() *HubT {
	return self
}

func (self *HubT) UpdateMethods(entrypoint string, methods []string) {
	self.EntryMethods[entrypoint] = methods
}
