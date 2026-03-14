import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'

export function Connections() {
  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold tracking-tight">Connections</h1>

      <Card>
        <CardHeader>
          <CardTitle>Active Connections</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8 text-muted-foreground">
            No active connections. Start a session to capture WebSocket traffic.
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
