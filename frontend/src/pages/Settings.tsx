import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'

export function Settings() {
  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold tracking-tight">Settings</h1>

      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Proxy Configuration</CardTitle>
            <CardDescription>Configure the WebSocket proxy settings</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Default Target URL</label>
              <Input placeholder="ws://localhost:3000" />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Proxy Port</label>
              <Input type="number" placeholder="8080" />
            </div>
            <Button>Save</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Storage</CardTitle>
            <CardDescription>Configure session storage</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Session Retention (days)</label>
              <Input type="number" placeholder="30" />
            </div>
            <Button>Save</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Security</CardTitle>
            <CardDescription>Configure security settings</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-2">
              <input type="checkbox" id="tls" />
              <label htmlFor="tls" className="text-sm">Enable TLS interception</label>
            </div>
            <div className="flex items-center gap-2">
              <input type="checkbox" id="masking" />
              <label htmlFor="masking" className="text-sm">Enable payload masking</label>
            </div>
            <Button>Save</Button>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
