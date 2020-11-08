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
	// copy elements here
	tmp := methods
	self.EntryMethods[entrypoint] = tmp
}

func (self *HubT) Subscribe(listener MethodDeclChan) {
	self.Listeners = append(self.Listeners, listener)
}

func (self *HubT) Unsubscribe(listener MethodDeclChan) {
	self.Listeners = delete(self.Listeners, listener)
}
