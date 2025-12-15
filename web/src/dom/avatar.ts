/*
 * Created Date: May 27th 2022, 3:54:40 pm
 * Author: zhoupengcheng
 * -----
 * Last Modified: June 9th 2022, 10:59:08 am
 */
import Avatar from '@/components/Right/TabPanel/Note/group/Avatar.vue';
import { createApp } from 'vue';
import Antd from 'ant-design-vue';
import {
  PDF_ANNOTATIONLAYER,
  PDF_ANNOTATIONLAYER_GROUP_NOTES,
} from '@idea/pdf-annotate-core';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import { useStore } from '@/store';
import { useAnnotationStore } from '../stores/annotationStore';

export const createAvatar = (
  item: {
    top: number;
    left: number;

    pdfAnnotateId: string;
    pageNumber: number;
  },
  pdfViewer: ViewerController
) => {
  const first = useAnnotationStore().controller.findAnnotation(
    item.pdfAnnotateId
  );
  if (first && first.pageNumber !== item.pageNumber) {
    return;
  }

  const app = createApp(Avatar, { item });

  app.use(Antd);
  app.use(useStore());

  const instance = app.mount(document.createElement('div'));

  const container = pdfViewer?.getDocumentViewer().container;
  const page = container?.querySelector(
    `.page[data-page-number="${item.pageNumber}"]`
  );

  page?.appendChild(instance.$el);

  const svg = page?.querySelector<SVGGElement>(
    `svg.${PDF_ANNOTATIONLAYER} .${PDF_ANNOTATIONLAYER_GROUP_NOTES}`
  );

  if (svg?.style.display === 'none') {
    instance.$el.style.display = 'none';
  }
};

export const clearAvatar = () => {
  const avatars = document.querySelectorAll('.group-avatar');

  avatars.forEach((ava) => {
    ava.parentNode?.removeChild(ava);
  });
};

export const showAvatar = () => {
  const avatars = document.querySelectorAll('.group-avatar');

  avatars.forEach((ava) => {
    (ava as any).style.display = 'block';
  });
};

export const hideAvatar = () => {
  const avatars = document.querySelectorAll('.group-avatar');

  avatars.forEach((ava) => {
    (ava as any).style.display = 'none';
  });
};

export const toggleAvatar = (
  curUserId: string,
  isShowGroupOther: boolean,
  isShowGroupMine: boolean,
  isShowGroupImage: boolean
) => {
  setTimeout(() => {
    const avatars = document.querySelectorAll('.group-avatar');

    avatars.forEach((ava: any) => {
      const userId = ava.getAttribute('data-userid');

      if (userId === curUserId) {
        ava.style.display =
          isShowGroupMine && isShowGroupImage ? 'block' : 'none';
      } else {
        ava.style.display =
          isShowGroupOther && isShowGroupImage ? 'block' : 'none';
      }
    });
  }, 0);
};
