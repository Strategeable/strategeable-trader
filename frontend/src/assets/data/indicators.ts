import { Indicator } from '@/types/Indicator'

export type indicatorKey = 'RSI' | 'CANDLE_OPEN' | 'CANDLE_CLOSE' | 'CANDLE_HIGH' | 'CANDLE_LOW' | 'NUMBER'

const indicators: Indicator[] = [
  {
    key: 'RSI',
    name: 'Relative Strength Index',
    shortName: 'RSI',
    hasTimeframe: true,
    hasSource: true,
    fields: [
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
    hasSource: true,
    fields: [
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
    key: 'EMA',
    name: 'Exponential Moving Average',
    shortName: 'EMA',
    hasTimeframe: true,
    hasSource: true,
    fields: [
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
    key: 'CANDLE_POSITION_VALUE',
    name: 'Candle Position Value',
    shortName: 'Candle Value',
    hasTimeframe: true,
    hasSource: false,
    fields: [
      {
        key: 'candlePosition',
        name: 'Candle Position',
        required: true,
        type: 'select',
        advanced: false,
        options: ['CLOSE', 'OPEN', 'HIGH', 'LOW'],
        default: 'CLOSE'
      }
    ]
  },
  {
    key: 'NUMBER',
    name: 'Number',
    shortName: 'Number',
    hasTimeframe: false,
    hasSource: false,
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
  },
  {
    key: 'POSITION_CHANGE',
    name: 'Position Change %',
    shortName: 'Position %',
    hasTimeframe: true,
    hasSource: false,
    fields: []
  },
  {
    key: 'POSITION_HOLD_TIME',
    name: 'Position Hold Time (s)',
    shortName: 'Position hold (s)',
    hasTimeframe: true,
    hasSource: false,
    fields: []
  }
]

export default indicators
