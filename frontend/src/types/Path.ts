export enum TimeFrame {
  m1 = '1m',
  m3 = '3m',
  m5 = '5m',
  m15 = '15m',
  m30 = '30m',
  h1 = '1h',
  h2 = '2h',
  h4 = '4h',
  h12 = '12h',
  d1 = '1d',
  d3 = '3d',
  w1 = '1w',
}

export const timeframes = [
  TimeFrame.m1, TimeFrame.m3, TimeFrame.m5, TimeFrame.m15, TimeFrame.m30,
  TimeFrame.h1, TimeFrame.h2, TimeFrame.h4, TimeFrame.h12, TimeFrame.d1,
  TimeFrame.d3, TimeFrame.w1
]

export enum Operand {
  GREATER_THAN = 'GREATER_THAN',
  LOWER_THAN = 'LOWER_THAN',
  GREATER_THAN_OR_EQUAL = 'GREATER_THAN_OR_EQUAL',
  LOWER_THAN_OR_EQUAL = 'LOWER_THAN_OR_EQUAL',
  EQUAL = 'EQUAL',
  NOT_EQUAL = 'NOT_EQUAL',
  CROSS_ABOVE = 'CROSS_ABOVE',
  CROSS_BELOW = 'CROSS_BELOW'
}

export enum StepType {
  SIGNAL_TILE = 'SIGNAL_TILE',
  ANY_SIGNAL_TILE = 'ANY_SIGNAL_TILE',
  CHUNK_ID = 'CHUNK_ID'
}

export interface IndicatorSettings {
  timeframe?: TimeFrame
  candlesBack: number
  realTime: boolean
  offset: number
  indicatorKey: string
  data: Record<string, any>
  symbol?: string
}

export interface SignalTile {
  id: string
  name: string
  persistance: number
  operand?: Operand
  indicatorA?: IndicatorSettings
  indicatorB?: IndicatorSettings
}

export interface AnySignal {
  signals: SignalTile[]
  amount: number
}

type ChunkId = string

export interface PathStep {
  id: string
  type: StepType
  data: SignalTile | AnySignal | ChunkId
}

export interface Chunk {
  id: string
  name: string
  steps: PathStep[];
}

export interface Path {
  id: string
  name: string | undefined
  whitelist: string[]
  steps: PathStep[]
  type: 'BUY' | 'SELL'
}
