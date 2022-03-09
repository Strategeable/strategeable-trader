import { Chunk, Path } from './Path'

export interface Strategy {
  id?: string
  name: string
  symbols: string[]
  chunks: Chunk[]
  paths: Path[]
}
