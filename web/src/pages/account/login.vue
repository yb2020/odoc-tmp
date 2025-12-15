<template>
  <div class="login-container">
      <div class="login-box">
        <div class="logo">
          <div class="logo-text">ODOC.AI</div>
        </div>
        <h1 class="title">{{ t('account.loginPage.title') }}</h1>
        
        <!-- Google 登录块 -->
        <div v-if="loginConfig?.isGoogleLoginEnabled" class="social-login">
          <button class="google-btn" @click="handleGoogleLogin">
            <img src="https://www.gstatic.com/firebasejs/ui/2.0.0/images/auth/google.svg" alt="Google" />
            {{ t('account.loginPage.continueWithGoogle') }}
          </button>
        </div>
        
        <!-- 用户名密码登录块 -->
        <div v-if="loginConfig?.isUsernamePasswordLoginEnabled" class="username-password-login">
          <div v-if="loginConfig?.isGoogleLoginEnabled" class="divider">
            <span>{{ t('account.loginPage.orContinueWithEmail') }}</span>
          </div>
          
          <a-form :model="formState" @finish="handleSubmit" class="login-form">
            <a-form-item
              name="email"
              :rules="[{ required: true, message: t('account.loginPage.emailRequired') }]"
            >
              <a-input 
                v-model:value="formState.email" 
                type="email" 
                :placeholder="t('account.loginPage.emailPlaceholder')"
                size="large"
                class="email-input"
                :style="{ lineHeight: '24px' }"
              >
                <template #prefix>
                  <mail-outlined class="form-icon" />
                </template>
              </a-input>
            </a-form-item>
            
            <a-form-item
              name="password"
              :rules="[{ required: true, message: t('account.loginPage.passwordRequired') }]"
            >
              <a-input-password 
                v-model:value="formState.password" 
                :placeholder="t('account.loginPage.passwordPlaceholder')"
                size="large"
                class="password-input"
                :style="{ lineHeight: '24px' }"
              >
                <template #prefix>
                  <lock-outlined class="form-icon" />
                </template>
              </a-input-password>
            </a-form-item>
            
            <div v-if="loginConfig?.isForgetPasswordEnabled" class="forgot-password">
              <a href="#" @click.prevent="handleForgotPassword">{{ t('account.loginPage.forgotPassword') }}</a>
            </div>
            
            <a-form-item>
              <a-button type="primary" html-type="submit" class="login-btn" block size="large">{{ t('account.loginPage.signIn') }}</a-button>
            </a-form-item>
          </a-form>
        </div>
        
        <!-- 注册链接 -->
        <div v-if="loginConfig?.isRegisterEnabled" class="signup">
          <p>{{ t('account.loginPage.newToOdoc') }} <a href="#" @click.prevent="handleSignUp">{{ t('account.loginPage.createAccount') }}</a></p>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { reactive, ref, onMounted } from 'vue';
  import { useRouter, useRoute } from 'vue-router';
  import { message } from 'ant-design-vue';
  import { MailOutlined, LockOutlined } from '@ant-design/icons-vue';
  import { useUserStore } from '@common/stores/user';
  import { getDomainOrigin } from '~/src/util/env';
  import { getEncryptKey, signIn } from '~/src/util/login';
  import { initTheme } from '@common/theme';
  import { useI18n } from 'vue-i18n';
  import { getLoginPageConfig } from '~/src/api/membership';
  import type { GetLoginPageConfigResponse } from 'go-sea-proto/gen/ts/membership/MembershipApi';

  // 确保主题正确初始化
  onMounted(() => {
    if (typeof window !== 'undefined') {
      // 初始化主题
      initTheme();
      console.log('登录页面主题已初始化');
    }
  });
  
  const router = useRouter();
  const route = useRoute();
  const userStore = useUserStore();
  const { t } = useI18n();
  
  const formState = reactive({
    email: '',
    password: ''
  });
  
  // 登录页面配置
  const loginConfig = ref<GetLoginPageConfigResponse | null>(null);
  
  // 如果已登录，跳转到来源页面或首页
  const checkLoginStatus = async () => {
    if (userStore.isLogin()) {
      const redirect = route.query.redirect as string || '/workbench';
      router.replace(redirect);
    }
  };

  // 使用 onMounted 确保只在组件首次挂载时检查一次
  onMounted(() => {
    checkLoginStatus();
    getConfig();
  });

  const getConfig = async () => {
    try {
      const res = await getLoginPageConfig();
      if (res) {
        loginConfig.value = res;
        console.log('登录页面配置:', res);
      } else {
        // 如果获取配置失败，使用默认配置
        loginConfig.value = {
          isGoogleLoginEnabled: true,
          isUsernamePasswordLoginEnabled: true,
          isRegisterEnabled: true,
          isForgetPasswordEnabled: true
        };
        console.warn('获取登录配置失败，使用默认配置');
      }
    } catch (error) {
      console.error('获取登录配置出错:', error);
      // 出错时使用默认配置
      loginConfig.value = {
        isGoogleLoginEnabled: true,
        isUsernamePasswordLoginEnabled: true,
        isRegisterEnabled: true,
        isForgetPasswordEnabled: true
      };
    }
  };
  
  // 处理表单提交
  const handleSubmit = async () => {
      const res = await signIn(formState.email, formState.password);
      if (!res) {
        message.error(t('account.loginPage.loginFailed'));
        return;
      }

      // 登录成功后，获取用户信息并更新到 store 中
      await userStore.getUserInfo();
      
      if (!userStore.isLogin()) {
        message.error('login failed, please try again');
        return;
      }

      //跳转到来的页面，如果有的话，否则跳转到工作台页面
      const redirect = route.query.redirect as string || '/workbench';
      
      // 使用 router.replace 进行跳转
      router.replace(redirect);
  };
  
  // 社交登录处理
  const handleGoogleLogin = () => {
    // 注释掉 uni-login 重定向
    // window.location.href = `${getDomainOrigin()}/apicdn/readpaper/uni-login/?type=google&redirect=${encodeURIComponent(window.location.href)}`;
    // message.info('Google 登录功能暂时不可用');
    window.location.href = '/api/oauth2/google/login';
  };
  
  const handleGithubLogin = () => {
    // 注释掉 uni-login 重定向
    // window.location.href = `${getDomainOrigin()}/apicdn/readpaper/uni-login/?type=github&redirect=${encodeURIComponent(window.location.href)}`;
    message.info('GitHub 登录功能暂时不可用');
  };
  
  const handleSSOLogin = () => {
    // 注释掉 uni-login 重定向
    // window.location.href = `${getDomainOrigin()}/apicdn/readpaper/uni-login/?type=sso&redirect=${encodeURIComponent(window.location.href)}`;
    message.info('SSO 登录功能暂时不可用');
  };
  
  // 其他功能
  const handleForgotPassword = () => {
    // 注释掉 uni-login 重定向
    // window.location.href = `${getDomainOrigin()}/apicdn/readpaper/uni-login/?action=reset&redirect=${encodeURIComponent(window.location.href)}`;
    message.info('密码重置功能暂时不可用');
  };
  
  const handleSignUp = () => {
    // 注释掉 uni-login 重定向
    // window.location.href = `${getDomainOrigin()}/apicdn/readpaper/uni-login/?action=signup&redirect=${encodeURIComponent(window.location.href)}`;
    message.info('注册功能暂时不可用');
  };
  </script>
  
  <style lang="less" scoped>
  /* 导入全局主题变量 */
  @import '../../assets/less/theme.less';
  /* Material Icons will be imported in the main CSS or replaced with local icons */
  /* Removed Google Fonts import to fix build error */

  /* 登录页面容器 */
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background: var(--site-theme-bg-primary);
    color: var(--site-theme-text-primary);
    transition: background-color 0.3s ease, color 0.3s ease;
  }

  .login-box {
    width: 420px;
    padding: 40px;
    background-color: var(--site-theme-bg-secondary);
    border-radius: 12px;
    box-shadow: var(--site-theme-shadow-3);
    border: 1px solid var(--site-theme-divider);
    transition: background-color 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease;
  }

  .logo {
    display: flex;
    justify-content: center;
    margin-bottom: 30px;
    
    .logo-text {
      font-size: 32px;
      font-weight: 700;
      background: linear-gradient(90deg, var(--site-theme-brand) 0%, var(--site-theme-brand-light) 100%);
      -webkit-background-clip: text;
      background-clip: text;
      -webkit-text-fill-color: transparent;
      letter-spacing: 1px;
    }
  }

  .title {
    text-align: center;
    font-size: 28px;
    margin-bottom: 24px;
    font-weight: 600;
    color: var(--site-theme-text-primary);
  }
  
  .social-login {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 24px;
    
    button {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 10px;
      padding: 12px 16px;
      border-radius: 6px;
      border: 1px solid var(--site-theme-divider);
      background-color: var(--site-theme-background-hover);
      color: var(--site-theme-text-primary);
      font-size: 15px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.3s;
      
      img {
        width: 20px;
        height: 20px;
      }
      
      &:hover {
        background-color: var(--site-theme-background-hover);
        border-color: var(--site-theme-brand);
        transform: translateY(-2px);
      }
    }
  }

  .username-password-login {
    margin-bottom: 24px;
  }
  
  .divider {
    display: flex;
    align-items: center;
    margin: 24px 0;
    
    &::before,
    &::after {
      content: '';
      flex: 1;
      height: 1px;
      background-color: var(--site-theme-divider);
    }
    
    span {
      padding: 0 16px;
      color: var(--site-theme-text-secondary-color);
      font-size: 14px;
      font-weight: 500;
    }
  }
  
  // 自定义 Ant Design 组件样式
  :deep(.ant-form) {
    .ant-form-item {
      margin-bottom: 20px;
      position: relative;
    }
    
    .ant-form-item-explain-error {
      color: var(--site-theme-danger-color);
      font-size: 12px;
      margin-top: 4px;
    }
    
    .ant-input,
    .ant-input-password {
      width: 100%;
      height: 48px;
      padding: 12px 16px 12px 45px;
      border-radius: 6px;
      border: 1px solid var(--site-theme-divider);
      background-color: var(--site-theme-bg-mute);
      color: var(--site-theme-text-primary);
      font-size: 15px;
      transition: all 0.3s;
      box-shadow: none;
      display: flex;
      align-items: center;
      
      &:focus,
      &-focused {
        outline: none;
        border-color: var(--site-theme-brand);
        box-shadow: 0 0 0 2px var(--site-theme-primary-color-fade);
      }
      
      &::placeholder {
        color: var(--site-theme-placeholder-color);
        line-height: 24px;
        vertical-align: middle;
      }
    }
    
    .ant-input-password {
      display: flex;
      align-items: center;
      padding-left: 45px;
      
      .ant-input {
        background-color: transparent;
        border: none;
        padding-left: 0;
        height: 24px;
        box-shadow: none;
        line-height: 24px;
      }
      
      .ant-input-affix-wrapper-status-error {
        background-color: transparent;
      }
    }
    
    .email-input,
    .password-input {
      .ant-input-prefix {
        position: absolute;
        left: 0;
        height: 100%;
        display: flex;
        align-items: center;
        padding-left: 16px;
      }
    }
    
    .email-input {
      border: 1px solid var(--site-theme-divider);
      border-radius: 6px;
      overflow: hidden;
      background-color: var(--site-theme-bg-mute) !important;
      height: 48px;
      padding-left: 45px;
      display: flex;
      align-items: center;
      
      .ant-input {
        border: none;
        box-shadow: none;
        background-color: transparent !important;
        padding-left: 0;
        height: 100%;
        line-height: 24px;
      }
      
      &.ant-input-affix-wrapper-focused,
      &:hover,
      &:focus,
      &.ant-input-affix-wrapper-status-error {
        background-color: var(--site-theme-bg-mute) !important;
      }
    }
    
    .password-input {
      .ant-input-password-icon {
        margin-right: 0px;
      }
    }
    
    .ant-input-password-icon {
      color: var(--site-theme-text-secondary-color);
    }
    
    .form-icon {
      color: var(--site-theme-text-secondary-color);
      font-size: 20px;
      margin-right: 8px;
      position: absolute;
      left: 16px;
      top: 50%;
      transform: translateY(-50%);
      z-index: 1;
    }
    
    .ant-btn {
      width: 100%;
      padding: 12px;
      height: auto;
      border-radius: 6px;
      border: none;
      background: linear-gradient(90deg, #6889ff 0%, #8662e9 50%, #e74694 100%);
      color: #fff;
      font-size: 16px;
      font-weight: 600;
      transition: all 0.3s;
      
      &:hover {
        background: linear-gradient(90deg, #7a97ff 0%, #9775f0 50%, #f05da6 100%);
        transform: translateY(-2px);
        box-shadow: var(--site-theme-shadow-2);
      }
    }
  }
  
  .forgot-password {
    text-align: right;
    margin-bottom: 20px;
    
    a {
      color: var(--site-theme-brand);
      font-size: 14px;
      text-decoration: none;
      transition: all 0.3s;
      
      &:hover {
        color: var(--site-theme-brand-light);
        text-decoration: underline;
      }
    }
  }
  
  .login-form {
    margin-bottom: 10px;
    
    :deep(.ant-form-item-control-input) {
      box-shadow: none;
      border: none;
    }
    
    :deep(.ant-form-item-has-error) {
      .ant-input, .ant-input-password {
        border-color: var(--site-theme-danger-color);
      }
    }
    
    /* 覆盖浏览器自动填充的背景色 */
    :deep(input:-webkit-autofill),
    :deep(input:-webkit-autofill:hover), 
    :deep(input:-webkit-autofill:focus),
    :deep(input:-webkit-autofill:active) {
      -webkit-box-shadow: 0 0 0 30px var(--site-theme-bg-mute) inset !important;
      -webkit-text-fill-color: var(--site-theme-text-color) !important;
      transition: background-color 5000s ease-in-out 0s;
    }
  }
  
  .login-btn {
    width: 100%;
    padding: 12px;
    border-radius: 6px;
    border: none;
    background: linear-gradient(90deg, #6889ff 0%, #8662e9 50%, #e74694 100%);
    color: var(--site-theme-text-inverse);
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
    
    &:hover {
      background: linear-gradient(90deg, #7a97ff 0%, #9775f0 50%, #f05da6 100%);
      transform: translateY(-2px);
      box-shadow: var(--site-theme-shadow-2);
    }
  }
  
  .signup {
    margin-top: 24px;
    text-align: center;
    font-size: 15px;
    color: var(--site-theme-text-secondary-color);
    
    a {
      color: var(--site-theme-brand);
      text-decoration: none;
      font-weight: 500;
      transition: all 0.3s;
      
      &:hover {
        color: var(--site-theme-brand-light);
        text-decoration: underline;
      }
    }
  }
  </style>