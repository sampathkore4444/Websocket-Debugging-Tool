import { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Link } from 'react-router-dom'

interface Session {
  id: number
  connection_id: string
  client_ip: string
  server_host: string
  status: string
  message_count: number
  start_time: string
}

export function Sessions() {
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/sessions')
      .then(res => res.json())
      .then(data => {
        setSessions(data.sessions || [])
        setLoading(false)
      })
      .catch(() => setLoading(false))
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">Sessions</h1>
        <Button>New Session</Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>All Sessions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <Input placeholder="Search sessions..." className="max-w-md" />
            {loading ? (
              <div className="text-center py-4">Loading...</div>
            ) : sessions.length === 0 ? (
              <div className="text-center py-4 text-muted-foreground">
                No sessions yet. Start a new session to begin capturing WebSocket traffic.
              </div>
            ) : (
              <div className="space-y-2">
                {sessions.map(session => (
                  <Link
                    key={session.id}
                    to={`/sessions/${session.id}`}
                    className="flex items-center justify-between p-4 rounded-lg border hover:bg-accent transition-colors"
                  >
                    <div className="space-y-1">
                      <p className="font-medium">{session.connection_id}</p>
                      <p className="text-sm text-muted-foreground">
                        {session.client_ip} → {session.server_host}
                      </p>
                    </div>
                    <div className="text-right">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        session.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                      }`}>
                        {session.status}
                      </span>
                      <p className="text-sm text-muted-foreground mt-1">
                        {session.message_count} messages
                      </p>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
