import { MousePointDir, BaseMousePoint } from "../type";

export const getPointDir = (point1: BaseMousePoint, point2: BaseMousePoint) => {
  /**
   * 根据滑动的两点计算出滑动方向是正向还是反向，
   * 这对于选区计算文字长度有影响，需要判断是从头算文字还是从尾算文字
   */

  const { x: x1, y: y1 } = point1
  const { x: x2, y: y2 } = point2

  let angle = 0

  if (x2 === x1) {
    angle = y2 > y1 ? Math.PI / 2 : Math.PI / 2 + Math.PI
  } else {
    const angleTan = (y2 - y1) / (x2 - x1)

    angle = Math.atan(angleTan)
  
    if (angleTan > 0) {
      angle = y2 > y1 ? angle : angle + Math.PI
    } else if (y2 < y1) {
      angle = 2 * Math.PI + angle
    } else {
      angle = angle +  Math.PI
    }
  }
  return angle >= 0 && angle <= Math.PI ? MousePointDir.forward : MousePointDir.backward

};

export const getStartPointAndEndPoint = (point1: BaseMousePoint, point2: BaseMousePoint) => {
  const dir = getPointDir(point1, point2)
  return dir === MousePointDir.forward ? [point1, point2, dir] as const : [point2, point1, dir] as const
};

export const detectCtrlPressing = (event: MouseEvent) => {
  const platform = (navigator as any)?.userAgentData?.platform || navigator?.platform || 'unknown';
  const ctrlPressing = /Mac/i.test(platform) ? event.metaKey : event.ctrlKey;
  return ctrlPressing;
};

export interface MouseDownMoveUpDoubleClickCallbacks {
  onMouseDown?: (e: MouseEvent) => void;
  onMouseMove?: (e: MouseEvent) => void;
  onMouseUp?: (e: MouseEvent) => void;
  onDoubleClick?: (e: MouseEvent) => void;
}

/**
 * 精准区分mousedown、mousemove、mouseup、dblclick
 * 
 * 本函数只在mousedown之后、mouseup之前触发mousemove
 * 
 * 如果希望单独触发mousemove，请不要使用本函数，直接.addEventListener('pointermove', ...)即可
 * 
 * @returns 清除本函数的事件监听
 */
export const useMouseDownMoveUpDoubleClick = (
  element: HTMLElement, 
  {
    onMouseDown,
    onMouseMove,
    onMouseUp,
    onDoubleClick,
  }: MouseDownMoveUpDoubleClickCallbacks,
  {
    time = 300,
    distance = 10,
  }
) => {
  const square = (number: number) => number * number;
  const squareDistance =  square(distance);

  element.addEventListener('pointerdown', mouseDownEntry);
  return () => {
    element.removeEventListener('pointerdown', mouseDownEntry);
  };

  function mouseDownEntry(stage1Event: MouseEvent) {
    const mouseDownEvent = new MouseEvent(stage1Event.type, stage1Event);
  
    const timeout1Id = setTimeout(mouseDown1OutOfTime, time);
    element.addEventListener('pointermove', mouseMove1OutOfRange);
    element.addEventListener('pointerup', mouseUpInTimeWithinRange, { once: true });

    function clearStage1() {
      element.removeEventListener('pointermove', mouseMove1OutOfRange);
      element.removeEventListener('pointerup', mouseUpInTimeWithinRange);
    }

    function faraway ({ pageX, pageY }: MouseEvent) {
      return square(pageX - mouseDownEvent.pageX) 
        + square(pageY - mouseDownEvent.pageY) 
        > squareDistance;
    }
  
    function listenMouseMoveMouseUp() {
      if (onMouseMove) {
        element.addEventListener('pointermove', onMouseMove);
      }

      element.addEventListener('pointerup', (event) => {
        if (onMouseMove) {
          element.removeEventListener('pointermove', onMouseMove);
        }
  
        onMouseUp?.(event);
      }, { once: true });
    }
  
    function mouseDown1OutOfTime () {
      clearStage1();
  
      onMouseDown?.(mouseDownEvent);
  
      listenMouseMoveMouseUp();
    }
  
    function mouseMove1OutOfRange (event: MouseEvent) {
      if (!faraway(event)) {
        return;
      }
  
      clearTimeout(timeout1Id);
      clearStage1();
  
      onMouseDown?.(mouseDownEvent);
      onMouseMove?.(event);
  
      listenMouseMoveMouseUp();
    }
  
    function mouseUpInTimeWithinRange (stage2Event: MouseEvent) {
      clearTimeout(timeout1Id);
      clearStage1();
      element.removeEventListener('pointerdown', mouseDownEntry);
  
      const mouseUpEvent = new MouseEvent(stage2Event.type, stage2Event);
      
      const timeout2Id = setTimeout(
        mouseUpOutOfTime,
        mouseDownEvent.timeStamp - mouseUpEvent.timeStamp + time
      );
      element.addEventListener('pointermove', mouseMove2OutOfRange);
      element.addEventListener('pointerdown', mouseDown2InTime, { once: true });
  
      function clearStage2() {
        element.removeEventListener('pointermove', mouseMove2OutOfRange);
        element.removeEventListener('pointerdown', mouseDown2InTime);
      }
      
      function mouseUpOutOfTime () {
        clearStage2();
  
        onMouseDown?.(mouseDownEvent);
        onMouseUp?.(mouseUpEvent);
  
        element.addEventListener('pointerdown', mouseDownEntry);
      }
  
      function mouseMove2OutOfRange (event: MouseEvent) {
        if (!faraway(event)) {
          return;
        }
  
        clearTimeout(timeout2Id);
        mouseUpOutOfTime();
      }
  
      function mouseDown2InTime () {
        clearTimeout(timeout2Id);
        clearStage2();
  

        if (onDoubleClick) {
          element.addEventListener('pointerup', onDoubleClick, { once: true });
        }
  
        element.addEventListener('pointerdown', mouseDownEntry);
      }
    }
  }
};

export const JS_IGNORE_MOUSE_OUTSIDE = 'js-ignore-mouse-outside'
