package schemas

type StartProxyRequest struct {
	Port   int    `json:"port" binding:"required"`
	Target string `json:"target" binding:"required"`
}

type ProxyStatusResponse struct {
	Status   string `json:"status"`
	Port     int    `json:"port"`
	Target   string `json:"target"`
	ActiveConnections int `json:"active_connections"`
}

type ProxyHealthResponse struct {
	Status        string `json:"status"`
	Uptime        int64  `json:"uptime"`
	TotalSessions int    `json:"total_sessions"`
	TotalMessages int64  `json:"total_messages"`
}
