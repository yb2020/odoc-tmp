import { reactive } from 'vue';

interface Fetch {
  (): void | Promise<void>;
}

export interface FetchState {
  error: Error | null;
  pending: boolean;
  timestamp: number;
}

export default function useFetch(
  requestFn: Fetch,
  immediate?: boolean
): {
  fetch: Fetch;
  fetchState: FetchState;
} {
  const fetchState = reactive<FetchState>({
    error: null,
    pending: false,
    timestamp: Date.now(),
  });

  const fetch = async () => {
    fetchState.pending = true;
    fetchState.error = null;
    fetchState.timestamp = Date.now();
    try {
      await requestFn();
      fetchState.pending = false;
      fetchState.error = null;
    } catch (error) {
      fetchState.error = error as Error;
      fetchState.pending = false;
    }
  };

  if (immediate !== false) {
    fetch();
  }

  return {
    fetch,
    fetchState,
  };
}
