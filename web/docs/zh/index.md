---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

# hero:
#   name: "odoc.ai"
#   text: "AI时代文献阅读新范式"
#   tagline: 不止阅读，更是思想碰撞
#   actions:
#     - theme: brand
#       text: 工作台
#       link: ../workbench/recent
#     - theme: alt
#       text: 用户手册
#       link: /guide

features:
  - title: 氛围阅读
    details: AI时代的阅读新体验
  - title: 沟通成文
    details: 阅读与AI的对话，生成文档
  - title: 翻译无处不在
    details: 随时随地的翻译
---

<script setup>
import ProductionFeatureAccordion from '../components/index/ProductionFeatureAccordion.vue'
import AiNativeReaderFeatureAccordion from '../components/index/AiNativeReaderFeatureAccordion.vue'
import AIAssistedAnalysisFeatureAccordion from '../components/index/AIAssistedAnalysisFeatureAccordion.vue'
import ProfessionalTranslationAccordion from '../components/index/ProfessionalTranslationAccordion.vue'
import TargetUserTabs from '../components/index/TargetUserTabs.vue'
import DifferWithChatGPT from '../components/index/DifferWithChatGPT.vue'
import UserTestimonials from '../components/index/UserTestimonials.vue'
import OurPartners from '../components/index/OurPartners.vue'
</script>
<ClientOnly>
  <AiNativeReaderFeatureAccordion />

  <AIAssistedAnalysisFeatureAccordion reverse />

  <ProductionFeatureAccordion />

  <ProfessionalTranslationAccordion reverse />

  <TargetUserTabs />

  <DifferWithChatGPT />

  <UserTestimonials />

  <OurPartners />
</ClientOnly>
