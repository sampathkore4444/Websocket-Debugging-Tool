import { useState, useEffect } from 'react'
import Editor from '@monaco-editor/react'
import { Button } from './ui/Button'
import { Card, CardContent, CardHeader, CardTitle } from './ui/Card'
import { format } from 'date-fns'

interface Message {
  id: number
  direction: 'client-to-server' | 'server-to-client'
  payload: string
  payload_format: string
  timestamp: string
  is_modified: boolean
  opcode: number
}

interface MessageEditorProps {
  message: Message | null
  onSave: (messageId: number, newPayload: string) => void
  onClose: () => void
}

export function MessageEditor({ message, onSave, onClose }: MessageEditorProps) {
  const [payload, setPayload] = useState('')
  const [isJsonValid, setIsJsonValid] = useState(true)
  const [isDirty, setIsDirty] = useState(false)

  useEffect(() => {
    if (message) {
      setPayload(message.payload)
      setIsDirty(false)
    }
  }, [message])

  const handleEditorChange = (value: string | undefined) => {
    if (value !== undefined) {
      setPayload(value)
      setIsDirty(true)
      
      // Validate JSON
      try {
        JSON.parse(value)
        setIsJsonValid(true)
      } catch {
        setIsJsonValid(false)
      }
    }
  }

  const handleSave = () => {
    if (message && isJsonValid) {
      onSave(message.id, payload)
      setIsDirty(false)
    }
  }

  const handleFormat = () => {
    try {
      const formatted = JSON.stringify(JSON.parse(payload), null, 2)
      setPayload(formatted)
      setIsJsonValid(true)
      setIsDirty(true)
    } catch {
      // Ignore
    }
  }

  const handleRevert = () => {
    if (message) {
      setPayload(message.payload)
      setIsDirty(false)
      setIsJsonValid(true)
    }
  }

  if (!message) {
    return (
      <Card>
        <CardContent className="p-6">
          <div className="text-center text-muted-foreground">
            Select a message to edit
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="h-full">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="text-sm">
            Edit Message #{message.id}
          </CardTitle>
          <div className="flex items-center gap-2">
            <span className={`text-xs px-2 py-1 rounded ${
              message.direction === 'client-to-server' 
                ? 'bg-blue-100 text-blue-800' 
                : 'bg-green-100 text-green-800'
            }`}>
              {message.direction === 'client-to-server' ? '→ Client' : '← Server'}
            </span>
            <span className="text-xs text-muted-foreground">
              {format(new Date(message.timestamp), 'HH:mm:ss')}
            </span>
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {!isJsonValid && (
          <div className="text-sm text-red-500 bg-red-50 p-2 rounded">
            Invalid JSON syntax
          </div>
        )}
        
        <div className="border rounded overflow-hidden" style={{ height: '300px' }}>
          <Editor
            height="100%"
            defaultLanguage="json"
            value={payload}
            onChange={handleEditorChange}
            theme="vs-light"
            options={{
              minimap: { enabled: false },
              fontSize: 13,
              lineNumbers: 'on',
              scrollBeyondLastLine: false,
              automaticLayout: true,
              tabSize: 2,
            }}
          />
        </div>

        <div className="flex items-center justify-between">
          <div className="flex gap-2">
            <Button 
              variant="outline" 
              size="sm" 
              onClick={handleFormat}
            >
              Format
            </Button>
            <Button 
              variant="outline" 
              size="sm" 
              onClick={handleRevert}
              disabled={!isDirty}
            >
              Revert
            </Button>
          </div>
          <div className="flex gap-2">
            <Button 
              variant="ghost" 
              size="sm" 
              onClick={onClose}
            >
              Cancel
            </Button>
            <Button 
              size="sm" 
              onClick={handleSave}
              disabled={!isJsonValid || !isDirty}
            >
              Save Changes
            </Button>
          </div>
        </div>

        {message.is_modified && (
          <div className="text-xs text-amber-600 bg-amber-50 p-2 rounded">
            This message has been modified
          </div>
        )}
      </CardContent>
    </Card>
  )
}
