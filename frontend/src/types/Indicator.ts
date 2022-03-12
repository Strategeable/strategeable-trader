export interface IndicatorField {
  name: string
  key: string
  type: 'text' | 'number' | 'radio' | 'checkbox' | 'select' | 'signal'
  required: boolean
  advanced: boolean
  options: string[]
  default: string
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value?: any
  placeholder?: string
  max?: number
  min?: number
}

export interface Indicator {
  key: string
  name: string
  shortName: string
  hasTimeframe: boolean
  hasSource: boolean
  fields: IndicatorField[]
}
