import { Chunk, Path } from './Path'

export interface Strategy {
  id?: string
  version: string
  createdAt: Date
  lastEdited: Date
  name: string
  symbols: string[]
  chunks: Chunk[]
  paths: Path[]
}
