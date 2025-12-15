<template>
  <div class="text-sm flex justify-center items-center leading-22px">
    <span class="mr-1">{{ $t('common.text.support') }}{{ $t('common.symbol.colon') }}</span>
    <ul class="flex items-center p-0 m-0 list-none">
      <li>
        <AlipayOutlined
          class="mr-1"
          :style="{ color: '#009FE8' }"
        />{{
          $t('common.text.alipay')
        }}
      </li>
      <li class="mx-1">
        /
      </li>
      <li>
        <WechatFilled
          class="mr-1"
          :style="{ color: '#2CBC1F' }"
        />{{
          $t('common.text.wechat')
        }}
      </li>
    </ul>
  </div>
</template>

<script lang="ts" setup>
import { AlipayOutlined, WechatFilled } from '@ant-design/icons-vue';
</script>
