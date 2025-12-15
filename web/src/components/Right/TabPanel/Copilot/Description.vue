<template>
  <div class="no-permission-container">
    <div class="copilot-description">
      <img class="ai-img" src="@/assets/images/right/ai.png" alt="" />
      <h2>Chat with Flow</h2>
      <p>Ask questions or request suggestions for your paper in general</p>
      <div class="text">
        We have made a summary for this paper and summarized several
        questions:
      </div>
    </div>
    <div class="go-to-chat">
      <a-button type="primary" @click="showChatPage">
        {{ $t('aiCopilot.startChat') }}
        <ArrowRightOutlined />
      </a-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ArrowRightOutlined } from '@ant-design/icons-vue'
import { PAGE_ROUTE_NAME } from '~/src/routes/type'
import { useClip } from '~/src/hooks/useHeaderScreenShot'
import { useCopilotStore } from '~/src/stores/copilotStore'

const props = withDefaults(
  defineProps<{
    clipSelecting: boolean
    clipAction: ReturnType<typeof useClip>['clipAction']
    page?: PAGE_ROUTE_NAME.CHAT | PAGE_ROUTE_NAME.NOTE | PAGE_ROUTE_NAME.WRITE
    textValue?: string
  }>(),
  {
    page: PAGE_ROUTE_NAME.NOTE,
    textValue: '',
  }
)

const copilotStore = useCopilotStore()

// Toggle between landing page and chat page
const showLandingPage = ref(!copilotStore.accessAiCopilot)

const showChatPage = () => {
  showLandingPage.value = false
}

</script>

<style lang="less" scoped>
.copilot-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #262626;
  color: #fff;
  position: relative;
}

.chat-container {
  height: 100%;
  position: relative;
}

.back-to-home {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 10;
}

.no-permission-container {
  .copilot-description {
    padding: 20px;
  }

  .ai-img {
    width: 80px;
  }

  .text {
    font-family: Noto Sans SC;
    font-size: 15px;
    font-weight: 400;
    line-height: 24px;
    color: #ffffffa6;
    margin: 13px 0 16px;
  }

  .ant-btn {
    min-width: 102px;
  }
}

.go-to-chat {
  margin-top: 20px;
  text-align: center;
}

.copilot {
  &-language {
    :deep(.ant-select-selector) {
      border: none !important;
    }
  }

  &-loading {
    text-align: center;
  }

  &-tip {
    color: #838c95;
    padding: 5px 16px;
    line-height: 22px;
    border-bottom: 1px solid #414548;
    margin-bottom: 16px;
    font-size: 13px;
    display: flex;
    justify-content: space-between;
    flex-wrap: wrap;

    .content {
      cursor: pointer;
    }

    .join-test-group {
      margin-left: 12px;
      color: #838c95;
    }
  }

  &-ps {
    flex: 1;
  }

  &-item {
    overflow: hidden;
    margin-bottom: 24px;

    .inbox {
      border-radius: 2px;
      max-width: 80%;
      padding: 10px 12px;
      clear: both;
      word-wrap: break-word;
    }

    &__simple {
      overflow: hidden;

      .question {
        float: right;
        background-color: #1f71e0;
        margin-bottom: 24px;
        margin-right: 12px;

        &.select {
          background-color: #404346;
          border-left: 4px solid #65676c;
          color: #abb2ba;
          padding: 0 8px;
          font-style: italic;
          -webkit-line-clamp: 10;
          display: -webkit-box;
          overflow: hidden;
          -webkit-box-orient: vertical;
        }

        .quote {
          border-left: 3px solid rgba(255, 255, 255, 0.15);
          color: rgba(255, 255, 255, 0.45);
          margin-bottom: 8px;
          padding-left: 8px;
          display: -webkit-box;
          overflow: hidden;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }
      }

      .question + .question {
        margin-top: -20px;
      }

      .answer {
        float: left;

        .error {
          color: #ffccc7;
        }
      }
    }
  }
}
</style>