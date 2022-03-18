import { StoreType } from '../store'

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $store: StoreType
  }
}
