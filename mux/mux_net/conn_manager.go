package mux_net


type ConnManager struct{
	connections []Connection
}

func (cm *ConnManager)AddConnection(connection Connection){
	cm.connections = append(cm.connections, connection)
}

func (cm *ConnManager) RemoveConnection(connectionId int) {
	cm.connections = append(cm.connections[:connectionId], cm.connections[:connectionId+1]...)
}
