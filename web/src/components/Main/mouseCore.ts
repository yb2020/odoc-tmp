import {
  PointOnPdf,
  detectCtrlPressing,
  useMouseDownMoveUpDoubleClick,
} from '@idea/pdf-annotate-viewer';
// import { DocumentViewer } from '@idea/pdf-annotate-viewer/typing/DocumentViewer';
import { DocumentViewer } from '~/src/pdf-annotate-viewer/DocumentViewer';
import {
  renderOldFrame,
  clearModelFrame,
  getModelByPoint,
  Canvas,
  Rect,
  Ellipse,
  colorMap,
  initShape,
  ShapeCallback,
  TextAnnotation,
  createTextDiv,
  Line,
  ARROWHEAD_CLASSNAME,
  renderSinglePageHandwrite,
  prepareTextsShapes,
  PageFcanvasMap,
  setTextRectByAnnotation,
  setEditingAndEditTextarea,
  initTextRect,
  getTextRectByAnnotation,
  deleteTextRectByAnnotation,
  cancelEditText,
  TextCallback,
  createShapeAtPoint,
  updateShapeWithPoint,
  getShapeBound,
  onMouseDownStartEditText,
  useStartEditText,
  getArrowHead,
  getShapeScale,
  SHAPE_DELETE,
  getIpadUnit,
  PDFJSAnnotate,
  NEW_ARROW_ID,
  getTextDivByAnnotation,
} from '@idea/pdf-annotate-core';

