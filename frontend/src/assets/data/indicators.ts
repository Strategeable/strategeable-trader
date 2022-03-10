import { Indicator } from '@/types/Indicator'

export type indicatorKey = 'RSI' | 'CANDLE_OPEN' | 'CANDLE_CLOSE' | 'CANDLE_HIGH' | 'CANDLE_LOW' | 'NUMBER'

const indicators: Indicator[] = [
  {
    key: 'RSI',
    name: 'Relative Strength Index',
    shortName: 'RSI',
    hasTimeframe: true,
    fields: [
      {
        key: 'source',
        name: 'Source',
        required: true,
        type: 'signal',
        advanced: false,
        options: ['CANDLE_OPEN', 'CANDLE_CLOSE', 'CANDLE_HIGH', 'CANDLE_LOW'],
        default: 'CANDLE_CLOSE'
      },
      {
        key: 'period',
        name: 'Period',
        required: true,
        type: 'number',
        advanced: false,
        options: [],
        default: '14',
        placeholder: '14',
        max: 100,
        min: 1
      }
    ]
  },
  {
    key: 'SMA',
    name: 'Simple Moving Average',
    shortName: 'SMA',
    hasTimeframe: true,
    fields: [
      {
        key: 'source',
        name: 'Source',
        required: true,
        type: 'signal',
        advanced: false,
        options: [],
        default: 'CANDLE_CLOSE'
      },
      {
        key: 'period',
        name: 'Period',
        required: true,
        type: 'number',
        advanced: false,
        options: [],
        default: '21',
        placeholder: '21',
        max: 1000,
        min: 2
      }
    ]
  },
  {
    key: 'CANDLE_OPEN',
    name: 'Candle Open',
    shortName: 'Candle Open',
    hasTimeframe: true,
    fields: []
  },
  {
    key: 'CANDLE_CLOSE',
    name: 'Candle Close',
    shortName: 'Candle Close',
    hasTimeframe: true,
    fields: []
  },
  {
    key: 'CANDLE_HIGH',
    name: 'Candle High',
    shortName: 'Candle High',
    hasTimeframe: true,
    fields: []
  },
  {
    key: 'CANDLE_LOW',
    name: 'Candle Low',
    shortName: 'Candle Low',
    hasTimeframe: true,
    fields: []
  },
  {
    key: 'NUMBER',
    name: 'Number',
    shortName: 'Number',
    hasTimeframe: false,
    fields: [
      {
        key: 'number',
        name: 'Number',
        advanced: false,
        default: '0',
        type: 'number',
        required: true,
        options: [],
        placeholder: '0'
      }
    ]
  }
]

export default indicators
