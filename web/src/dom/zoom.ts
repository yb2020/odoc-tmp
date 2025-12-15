const WHEEL_ZOOM_DISABLED_TIMEOUT = 1000; // ms

export interface ZoomFn {
  (ticks: number, evt: WheelEvent): void
}

export interface ValidFn {
  (evt: WheelEvent): boolean
}

function normalizeWheelEventDirection(evt: WheelEvent) {
  let delta = Math.hypot(evt.deltaX, evt.deltaY);
  const angle = Math.atan2(evt.deltaY, evt.deltaX);

  if (-0.25 * Math.PI < angle && angle < 0.75 * Math.PI) {
    delta = -delta;
  }

  return delta;
}

export default class ZoomController {
  private zoomDisabledTimeout = 0

  private _wheelUnusedTicks = 0

  private isValid

  private zoomIn
  private zoomOut

  constructor({ isValid, zoomIn, zoomOut }: { isValid: ValidFn, zoomIn: ZoomFn, zoomOut: ZoomFn }) {
    this.isValid = isValid
    this.zoomIn = zoomIn
    this.zoomOut = zoomOut
  }

  private setZoomDisabledTimeout() {
    if (this.zoomDisabledTimeout) {
      clearTimeout(this.zoomDisabledTimeout);
    }
    this.zoomDisabledTimeout = window.setTimeout(() => {
      this.zoomDisabledTimeout = 0;
    }, WHEEL_ZOOM_DISABLED_TIMEOUT);
  }

  private accumulateWheelTicks(ticks: number) {
    if (
      (this._wheelUnusedTicks > 0 && ticks < 0) ||
      (this._wheelUnusedTicks < 0 && ticks > 0)
    ) {
      this._wheelUnusedTicks = 0;
    }
    this._wheelUnusedTicks += ticks;
    const wholeTicks =
      Math.sign(this._wheelUnusedTicks) *
      Math.floor(Math.abs(this._wheelUnusedTicks));
    this._wheelUnusedTicks -= wholeTicks;
    return wholeTicks;
  }

  onWheel(evt: WheelEvent) {

    if (evt.ctrlKey || evt.metaKey) {
      if (!this.isValid(evt)) {
        return
      }
      // Only zoom the pages, not the entire viewer.
      evt.preventDefault();
      // NOTE: this check must be placed *after* preventDefault.
      if (this.zoomDisabledTimeout || document.visibilityState === 'hidden') {
        return;
      }

      const delta = normalizeWheelEventDirection(evt);
      let ticks = 0;
      if (
        evt.deltaMode === WheelEvent.DOM_DELTA_LINE ||
        evt.deltaMode === WheelEvent.DOM_DELTA_PAGE
      ) {
        // For line-based devices, use one tick per event, because different
        // OSs have different defaults for the number lines. But we generally
        // want one "clicky" roll of the wheel (which produces one event) to
        // adjust the zoom by one step.
        if (Math.abs(delta) >= 1) {
          ticks = Math.sign(delta);
        } else {
          // If we're getting fractional lines (I can't think of a scenario
          // this might actually happen), be safe and use the accumulator.
          ticks = this.accumulateWheelTicks(delta);
        }
      } else {
        // pixel-based devices
        const PIXELS_PER_LINE_SCALE = 150;
        ticks = this.accumulateWheelTicks(delta / PIXELS_PER_LINE_SCALE);
      }

      if (ticks < 0) {
        this.zoomOut(-ticks, evt);
      } else if (ticks > 0) {
        this.zoomIn(ticks, evt);
      }

    } else {
      this.setZoomDisabledTimeout();
    }
  }

  bindWheelToZoom() {

    /**
     * 这里使用perfectScroll之后有个问题，在缩放的时候无法禁止掉ps的滚动，所以就改了ps的源码，做了简单的处理
     * 不知道这里有没有更好的解法方案
     */
    document.addEventListener('wheel', this.onWheel, {
      passive: false,
    });
  }

  destroy() {
    document.removeEventListener('wheel', this.onWheel)
  }

}