import { noteBuffer, useAnnotationStore } from '~/src/stores/annotationStore';
import { WebDrawItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import {
  AnnotationColor,
  IDEAAnnotateType,
  ShapeAnnotation,
  ShapeType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { computed, createApp, ref, watch } from 'vue';
import { PersonAnnotationController } from '~/src/stores/annotationStore/PersonAnnotationController';
import {
  DEFAULT_FONT_SIZE,
  useCommentGlobalState,
} from '~/src/hooks/useNoteState';
import ShapeColor from '~/src/components/ToolTip/ShapeColor.vue';
import { getNewAnnotation } from '../Right/TabPanel/Note/annotation-state';
import { deleteHandwrite } from '~/src/api/annotations';
import { currentNoteInfo, isOwner } from '~/src/store';
import { Modal } from 'ant-design-vue';
import hotkeys from 'hotkeys-js';
import { usePdfStore } from '~/src/stores/pdfStore';
// import PDFPageView from '@idea/pdf-annotate-viewer/typing/PDFPageView';
import PDFPageView from '~/src/pdf-annotate-viewer/PDFPageView';

hotkeys('delete', () => {
  document.querySelector<HTMLDivElement>('.' + SHAPE_DELETE)?.click();
});

const getHandwriteModel = async (
  point: PointOnPdf,
  pdfAnnotater?: PDFJSAnnotate
) => {
  const canvas = pdfAnnotater?.canvasElements.get(point.cur);

  if (!canvas) {
    return null;
  }

  const pageNumberPlus1 = point.cur + 1;
  const modelList: WebDrawItem[] = [];

  (await noteBuffer.handwriteBuffer)?.[pageNumberPlus1]?.forEach((item) => {
    modelList.push(...item.webDrawItems);
  });

  const model = getModelByPoint(
    canvas,
    point.left / point.pv.scale,
    point.top / point.pv.scale,
    modelList,
    getIpadUnit(pdfAnnotater!, point.cur)
  );
  return model;
};

export const shapeState = {
  creating: ref<ShapeType | 'text'>(ShapeType.UNRECOGNIZED),
  startPoint: null as null | PointOnPdf,
  shape: null as null | Rect | Ellipse | Line,
  deleteTitle: '',
  deleteOk: '',
};

export const cancelShape = () => {
  if (shapeState.shape && shapeState.startPoint) {
    const fcanvas = PageFcanvasMap.get(shapeState.startPoint.cur + 1) as Canvas;

    if (shapeState.shape instanceof Line) {
      const arrowhead = getArrowHead(NEW_ARROW_ID);
      arrowhead.remove();
    }

    fcanvas.remove(shapeState.shape);
  }

  shapeState.creating.value = ShapeType.UNRECOGNIZED;
  shapeState.startPoint = null;
  shapeState.shape = null;

  Array.from(PageFcanvasMap.values()).forEach((fc) => {
    fc.discardActiveObject();
    fc.requestRenderAll();
  });
};

watch(shapeState.creating, () => {
  const pdfStore = usePdfStore();
  const pdfViewerRef = computed(() => {
    return pdfStore.getViewer(currentNoteInfo.value?.pdfId);
  });
  const pdfAnnotater = computed(() => {
    return pdfStore.getAnnotater(currentNoteInfo.value?.noteId);
  });
  if (!pdfViewerRef.value || !pdfAnnotater.value) {
    return;
  }
  const viewer = pdfViewerRef.value.getDocumentViewer().getPdfViewer();
  const pages = (viewer._pages ?? []) as PDFPageView[];
  pages.forEach((pageView) => {
    const { pageNumber } = pageView.pdfPage;
    if (pageView.div.getAttribute('data-loaded') === 'true') {
      prepareTextsShapes({
        pageNumber,
        pageView,
        instance: pdfAnnotater.value!,
      });
    }
  });

  const defaultCursor: CSSStyleDeclaration['cursor'] =
    shapeState.creating.value === ShapeType.UNRECOGNIZED
      ? 'default'
      : 'crosshair';

  Array.from(PageFcanvasMap.values()).forEach((fcanvas) => {
    fcanvas.defaultCursor = defaultCursor;
  });
});

export const useMouse = (documentViewer: DocumentViewer) => {
  const localStorageColor = useCommentGlobalState();
  const annotationStore = useAnnotationStore();
  const textCallback = useTextCallback(annotationStore);
  const pdfStore = usePdfStore();
  const pdfAnnotaterRef = computed(() =>
    pdfStore.getAnnotater(currentNoteInfo.value?.noteId)
  );

  let latestMouseDown: Promise<Parameters<
    DocumentViewer['onMouseDown']
  > | null> = (async () => null)();
  const listOfMouseMoveEvent: Array<PointOnPdf | null> = [];
  const clearMouseEvent = () => {
    latestMouseDown = (async () => null)();
    listOfMouseMoveEvent.length = 0;
  };

  useMouseDownMoveUpDoubleClick(
    documentViewer.container,
    {
      onMouseDown: (event) => {
        clearAllTextRect(documentViewer, textCallback);
        clearModelFrame();

        const point = documentViewer.getPoint(event);
        const ctrlPressing = detectCtrlPressing(event);

        latestMouseDown = (async () => {
          if (point && shapeState.creating.value !== ShapeType.UNRECOGNIZED) {
            shapeState.startPoint = point;
            shapeState.shape = createShapeAtPoint(
              point,
              shapeState.creating.value,
              localStorageColor.value.shapeStyleId || AnnotationColor.blue,
              getScale(documentViewer, point.cur),
              documentViewer.container
            );

            return null;
          }

          if (
            !isOwner.value ||
            ctrlPressing ||
            !point ||
            !(await getHandwriteModel(point, pdfAnnotaterRef.value))
          ) {
            documentViewer.onMouseDown(point, ctrlPressing);
            return null;
          }

          return [point, ctrlPressing];
        })();
      },

      onMouseMove: async (event) => {
        const point = documentViewer.getPoint(event);

        if (point && shapeState.creating.value !== ShapeType.UNRECOGNIZED) {
          if (point.cur !== shapeState.startPoint?.cur) {
            return;
          }

          const shape = shapeState.shape as Rect | Ellipse | Line;
          const scale = getScale(documentViewer, point.cur);
          updateShapeWithPoint(
            shapeState.startPoint as PointOnPdf,
            point,
            shape,
            scale
          );
          return;
        }

        const mouseDown = await latestMouseDown;

        if (!mouseDown?.[0]) {
          documentViewer.onMouseMove(point);
          return;
        }

        if (point) {
          const dx = point.left - mouseDown[0].left;
          const dy = point.top - mouseDown[0].top;
          const out = dx * dx + dy * dy > 100;

          if (!out) {
            listOfMouseMoveEvent.push(point);
            return;
          }
        }

        documentViewer.onMouseDown(...mouseDown);
        listOfMouseMoveEvent.forEach(documentViewer.onMouseMove);
        documentViewer.onMouseMove(point);

        clearMouseEvent();
      },

      onMouseUp: async (event) => {
        const point = documentViewer.getPoint(event);
        const shape = shapeState.shape as Rect | Ellipse | Line;

        const localStorageColor = useCommentGlobalState();
        const scale = getScale(documentViewer, shapeState.startPoint?.cur ?? 0);

        if (shapeState.creating.value !== ShapeType.UNRECOGNIZED) {
          if (!point || point.cur !== shapeState.startPoint?.cur) {
            cancelShape();
            return;
          }

          const fcanvas = PageFcanvasMap.get(
            shapeState.startPoint!.cur + 1
          ) as Canvas;
          const bound = getShapeBound(
            shapeState.startPoint as PointOnPdf,
            point
          );
          let { width, height, left, top } = bound;
          if (shape instanceof Rect) {
            if (shapeState.creating.value === 'text') {
              const ipadUnit = getIpadUnit(pdfAnnotaterRef.value!, point.cur);
              const fontSize =
                (localStorageColor.value.shapeFontSize || DEFAULT_FONT_SIZE) *
                ipadUnit *
                scale;
              const MIN_HEIGHT = fontSize * 2;
              if (height < MIN_HEIGHT) {
                height = MIN_HEIGHT;
                const fheight = fcanvas.getHeight();
                if (top + MIN_HEIGHT > fheight) {
                  top = fheight - MIN_HEIGHT;
                }
              }

              const MIN_WIDTH = fontSize * 6;
              if (width < MIN_WIDTH) {
                width = MIN_WIDTH;
                const fwidth = fcanvas.getWidth();
                if (left + MIN_WIDTH > fwidth) {
                  left = fwidth - MIN_WIDTH;
                }
              }

              shape.set({
                left,
                top,
                width,
                height,
              });
            } else {
              shape.set({
                width,
                height,
              });
            }
          } else if (shape instanceof Ellipse) {
            const rx = width / 2;
            const ry = height / 2;

            shape.set({
              left,
              top,
              rx,
              ry,
            });
          }

          fcanvas.requestRenderAll();

          const { controller } = annotationStore;
          if (controller instanceof PersonAnnotationController) {
            const xy = {
              x: (shape.flipX ? shape.left - shape.width : shape.left) / scale,
              y: (shape.flipY ? shape.top - shape.height : shape.top) / scale,
            };

            const hw = {
              width: shape.width / scale,
              height: shape.height / scale,
            };

            if (shapeState.creating.value === 'text') {
              const textItem: TextAnnotation = {
                type: IDEAAnnotateType.IDEAAnnotateTypeTextBox,
                pageNumber: point.cur + 1,
                textBox: {
                  id: '',
                  ...xy,
                  ...hw,
                  fonts: [
                    {
                      fontSize:
                        localStorageColor.value.shapeFontSize ||
                        DEFAULT_FONT_SIZE,
                      color:
                        colorMap[
                          localStorageColor.value.shapeStyleId ||
                            AnnotationColor.blue
                        ],
                      content: '',
                    },
                  ],
                  borderColor: '',
                },
              };

              const textAnnotation = getNewAnnotation(point.cur + 1);

              console.log({
                ...textAnnotation,
                textItem,
              });

              const id = await controller.onlineSaveAnnotation(
                textAnnotation.documentId,
                textAnnotation.pageNumber,
                {
                  ...textAnnotation,
                  ...textItem,
                } as any
              );

              textItem.textBox!.id = id as string;

              setTextRectByAnnotation(textItem, shape as Rect);
              const { container } = documentViewer;
              const textDiv = createTextDiv(textItem, container);
              getShapeScale(point.cur + 1, container).appendChild(textDiv);
              const startEditText = useStartEditText(
                scale,
                textItem,
                textDiv,
                textCallback
              );
              onMouseDownStartEditText(textDiv, startEditText);
              initTextRect(scale, textItem, textDiv, textCallback);
              setEditingAndEditTextarea(id as string, textDiv);
              shape.set({
                stroke: 'transparent',
                strokeWidth: 0,
              });
              fcanvas.setActiveObject(shape);
              fcanvas.requestRenderAll();
            } else {
              const shapeAnnotation: ShapeAnnotation = {
                type: shapeState.creating.value as ShapeType,
                ...xy,
                width: 0,
                height: 0,
                uuid: '',
                pageNumber: point.cur + 1,
                strokeColor:
                  localStorageColor.value.shapeStyleId || AnnotationColor.blue,
                radiusX: 0,
                radiusY: 0,
                endX: 0,
                endY: 0,
              };

              if (shapeState.creating.value === ShapeType.rectangle) {
                shapeAnnotation.width = hw.width;
                shapeAnnotation.height = hw.height;
              } else if (shape instanceof Ellipse) {
                shapeAnnotation.radiusX = shape.rx / scale;
                shapeAnnotation.radiusY = shape.ry / scale;
                shapeAnnotation.x += shapeAnnotation.radiusX;
                shapeAnnotation.y += shapeAnnotation.radiusY;
              } else if (shape instanceof Line) {
                shapeAnnotation.x = shape.x1 / scale;
                shapeAnnotation.y = shape.y1 / scale;
                shapeAnnotation.endX = shape.x2 / scale;
                shapeAnnotation.endY = shape.y2 / scale;

                const length = Math.sqrt(
                  (shapeAnnotation.endX - shapeAnnotation.x) ** 2 +
                    (shapeAnnotation.endY - shapeAnnotation.y) ** 2
                );

                if (length < 10) {
                  cancelShape();
                  return;
                }
              }

              (async () => {
                const shapeId = await controller.createShape(shapeAnnotation);
                console.log({ shapeId, shapeAnnotation });
                if (!shapeId) {
                  fcanvas.remove(shape);
                  return;
                }

                shapeAnnotation.shapeId = shapeId;

                if (shape instanceof Line) {
                  const arrowhead = getArrowHead(NEW_ARROW_ID);
                  arrowhead?.setAttribute(
                    `data-${ARROWHEAD_CLASSNAME}`,
                    shapeId
                  );
                }

                console.warn({ shape, fcanvas });

                initShape(
                  shapeAnnotation,
                  shape,
                  documentViewer
                    .getPdfViewer()
                    .getPageView(shapeAnnotation.pageNumber),
                  useShapeCallback()
                );

                fcanvas.setActiveObject(shape);
                fcanvas.requestRenderAll();
              })();
            }
          }

          shapeState.creating.value = ShapeType.UNRECOGNIZED;
          shapeState.startPoint = null;
          shapeState.shape = null;
          return;
        }

        if (!(await latestMouseDown)) {
          documentViewer.onMouseUp(point);
        } else {
          const model = await getHandwriteModel(
            point as PointOnPdf,
            pdfAnnotaterRef.value
          );
          const handMap = await noteBuffer.handwriteBuffer;

          if (!model || !handMap) {
            return;
          }

          Object.values(handMap).some((list) => {
            return list.some((item) => {
              return item.webDrawItems.some((webDrawItem) => {
                if (webDrawItem.id !== model.id) {
                  return false;
                }

                const frame = renderOldFrame(
                  item.pageNumber,
                  getScale(documentViewer, item.pageNumber - 1),
                  model,
                  () => {
                    Modal.confirm({
                      title: shapeState.deleteTitle,
                      onOk() {
                        removeSingleHandwrite(
                          pdfAnnotaterRef.value!,
                          item.pageNumber,
                          model.id
                        );
                        frame.remove();
                      },
                      okButtonProps: {
                        danger: true,
                      },
                      cancelButtonProps: { type: 'primary' },
                      okText: shapeState.deleteOk,
                    });
                  }
                );
                return true;
              });
            });
          });

          clearMouseEvent();
        }
      },

      onDoubleClick: async (event) => {
        const point = documentViewer.getPoint(event);
        const ctrlPressing = detectCtrlPressing(event);

        const todocheckTextShape = async (p: PointOnPdf | null) => !p;

        if (ctrlPressing || !(await todocheckTextShape(point))) {
          documentViewer.onDoubleClick(event);
        } else {
          // TODO EDIT TEXT SHAPE
        }
      },
    },
    {}
  );
};

export const useTextCallback = (
  annotationStore: ReturnType<typeof useAnnotationStore>
): TextCallback => {
  const deleted: PersonAnnotationController['deleteTextBox'] = (...args) => {
    return (
      annotationStore.controller as PersonAnnotationController
    ).deleteTextBox(...args);
  };

  const edited: PersonAnnotationController['editTextBox'] = (...args) => {
    return (
      annotationStore.controller as PersonAnnotationController
    ).editTextBox(...args);
  };

  return {
    edited,
    changeColor(color) {
      const localStorageColor = useCommentGlobalState();
      localStorageColor.value.shapeStyleId = color;
    },
    deleted,
    confirmDelete: (textItem: TextAnnotation) =>
      new Promise((resolve) => {
        Modal.confirm({
          title: shapeState.deleteTitle,
          onOk() {
            deleted(textItem.textBox?.id ?? '');
            resolve(true);
          },
          onCancel() {
            resolve(false);
          },
          okButtonProps: {
            danger: true,
          },
          cancelButtonProps: { type: 'primary' },
          okText: shapeState.deleteOk,
        });
      }),
    createToolbarApp(toolbarProps, toolbarDiv) {
      const scApp = createApp(ShapeColor, toolbarProps as never);
      scApp.mount(toolbarDiv);
    },
  };
};

const clearAllTextRect = async (
  documentViewer: DocumentViewer,
  textCallback: TextCallback
) => {
  const [, textMap] = (await noteBuffer.annotationBuffer) || [];
  if (!textMap) {
    return;
  }

  // console.log('clear start', textMap);

  Object.entries(textMap).forEach(([pageNumber, list]) => {
    const fcanvas = PageFcanvasMap.get(Number(pageNumber)) as Canvas;
    list.forEach((item) => {
      const textRect = getTextRectByAnnotation(item);

      if (!textRect) {
        return;
      }

      fcanvas.remove(textRect);
      deleteTextRectByAnnotation(item);
      cancelEditText(
        item,
        textCallback,
        getTextDivByAnnotation(documentViewer.container, item)
      );
    });
  });
};

export const useShapeCallback = (
  annotationStore = useAnnotationStore(),
  localStorageColor = useCommentGlobalState()
): ShapeCallback => ({
  onUpdate(shapeAnnotation) {
    const controller = annotationStore.controller as PersonAnnotationController;
    controller.updateShape(shapeAnnotation);
  },
  onDelete(shapeId: string) {
    const controller = annotationStore.controller as PersonAnnotationController;
    controller.deleteShape(shapeId);
  },
  createToolbar(toolbarDiv, shapeAnnotation, onChangeColor) {
    const scApp = createApp(ShapeColor, {
      originColor: shapeAnnotation.strokeColor,
      onColorEnter(color: AnnotationColor) {
        onChangeColor(color);
      },
      onColorLeave() {
        onChangeColor(shapeAnnotation.strokeColor);
      },
      selectColor(color: AnnotationColor) {
        shapeAnnotation.strokeColor = color;
        onChangeColor(color);
        localStorageColor.value.shapeStyleId = color;
        const controller =
          annotationStore.controller as PersonAnnotationController;
        controller.updateShape(shapeAnnotation);
        const fcanvas = PageFcanvasMap.get(shapeAnnotation.pageNumber);
        fcanvas?.discardActiveObject();
        fcanvas?.requestRenderAll();
      },
    });
    scApp.mount(toolbarDiv);
  },
});

const getScale = (
  documentViewer: DocumentViewer,
  pageNumberFrom0 = 0
): number => {
  return documentViewer.getPdfViewer().getPageView(pageNumberFrom0).viewport
    .scale;
};

const removeSingleHandwrite = async (
  pdfAnnotater: PDFJSAnnotate,
  pageNumber: number,
  id: string
) => {
  deleteHandwrite({ id: [id] });

  const pageMap = await noteBuffer.handwriteBuffer;
  const pageList = pageMap?.[pageNumber] ?? [];

  const findAndDeleted = pageList.some((item) => {
    const index = item.webDrawItems.findIndex(
      (webDrawItem) => webDrawItem.id === id
    );
    if (index !== -1) {
      item.webDrawItems.splice(index, 1);
      return true;
    }
  });

  if (!findAndDeleted) {
    return;
  }

  renderSinglePageHandwrite(pageNumber, pageList, pdfAnnotater, true);
};
