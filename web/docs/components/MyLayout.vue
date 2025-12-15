<script setup lang="ts">
// =================================================================
// 区域 1: IMPORTS (导入)
// =================================================================
import { computed, onMounted, watch, onUnmounted } from 'vue'
import DefaultTheme from 'vitepress/theme'
import { useData } from 'vitepress'
import { UserOutlined, SettingOutlined, LogoutOutlined } from '@ant-design/icons-vue'

// 导入您的项目公共模块和 Store
import { useUserStore } from '@common/stores/user'
import { useVitePressI18n } from '../i18n'
import CustomFooter from './CustomFooter.vue'
import { getCurrentTheme, setTheme, onThemeChange, isDarkTheme } from '@common/theme'

// ⭐️ 新增：导入我们将要渲染的 Hero 组件
import Hero from '../components/index/Hero.vue'

// =================================================================
// 区域 2: STATE & LOGIC (状态与逻辑)
// =================================================================

// --- 初始化 VitePress 和 Pinia ---
const { Layout } = DefaultTheme
// 从 useData() 中获取渲染所需的所有响应式数据
const { frontmatter, isDark, lang } = useData() 
const userStore = useUserStore()
const { t } = useVitePressI18n()

// ⭐️ 新增：Hero 组件的翻译内容 (使用标准语言格式)
const heroTranslations = {
  'en-US': {
    text: "A new experience of doc reading in AI age",
    tagline: "From reading to conversation.",
    workbenchText: "My Library",
    guideText: "User Guide"
  },
  'zh-CN': {
    text: "AI时代文献阅读新范式",
    tagline: "不止阅读，更是思想碰撞",
    workbenchText: "我的文库",
    guideText: "用户手册"
  }
}

// ⭐️ 新增：最终的、健壮的 heroData 计算属性
const heroData = computed(() => {
  // 1. 获取当前语言，并设置默认回退到英文
  const currentLang = lang.value;
  console.log('[heroData] 当前语言:', currentLang); // 调试日志
  
  // 2. 标准化语言标识符
  let standardLang = currentLang;
  if (currentLang === 'zh') {
    standardLang = 'zh-CN';
  } else if (currentLang === 'en') {
    standardLang = 'en-US';
  }
  
  console.log('[heroData] 标准语言:', standardLang); // 调试日志
  
  // 3. 获取对应的翻译内容
  const translations = heroTranslations[standardLang] || heroTranslations['en-US'];

  // 2. 获取用户登录状态
  const isLoggedIn = userStore.isLogin();

  // 3. 动态构建能够正确跳出 /docs/ 的完整 URL
  let workbenchLink = '#'; // 提供一个在服务器端渲染时的备用值
  
  // 这段逻辑只在浏览器客户端执行，以安全地访问 window 对象
  if (typeof window !== 'undefined') {
    const origin = window.location.origin; // e.g., "http://localhost:3000" or "https://your-domain.com"
    
    // 根据登录状态，构建指向主应用的完整链接
    workbenchLink = isLoggedIn 
        ? `${origin}/workbench/recent` 
        : `${origin}/account/login?redirect_url=${encodeURIComponent(window.location.href)}`;
  }
  
  // 4. 返回最终的数据对象，供 Hero 组件渲染
  return {
    name: "ODOC.AI", // 品牌名保持不变
    text: translations?.text || "A new experience of doc reading in AI age",
    tagline: translations?.tagline || "From reading to conversation.",
    actions: [
      // 工作台链接，指向主应用
      { theme: 'brand', text: translations?.workbenchText || "Workbench", link: workbenchLink },
      // 指南链接，仍然是 VitePress 内部链接
      { theme: 'alt', text: translations?.guideText || "User Guide", link: `/docs/${currentLang === 'ZH_CN' ? 'zh/' : ''}guide` }
    ]
  }
})


// --- 以下是您原文件中所有的既有逻辑，保持完整无缺 ---

const isLoggedIn = computed(() => userStore.isLogin())
const userName = computed(() => userStore.userInfo?.nickname || userStore.userInfo?.username || t('layout.userProfile'))
const userAvatar = computed(() => userStore.userInfo?.avatarUrl || '')

let themeCleanup: (() => void) | null = null

onMounted(() => {
  const updateWorkbenchLink = () => {
    if (typeof window === 'undefined') return;

    // 查找包含占位符的导航链接
    const navLink = document.querySelector<HTMLAnchorElement>('a.VPNavBarMenuLink[href*="/workbench_placeholder"]');
    
    if (navLink) {
      const origin = window.location.origin;
      const correctUrl = `${origin}/workbench`;
      
      // 仅在链接不正确时才更新，避免不必要的 DOM 操作
      if (navLink.href !== correctUrl) {
        navLink.href = correctUrl;
      }
    }
  };

  // 初始加载时执行一次
  updateWorkbenchLink();

  // 使用 MutationObserver 监听 DOM 变化，以应对 VitePress 的客户端路由
  // 这确保了即使用户在文档页面之间切换，链接也能被及时修正
  const observer = new MutationObserver((mutations) => {
    // 简单地在任何变化后重新运行更新逻辑
    updateWorkbenchLink();
  });

  // 观察整个文档的变化，因为我们不确定 VitePress 具体会更新哪个部分
  observer.observe(document.body, { childList: true, subtree: true });

  // 在组件卸载时停止观察，防止内存泄漏
  onUnmounted(() => {
    observer.disconnect();
  });

  userStore.getUserInfo()
  const currentTheme = getCurrentTheme()
  isDark.value = isDarkTheme(currentTheme)
  themeCleanup = onThemeChange((newTheme) => {
    const shouldBeDark = isDarkTheme(newTheme)
    if (isDark.value !== shouldBeDark) {
      isDark.value = shouldBeDark
    }
  })
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('data-theme', isDark.value ? 'dark' : 'light')
  }
})

