import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'
import { Button } from '../components/ui/Button'

export function FuzzTesting() {
  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold tracking-tight">Fuzz Testing</h1>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Create Fuzz Test</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Test Name</label>
              <input className="w-full p-2 border rounded" placeholder="Enter test name" />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Template JSON</label>
              <textarea 
                className="w-full h-32 p-2 border rounded font-mono text-sm" 
                placeholder='{"type": "test"}'
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Strategy</label>
              <select className="w-full p-2 border rounded">
                <option value="random">Random</option>
                <option value="mutation">Mutation</option>
                <option value="boundary">Boundary</option>
                <option value="invalid">Invalid</option>
              </select>
            </div>
            <Button className="w-full">Start Fuzz Test</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Previous Tests</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-center py-8 text-muted-foreground">
              No fuzz tests yet.
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
