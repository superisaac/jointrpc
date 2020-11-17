package tube

func NewConn() *ConnT {
	connId := NextCID()
	ch := make(MsgChannel, 100)
	methods := make(map[string]bool)
	conn := &ConnT{ConnId: connId, RecvChannel: ch, Methods: methods}
	return conn
}

func (self ConnT) GetMethods() []string {
	var keys []string
	for k := range self.Methods {
		keys = append(keys, k)
	}
	return keys
}
