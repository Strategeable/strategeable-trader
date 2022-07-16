export default class EventEmitter {
  public _events: any

  constructor () {
    this._events = {}
  }

  on (name: string, listener: any) {
    if (!this._events[name]) {
      this._events[name] = []
    }

    this._events[name].push(listener)
  }

  emit (name: string, data: any) {
    if (!this._events[name]) {
      throw new Error(`Can't emit an event. Event "${name}" doesn't exits.`)
    }

    const fireCallbacks = (callback: any) => {
      callback(data)
    }

    this._events[name].forEach(fireCallbacks)
  }
}
