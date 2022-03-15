import { Chunk, Path } from './Path'

export interface Variable {
  type: 'number' | 'timeframe'
  id: string
  key: string
  value: any
}

export interface Strategy {
  id?: string
  version: string
  exchange: string
  createdAt: Date
  lastEdited: Date
  name: string
  symbols: string[]
  chunks: Chunk[]
  paths: Path[]
  variables: Variable[]
}
