import { useParams, Link } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { ArrowLeft, Play, Pause, Trash2 } from 'lucide-react'

export function SessionDetail() {
  const { id } = useParams<{ id: string }>()

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/sessions">
          <Button variant="ghost" size="icon">
            <ArrowLeft className="h-5 w-5" />
          </Button>
        </Link>
        <h1 className="text-3xl font-bold tracking-tight">Session {id}</h1>
      </div>

      <div className="flex gap-2">
        <Button><Play className="mr-2 h-4 w-4" />Replay</Button>
        <Button variant="outline"><Pause className="mr-2 h-4 w-4" />Pause</Button>
        <Button variant="destructive"><Trash2 className="mr-2 h-4 w-4" />Delete</Button>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Session Info</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Connection ID:</span>
              <span className="font-mono text-sm">abc-123-def</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Client IP:</span>
              <span>127.0.0.1</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Server:</span>
              <span>ws://localhost:3000</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Status:</span>
              <span className="text-green-600">Active</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Messages</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-64 overflow-y-auto space-y-2">
              <div className="p-2 bg-blue-50 rounded text-sm">
                <span className="text-blue-600 font-medium">→ Client:</span>
                <pre className="mt-1 text-xs">{"{\"type\": \"login\"}"}</pre>
              </div>
              <div className="p-2 bg-green-50 rounded text-sm">
                <span className="text-green-600 font-medium">← Server:</span>
                <pre className="mt-1 text-xs">{"{\"status\": \"ok\"}"}</pre>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
