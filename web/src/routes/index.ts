import type { RouteRecordRaw } from 'vue-router';
import { createWebHistory, createRouter } from 'vue-router';

// 使用懒加载导入所有页面组件
const DefaultLayout = () => import('@/layouts/default.vue');
const WorkBenchPage = () => import('@/pages/workBench.vue');
const NotePage = () => import('@/pages/note.vue');
const ErrorPage = () => import('@/pages/exception/index.vue');
const ForbiddenPage = () => import('@/pages/exception/forbidden.vue');
const LibraryPage = () => import('@/pages/library/index.vue');
import { selfNoteInfo, store } from '../store';
import { BaseActionTypes } from '../store/base';
import { ResponseError } from '../api/type';
import {
  ERROR_CODE_UNLOGIN,
  UNKNOWN_ERROR_CODE,
  UNKNOWN_ERROR_MESSAGE,
} from '../api/const';
import { bridgeAdaptor } from '../adaptor/bridge';
import { setSeo } from '../util/seo';
import { getDomainOrigin, isInOverseaseElectron } from '../util/env';
import { DocumentsActionTypes } from '../store/documents';
import { PAGE_ROUTE_NAME } from './type';
import i18n, { LanguageCodeMap } from '../locals/i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import langs from '../locals/lang.json';
import { useEnvStore } from '../stores/envStore';
import { setLanguageCookie, standardToLanguageEnum } from '../shared/language/service';
import { useBaseStore } from '../stores/baseStore';
import { useUserStore } from '../common/src/stores/user';

// 懒加载 EmbedLayout
const EmbedLayout = () => import('@/layouts/embed.vue');

export const ROOT_PATH = '';

const createHistory = () => {
  return createWebHistory(ROOT_PATH);
};

const createI18nRules = () => {
  const meta = {
    auth: false,
  };
  const getAlias = (path: string) => {
    const alias = [`${path}/`];
    langs.lang.map((lang) => {
      alias.push(`${lang}/${path}`);
      alias.push(`${lang}/${path}/`);
    });
    return alias;
  };
  const rules = [
    {
      path: 'workbench',
      component: WorkBenchPage,
      name: PAGE_ROUTE_NAME.WORKBENCH,
      alias: getAlias('workbench'),
      meta: { ...meta, auth: true }, // 需要认证
      children: [
        {
          path: 'recent',
          component: () => import('@/pages/RecentReading.vue'),
          name: PAGE_ROUTE_NAME.RECENT_READING_IN_WORKBENCH,
          meta: { ...meta, auth: true, title: i18n.global.t('recent') },
        },
        {
          path: 'library',
          component: () => import('@/pages/library/index.vue'),
          name: PAGE_ROUTE_NAME.LIBRARY_IN_WORKBENCH,
          meta: { ...meta, auth: true, title: i18n.global.t('library') },
        },
        {
          path: 'notes',
          component: () => import('@/pages/notes/index.vue'),
          name: PAGE_ROUTE_NAME.NOTES_IN_WORKBENCH,
          meta: { ...meta, auth: true, title: i18n.global.t('notes') },
        }
      ],
      redirect: '/workbench/recent'
    },
    {
      path: 'library',
      component: () => import('@/pages/library/index.vue'),
      name: PAGE_ROUTE_NAME.LIBRARY,
      alias: getAlias('library'),
      meta: { ...meta, auth: true }, // 需要认证
    },
    {
      path: 'note',
      component: NotePage,
      name: PAGE_ROUTE_NAME.NOTE,
      alias: getAlias('note'),
      meta: { ...meta, auth: true }, // 需要认证
    },
    {
      path: 'error',
      component: ErrorPage,
      name: PAGE_ROUTE_NAME.EXCEPTION,
      alias: getAlias('error'),
      meta,
    },
    {
      path: '403',
      component: ForbiddenPage,
      name: PAGE_ROUTE_NAME.FORBIDDEN,
      alias: getAlias('403'),
      meta: meta,
    },
    {
      path: '/:path(.*)*',
      name: PAGE_ROUTE_NAME.EXCEPTION_404,
      component: ErrorPage,
      meta,
    },
  ];

  return {
    path: '/',
    component: DefaultLayout,
    redirect: '/workbench',
    children: rules,
  };
};

const createAuthRoutes = () => {
  return {
    path: '/account',
    component: EmbedLayout,
    children: [
      {
        path: 'login',
        component: () => import('@/pages/account/login.vue'),
        name: PAGE_ROUTE_NAME.LOGIN,
        meta: { auth: false },
      },
    ],
  };
};

