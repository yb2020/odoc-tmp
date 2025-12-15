---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

# hero:
#   name: "odoc.ai"
#   text: "An new experience of reading doc for AI age"
#   tagline: More than reading, fosters conversation.
#   actions:
#     - theme: brand
#       text: Workbench
#       link: ../workbench/recent
#     - theme: alt
#       text: User Guide
#       link: /guide
   
features:
  - title: Vibe Reading
    details: A fresh reading experience with AI 
  - title: Conversation with Doc
    details: Reading with ai, Conversation with Doc
  - title: Translate anytime
    details: Translate anytime you need help
---

<script setup>
import ProductionFeatureAccordion from './components/index/ProductionFeatureAccordion.vue'
import AiNativeReaderFeatureAccordion from './components/index/AiNativeReaderFeatureAccordion.vue'
import AIAssistedAnalysisFeatureAccordion from './components/index/AIAssistedAnalysisFeatureAccordion.vue'
import ProfessionalTranslationAccordion from './components/index/ProfessionalTranslationAccordion.vue'
import TargetUserTabs from './components/index/TargetUserTabs.vue'
import DifferWithChatGPT from './components/index/DifferWithChatGPT.vue'
import UserTestimonials from './components/index/UserTestimonials.vue'
import OurPartners from './components/index/OurPartners.vue'
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
