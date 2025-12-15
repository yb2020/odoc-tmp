<template>
  <footer class="custom-footer">
    <div class="footer-container">
      <div class="footer-logo">
        <!-- <img src="/logo.png" alt="odoc.ai" class="logo-image" /> -->
        <!-- <span class="logo-text">odoc.ai</span> -->
        <img src="/logo_text.svg" alt="odoc.ai" class="logo-image" />
      </div>
      
      <div class="footer-sections">
        <div class="footer-section">
          <h3>{{ t('footer.product.title') }}</h3>
          <ul>
            <li><a :href="productFeaturesLink">{{ t('footer.product.features') }}</a></li>
            <li><a :href="productUsageLink">{{ t('footer.product.usage') }}</a></li>
            <li><a :href="productUserAgreementLink">{{ t('footer.product.userAgreement') }}</a></li>
            <li><a :href="productPrivacyLink">{{ t('footer.product.privacy') }}</a></li>
          </ul>
        </div>
        
        <div class="footer-section">
          <h3>{{ t('footer.team.title') }}</h3>
          <ul>
            <li><a :href="aboutUsLink">{{ t('footer.team.aboutUs') }}</a></li>
            <li><a :href="updatesLink">{{ t('footer.team.updates') }}</a></li>
            <li><a href="https://github.com/yb2020/odoc.ai.feedback/issues" target="_blank">{{ t('footer.team.contactUs') }}</a></li>
            <li><a :href="joinUsLink">{{ t('footer.team.joinUs') }}</a></li>
          </ul>
        </div>
        
        <div class="footer-section">
          <h3>{{ t('footer.social.title') }}</h3>
          <ul class="social-links">
            <li><a :href="oneLink" target="_blank">{{ t('footer.social.one') }}</a></li>
            <li><a :href="twoLink" target="_blank">{{ t('footer.social.two') }}</a></li>
            <li><a :href="threeLink" target="_blank">{{ t('footer.social.three') }}</a></li>
          </ul>
        </div>
      </div>
    </div>
    
    <div class="footer-bottom">
      <div class="copyright">
        &copy; {{ new Date().getFullYear() }} odoc.ai. {{ t('footer.copyright') }}
      </div>
    </div>
  </footer>
</template>

<script setup>
// ⭐️ 1. Import `computed` from Vue and `useData` from VitePress
import { computed } from 'vue';
import { useData } from 'vitepress';
import { useVitePressI18n } from '../i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '../../src/hooks/useLanguage';

// Use VitePress's i18n system for text
const { t } = useVitePressI18n();

// ⭐️ 2. Get reactive language data from VitePress
const { lang } = useData();

// 语言管理
const { isCurrentLanguage } = useLanguage();

// Create a computed property for the "About Us" link
const aboutUsLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide/team/about'; 
  }
  return '/docs/guide/team/about';
});

const updatesLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide/team/update'; 
  }
  return '/docs/guide/team/update';
});

const joinUsLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide/team/join'; 
  }
  return '/docs/guide/team/join';
});

const productFeaturesLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide/product/features'; 
  }
  return '/docs/guide/product/features';
});

const productUsageLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide/product/usage'; 
  }
  return '/docs/guide/product/usage';
});

const productUserAgreementLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide/product/userAgreement'; 
  }
  return '/docs/guide/product/userAgreement';
});

const productPrivacyLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide/product/privacy'; 
  }
  return '/docs/guide/product/privacy';
});

// 社交媒体
const oneLink = computed(() => {
  return 'https://www.facebook.com';
});

const twoLink = computed(() => {
  return 'https://www.reddit.com';
});

const threeLink = computed(() => {
  return 'https://www.whatsapp.com';
});
</script>

<style scoped>
.custom-footer {
  background-color: #1A1A1A;
  color: #fff;
  padding: 40px 0 20px;
  font-family: var(--vp-font-family-base);
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  flex-wrap: wrap;
}

.footer-logo {
  display: flex;
  align-items: center;
  margin-right: 60px;
  margin-bottom: 30px;
}

.logo-image {
  margin-right: 10px;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
}

.footer-sections {
  display: flex;
  flex: 1;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 50px;
}

.footer-section {
  margin-right: 0;
  margin-bottom: 30px;
  min-width: 120px;
}

.footer-section h3 {
  font-size: 16px;
  margin-bottom: 15px;
  color: #fff;
  font-weight: 500;
}

.footer-section ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-section li {
  margin-bottom: 10px;
}

.footer-section a {
  color: #aaa;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}

.footer-section a:hover {
  color: #fff;
}

.social-links {
  display: flex;
  flex-direction: column;
}

.footer-bottom {
  margin: 0 auto;
  padding: 20px 20px 0;
  border-top: 1px solid #262a3c;
  margin-top: 10px;
}

.copyright {
  font-size: 12px;
  color: #888;
  text-align: center;
}

@media (max-width: 768px) {
  .footer-container {
    flex-direction: column;
  }
  
  .footer-sections {
    flex-direction: column;
  }
  
  .footer-section {
    margin-right: 0;
    margin-bottom: 20px;
    width: 100%;
  }
}
</style>
