<!-- docs/.vitepress/theme/components/index/Hero.vue -->
<script setup>
import { computed } from 'vue'

// 定义组件期望接收的 props
const props = defineProps({
  name: { type: String },
  text: { type: String },
  tagline: { type: String },
  actions: { type: Array, default: () => [] }
})

// 格式化 actions 数据以供模板使用
const formattedActions = computed(() => {
  return props.actions.map(action => ({
    theme: action.theme || 'brand',
    text: action.text,
    link: action.link
  }))
})
</script>

<template>
  <!-- 使用 VitePress 默认主题的类名来确保样式一致 -->
  <div class="VPHero VPHomeHero">
    <div class="container">
      <div class="main">
        <h1 v-if="name" class="name">
          <span class="clip" v-html="name"></span>
        </h1>
        <p v-if="text" class="text" v-html="text"></p>
        <p v-if="tagline" class="tagline" v-html="tagline"></p>

        <div v-if="formattedActions.length" class="actions">
          <div v-for="action in formattedActions" :key="action.link" class="action">
            <a
              :class="['VPButton', action.theme]"
              :href="action.link"
              >{{ action.text }}</a
            >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<style scoped>
/* 基础布局样式，保持不变 */
.VPHomeHero {
  margin-top: calc((var(--vp-nav-height) + var(--vp-layout-top-height, 0px)) * -1);
  padding-top: calc(var(--vp-nav-height) + var(--vp-layout-top-height, 0px) + 48px);
  padding-bottom: 48px;
}

@media (min-width: 640px) {
  .VPHomeHero {
    padding-top: calc(var(--vp-nav-height) + var(--vp-layout-top-height, 0px) + 80px);
    padding-bottom: 80px;
  }
}

.container {
  margin: 0 auto;
  max-width: 1152px;
}

.main {
  text-align: left;
  padding: 0 24px;
}

/* 标题样式，保持不变 */
.name, .text {
  margin: 0;
  max-width: 768px; 
  font-size: 56px;
  font-weight: 700;
  line-height: 1.25;
  letter-spacing: -0.5px;
}

.clip {
  background: var(--vp-home-hero-name-background);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: var(--vp-home-hero-name-color);
}

/* 
  ⭐️ 变更 1: 修改 .tagline 样式 
*/
.tagline {
  margin: 0 0 0 0;
  max-width: 768px;
  /* 应用您从 F12 中找到的样式 */
  font-size: 24px;
  line-height: 36px;
  /* 添加加粗效果 */
  font-weight: 600; 
  color: var(--vp-c-text-2);
}

.actions {
  display: flex;
  justify-content: flex-start;
  margin-top: 24px;
  gap: 12px;
}

/* 
  ⭐️ 变更 2: 为按钮（.VPButton）添加明确的、完整的样式定义 
  这能确保我们的按钮和 VitePress 默认主题中的按钮外观完全一致
*/
.VPButton {
  display: inline-block;
  border: 1px solid transparent;
  text-align: center;
  font-weight: 600;
  white-space: nowrap;
  transition: color 0.25s, border-color 0.25s, background-color 0.25s;
  /* 应用您从 F12 中找到的 .medium 尺寸样式 */
  border-radius: 20px;
  padding: 0 20px;
  line-height: 38px;
  font-size: 14px;
}

/* 为 brand 主题按钮（Workbench）应用颜色 */
.VPButton.brand {
  border-color: var(--vp-button-brand-border);
  color: var(--vp-button-brand-text);
  background-color: var(--vp-button-brand-bg);
}

/* 为 alt 主题按钮（User Guide）应用颜色 */
.VPButton.alt {
  border-color: var(--vp-button-alt-border);
  color: var(--vp-button-alt-text);
  background-color: var(--vp-button-alt-bg);
}

/* 添加 hover 和 active 状态的反馈，使其更具交互性 */
.VPButton:hover {
  border-color: var(--vp-button-brand-hover-border);
  color: var(--vp-button-brand-hover-text);
  background-color: var(--vp-button-brand-hover-bg);
}

.VPButton.alt:hover {
  border-color: var(--vp-button-alt-hover-border);
  color: var(--vp-button-alt-hover-text);
  background-color: var(--vp-button-alt-hover-bg);
}
</style>