onUnmounted(() => {
  if (themeCleanup) {
    themeCleanup()
  }
})

watch(isDark, (newValue) => {
  const themeValue = newValue ? 'dark' : 'default'
  setTheme(themeValue)
})

function handleLogin() {
  userStore.openLogin()
  if (typeof window !== 'undefined') {
    const baseUrl = window.location.origin
    const loginUrl = `${baseUrl}/account/login?redirect_url=${encodeURIComponent(window.location.href)}`
    setTimeout(() => {
      window.location.replace(loginUrl)
    }, 100)
  }
}

async function handleLogout() {
  try {
    await userStore.logout()
    if (typeof window !== 'undefined') {
      window.location.reload()
    }
  } catch (error) {
    console.error('Logout failed:', error)
  }
}

function handleProfile() {
  if (!isLoggedIn.value) {
    handleLogin()
    return
  }
  const userId = localStorage.getItem('userId')
  if (userId && typeof window !== 'undefined') {
    window.open(`/user/${userId}`, '_blank')
  }
}

function handleSettings() {
  if (typeof window !== 'undefined') {
    window.open('/settings', '_blank')
  }
}
</script>

<template>
  <Layout>
    <!-- 
      ⭐️ 变更：使用 #home-hero-before 插槽。
      它由 VitePress 默认主题提供，专门用于在首页渲染自定义的 Hero 区域。
      这个插槽会在默认的 features 内容之前显示。
    -->
    <template #home-hero-before>
      <!-- 我们只在 VitePress 识别为首页时才渲染这个 Hero 组件 -->
      <div v-if="frontmatter.layout === 'home'">
        <Hero v-bind="heroData" />
      </div>
    </template>

    <!-- 在站点标题右侧添加 BEAT 微标 -->
    <template #nav-bar-title-after>
      <span class="vp-site-badge">BETA</span>
    </template>
    
    <!-- 您现有的导航栏用户菜单功能，保持不变 -->
    <template #nav-bar-content-after>
      <div class="user-avatar-container">
        <template v-if="isLoggedIn">
          <a-dropdown placement="bottomRight">
            <div class="avatar-wrapper">
              <a-avatar :src="userAvatar" class="user-avatar">
                <template #icon v-if="!userAvatar">
                  {{ userName.charAt(0).toUpperCase() }}
                </template>
              </a-avatar>
            </div>
            <template #overlay>
              <a-menu class="user-dropdown-menu">
                <!-- <a-menu-item key="profile" @click="handleProfile">
                  <template #icon><user-outlined /></template>
                  <span>{{ t('layout.userProfile') }}</span>
                </a-menu-item>
                <a-menu-item key="settings" @click="handleSettings">
                  <template #icon><setting-outlined /></template>
                  <span>{{ t('layout.settings') }}</span>
                </a-menu-item>
                <a-menu-divider /> -->
                <a-menu-item key="logout" @click="handleLogout">
                  <template #icon><logout-outlined /></template>
                  <span>{{ t('layout.logout') }}</span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </template>
        <template v-else>
          <a class="nav-link" @click="handleLogin">{{ t('layout.login') }}</a>
        </template>
      </div>
    </template>
    
    <!-- 您现有的自定义页脚，保持不变 -->
    <template #layout-bottom>
      <CustomFooter />
    </template>
  </Layout>
</template>

<!-- 您的所有现有样式，保持不变 -->
<style scoped>
.user-avatar-container {
  display: flex;
  align-items: center;
  margin-left: 16px;
}

.avatar-wrapper {
  cursor: pointer;
}

:deep(.ant-avatar) {
  background-color: #1e88e5;
  color: white;
}

:deep(.ant-dropdown-menu) {
  background-color: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
}

:deep(.ant-dropdown-menu-item) {
  color: var(--vp-c-text-1);
}

:deep(.ant-dropdown-menu-item:hover) {
  background-color: var(--vp-c-bg-mute);
}

:deep(.ant-menu-item .anticon) {
  margin-right: 8px;
}

.login-button {
  background: none;
  border: none;
  padding: 0;
  font: inherit;
  cursor: pointer;
  outline: inherit;
}

.nav-link {
  display: block;
  margin-left: 16px;
  font-size: 14px;
  font-weight: 500;
  color: var(--vp-c-text-1);
  transition: color 0.25s;
  cursor: pointer;
}

.nav-link:hover {
  color: var(--vp-c-brand);
}
</style>