function createVueRouter() {
  const routes: RouteRecordRaw[] = [createI18nRules(), createAuthRoutes()];

  // 添加 Notes 路由
  routes.push(
    {
      path: '/notes',
      component: () => import(/* webpackChunkName: "notes" */ '@/pages/notes/index.vue'),
      name: PAGE_ROUTE_NAME.NOTES,
      meta: { auth: true, title: i18n.global.t('notes') },
    }
  );

  const router = createRouter({
    history: createHistory(),
    routes: routes,
    strict: true,
    scrollBehavior: () => ({ left: 0, top: 0 }),
  });

  router.beforeEach(async (to) => {
    const matches = to.fullPath.match(/\/(.+)\/.+/);

    let lang = matches?.[1];

    if (lang && !langs.lang.includes(lang)) {
      lang = undefined;
    } else if (lang) {
      i18n.global.locale.value = lang;
    }

    if (isInOverseaseElectron()) {
      // 使用统一的语言管理服务设置 Cookie
      const currentLocale = i18n.global.locale.value;
      // 将标准格式转换为 proto 枚举，然后通过统一服务设置 Cookie
      try {
        const languageEnum = standardToLanguageEnum(currentLocale);
        setLanguageCookie(currentLocale as 'en-US' | 'zh-CN');
      } catch (error) {
        // 如果转换失败，使用默认语言
        console.warn('[routes] Invalid language format:', currentLocale);
      }
    }

    const messages = i18n.global.messages.value;
    const locale = i18n.global.locale.value;
    const htmlMessages = messages[locale]?.html;
    const metaTitle = htmlMessages?.pageTitle;

    if (metaTitle && to.name) {
      setSeo(
        (metaTitle[to.name as keyof typeof metaTitle] as string) || 'vibe reading - odoc.ai'
      );
    }

    const envStore = useEnvStore();
    envStore.initViewerConfig();

    // 检查路由是否需要认证
    const requiresAuth = to.meta.auth !== false;
    
    // 如果路由需要认证，检查用户是否已登录
    if (requiresAuth) {
      try {
        const userStore = useUserStore();
        await userStore.getUserInfo();
        
        if (!userStore.isLogin()) {
          // 用户未登录，重定向到登录页面
          return {
            name: PAGE_ROUTE_NAME.LOGIN,
            query: { redirect: to.fullPath },
          };
        }
      } catch (error) {
        // 获取用户信息失败，重定向到登录页面
        return {
          name: PAGE_ROUTE_NAME.LOGIN,
          query: { redirect: to.fullPath },
        };
      }
    }

    // 继续原有的路由处理逻辑...
    if (to.name !== PAGE_ROUTE_NAME.NOTE) {
      return true;
    }

    if (!to.query.pdfId && !to.query.noteId) {
      function handleException(error: ResponseError, lang?: string) {
        const pathOrName = lang
          ? { path: `/${lang}/error` }
          : { name: PAGE_ROUTE_NAME.EXCEPTION };
        const { fullPath } = router.resolve({
          ...pathOrName,
          query: {
            status: (error as ResponseError).code || UNKNOWN_ERROR_CODE,
            message: encodeURIComponent(
              (error as ResponseError).message || UNKNOWN_ERROR_MESSAGE
            ),
            source: encodeURIComponent(window.location.href),
          },
        });
        router.push(fullPath);
      }

      handleException(
        new ResponseError({
          code: 400,
          message: 'Bad Request',
        }),
        lang
      );
      return false;
    }

    /**
     * https://www.tapd.cn/57741831/markdown_wikis/show/#1157741831001000305
     * PDF渲染和权限以链接参数pdfId为准
     * 笔记相关以链接参数noteId为准
     * 可能存在展示的PDF和note挂的PDF不是同一个的情况，也就是笔记错位，这个anson已知
     * noteId仅和paperId一一对应，也就是对于同一个paperId用户的noteId仅为1个，但是一个paperId下面用户有权限的pdfId可能多个
     *
     */
    if (!to.query.noteId) {
      try {
        const noteId = await store.dispatch(
          `base/${BaseActionTypes.GET_NOTEID}`,
          {
            pdfId: to.query.pdfId,
            groupId: to.query.groupId,
          }
        );
        const fullPath = resolveNote(
          {
            noteId,
            pdfId: to.query.pdfId as string,
          },
          lang
        );
        router.push(fullPath);
      } catch (error) {
        handleError(error as ResponseError);
      }

      return true;
    }

    try {
      const pdfId = await store.dispatch(
        `base/${BaseActionTypes.GET_SELF_NOTEINFO}`,
        {
          noteId: to.query.noteId,
          groupId: to.query.groupId,
        }
      );
      // getPdfStatusInfo需要参数pdfId和noteId
      await Promise.all([
        store.dispatch(`base/${BaseActionTypes.GET_STATUSINFO}`, {
          pdfId: to.query.pdfId || pdfId,
          noteId: to.query.noteId,
        }),
        // 使用 Pinia 用户 store 替代 Vuex
        useUserStore().getUserInfo(),
      ]);
      await store.dispatch(
        `documents/${DocumentsActionTypes.INIT_PDFVIEWER_WITH_SETTING}`,
        to.query.noteId
      );
    } catch (error) {
      handleError(error as ResponseError);
      return true;
    }

    // 这里GET_NOTEINFO下发的noteId可能和query上的不一样，做一次重定向
    if (
      selfNoteInfo.value?.noteId &&
      String(selfNoteInfo.value?.noteId) !== String(to.query.noteId)
    ) {
      const fullPath = resolveNote({
        noteId: String(selfNoteInfo.value?.noteId),
      });
      window.location.replace(location.origin + ROOT_PATH + fullPath);
      return false;
    }

    const title = selfNoteInfo.value?.docName || selfNoteInfo.value?.paperTitle;
    if (title) {
      setSeo(title);
    }
    return true;

    function handleError(error: ResponseError) {
      if (error.code === ERROR_CODE_UNLOGIN) {
        const redirectUrl = `${getDomainOrigin()}${ROOT_PATH}${to.fullPath}`;
        bridgeAdaptor.login(redirectUrl);
      } else {
        useBaseStore().pageError = error;
      }
    }

    function resolveNote(
      params: { noteId?: string; pdfId?: string },
      lang?: string
    ) {
      const pathOrName = lang
        ? { path: `/${lang}/note` }
        : { name: PAGE_ROUTE_NAME.NOTE };
      const { fullPath } = router.resolve({
        ...pathOrName,
        query: {
          ...to.query,
          ...params,
        },
      });
      return fullPath;
    }
  });

  return router;
}

const router = createVueRouter();

export default router;
