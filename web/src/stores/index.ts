import { createPinia } from "pinia";

const store = createPinia();

export function usePinaStore() {
  return store;
}