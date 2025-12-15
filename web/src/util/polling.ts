interface PollingArgs<P, R> {
  fn: (p: P) => Promise<R>;
  maxAttempts: number;
  interval: number;
  validate: (r: R) => boolean;
  params: P;
}

export function polling<P, R>({
  fn,
  validate,
  interval,
  maxAttempts,
  params,
}: PollingArgs<P, R>) {
  let attempts = 0;

  let abort = false;

  const cancel = () => {
    abort = true;
  };

  const executePoll = async (
    resolve: (value: R) => void,
    reject: (reason?: Error) => void
  ) => {
    if (abort) {
      return reject(new Error('Abort'));
    }
    try {
      const result = await fn(params);
      if (validate(result)) {
        return resolve(result);
      }
    } catch (error) {
      attempts++;
    }
    if (maxAttempts && attempts === maxAttempts) {
      return reject(new Error('Exceeded max attempts'));
    } else {
      window.setTimeout(executePoll, interval, resolve, reject);
    }
  };

  return {
    request: new Promise<R>(executePoll),
    cancel,
  };
}
