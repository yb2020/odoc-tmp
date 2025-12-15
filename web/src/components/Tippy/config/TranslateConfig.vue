<template>
  <div class="config-containter">
    <div>
      <span class="config-title"
        >当前常用翻译渠道({{ modelList ? modelList.length : 0 }})</span
      >
      <br />
      <span class="config-desc">最多支持展示前5个</span>
    </div>
    <div
      v-for="(model, index) in modelList"
      :key="index"
      :class="current_idx === index ? 'config-item-selected' : 'config-item'"
    >
      <div class="config-content">
        <a-button
          :class="
            current_idx === index
              ? 'config-button-item-selected'
              : 'config-button-item'
          "
          style="text"
          @click="itemClick(index)"
          :disabled="model.type === 'other'"
          >{{ model.title }}
        </a-button>
        <div v-if="model.type !== 'other'">
          <a-button
            style="link"
            size="small"
            class="config-delete-button"
            @click="deleteAction(index)"
          >
            <template #icon>
              <DeleteOutlined />
            </template>
          </a-button>
        </div>
      </div>
    </div>
    <div class="config-dropdown-wrap" v-if="modelList.length < 5">
      <a-dropdown :trigger="['hover']" class="config-dropdown" :zIndex="1050">
        <template #overlay>
          <a-menu @click="menuChoose" style="background: white">
            <a-menu-item
              v-for="(model, index) in configList"
              :key="index"
              style="color: #999; font-size: 12px"
              @click="menuitemClick(index)"
            >
              <div class="config-menu-item">
                <span>{{ model.title }}</span>
              </div>
            </a-menu-item>
          </a-menu>
        </template>
        <a-button class="config-button">
          <template #icon>
            <PlusOutlined />
          </template>
          添加自定义翻译窗口
        </a-button>
      </a-dropdown>
    </div>
  </div>
</template>

<script lang="ts">
import {
  useTranslateStore,
  TranslateTabKey,
  GoogleTranslateStyle,
} from '~/src/stores/translateStore';
import {
  CustomTranslateChannel,
  DeleteInterfaceReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/CustomTranslateInterface';
import { get } from '@svgdotjs/svg.js';
import { ref, VNodeChild, defineComponent, onUnmounted } from 'vue';
import { DelegateInstance } from 'tippy.js';
import {
  DeleteOutlined,
  UserOutlined,
  DownOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue';
import {
  AddTxTranslateType,
  DeleteCustomerTranslateTab,
} from '@/api/translate';
import {
  emitter,
  CONFIG_TYPE,
  CONFIG_RESET_TYPE,
  CONFIG_ADD_TYPE,
  CONFIG_CLOSE_ACTION,
} from './config';
import { Modal } from 'ant-design-vue';
import { gteElectronVersion } from '~/src/util/env';

import { Config } from '@idea/aiknowledge-report';
interface MenuInfo {
  key: string;
  item: VNodeChild;
  domEvent: MouseEvent;
}

interface ConfigInfo {
  title: string;
  type: TranslateTabKey;
  verified: boolean;
  id: string;
}

export default defineComponent({
  components: {
    DeleteOutlined,
    UserOutlined,
    DownOutlined,
    PlusOutlined,
  },

  setup() {
    const store = useTranslateStore();
    const modelList = ref<ConfigInfo[]>([]);
    const current_idx = ref(-1);
    const showGoogleSetting = gteElectronVersion('1.20.6');
    const tencentModel = <ConfigInfo>{
      title: store.txConfig.name,
      type: TranslateTabKey.tencent,
      verified: false,
      id: '',
    };

    const aliModel = <ConfigInfo>{
      title: store.aliConfig.name,
      type: TranslateTabKey.ali,
      verified: false,
      id: '',
    };
    const googlModel = <ConfigInfo>{
      title: store.googleConfig.name,
      type: TranslateTabKey.google,
      verified: false,
      id: '',
    };

    const deeplModel = <ConfigInfo>{
      title: store.deeplConfig.name,
      type: TranslateTabKey.deepl,
      verified: false,
      id: '',
    };

    const modelArray = [
      aliModel,
      tencentModel,
      // googlModel,
      deeplModel,
    ];
    if (showGoogleSetting) {
      modelArray.push(googlModel);
    }

    const configList = ref(modelArray);

    const titles: { [k in TranslateTabKey]?: string } = {
      [TranslateTabKey.youdao]: '有道',
      [TranslateTabKey.baidu]: '百度',
      [TranslateTabKey.idea]: 'IDEA',
    };

    const addTypeListen = (e: any) => {
      const configInfo = JSON.parse(JSON.stringify(e));
      const addInfo = <ConfigInfo>{
        title: '',
        type: TranslateTabKey.other,
        verified: true,
        id: '',
      };
      if (configInfo.txSecretId?.length > 0) {
        /// 腾讯翻译
        addInfo.type = TranslateTabKey.tencent;
      }
      if (configInfo.aliAccessKeyId?.length > 0) {
        /// 阿里翻译
        addInfo.type = TranslateTabKey.ali;
      }
      if (
        configInfo.googleApiKey?.length > 0 ||
        store.googleConfigVersion !== GoogleTranslateStyle.none
      ) {
        /// 谷歌翻译
        addInfo.type = TranslateTabKey.google;
      }
      if (configInfo.deepLKey?.length > 0) {
        /// deepl
        addInfo.type = TranslateTabKey.deepl;
      }
      addInfo.id = configInfo.id;
      addInfo.title = configInfo.name;
      const item = modelList.value.findIndex(
        (item) => item.type === addInfo.type
      );

      if (item === -1) {
        modelList.value.push(addInfo);
      } else {
        modelList.value.splice(item, 1, addInfo);
      }
      current_idx.value = modelList.value.length - 1;
      const cusInfo = {
        name: addInfo.title,
        type: addInfo.type,
        verified: addInfo.verified,
        id: addInfo.id,
      };
      const idx = store.cusTabs.findIndex((item: any) => {
        return item.type === configInfo.type;
      });
      if (idx !== -1) {
        store.cusTabs.splice(idx, 1, cusInfo);
      } else {
        store.cusTabs.push(cusInfo);
      }
    };

    onUnmounted(() => {
      emitter.off(CONFIG_TYPE, addTypeListen);
    });

    emitter.on(CONFIG_ADD_TYPE, addTypeListen);

    const initConfigList = () => {
      store.tabs.map((item: any) => {
        if (titles[item as TranslateTabKey]) {
          const configInfo = <ConfigInfo>{
            title: titles[item as TranslateTabKey],
            type: TranslateTabKey.other,
            verified: true,
          };
          const configItem = modelList.value!.find(
            (obj) => obj.title === configInfo.title
          );
          if (!configItem) {
            modelList.value?.push(configInfo as ConfigInfo);
          }
        }
      });

      // modelList.value = modelList.value!.filter((item) => {
      //   return !titles[item.type as TranslateTabKey];
      // });

      for (let i = store.cusTabs.length - 1; i >= 0; i--) {
        const item = store.cusTabs[i];
        const configInfo = <ConfigInfo>{
          title: item.name as string,
          type: item.type as TranslateTabKey,
          verified: item.verified as boolean,
          id: item.id as string,
        };
        const idx = modelList.value.findIndex((item) => {
          return item.type === configInfo.type;
        });
        if (idx === -1 && item.verified) {
          modelList.value!.push(configInfo);
        }

        const filteredList = configList.value.filter(
          (item) => item.title !== configInfo.title
        );
        configList.value = filteredList;
      }
    };

    initConfigList();

    const menuitemClick = (index: number) => {
      if (modelList.value?.length >= 5) {
        return;
      }
      const value = configList.value[index];
      modelList.value!.push(value);
      current_idx.value = modelList.value.length - 1;
      emitter.emit(CONFIG_TYPE, configList.value[index]);

      configList.value.splice(index, 1);
    };

    const menuChoose = ({ key }: MenuInfo) => {};

    const itemClick = (index: number) => {
      current_idx.value = index;
      const chooseInfo = modelList.value![index];
      emitter.emit(CONFIG_TYPE, chooseInfo);
    };

    const qcConfigList = () => {
      const set = new Set();
      const qc: ConfigInfo[] = [];
      configList.value.forEach((item) => {
        if (!set.has(item.title)) {
          set.add(item.title);
          qc.push(item);
        }
      });
      configList.value = qc;
    };

    const deleteAction = async (index: number) => {
      const deleteInfo = JSON.parse(
        JSON.stringify(modelList.value![index])
      ) as ConfigInfo;
      Modal.confirm({
        title: '确定要删除翻译配置吗',
        zIndex: 1051,
        onOk: async () => {
          modelList.value!.splice(index, 1);
          current_idx.value = -1;
          emitter.emit(CONFIG_RESET_TYPE, deleteInfo);
          const cusInfo = {
            name: deleteInfo.title,
            type: deleteInfo.type,
            verified: false,
            id: '',
          };
          // 同步自定义接口数据
          const idx = store.cusTabs.findIndex((item: any) => {
            return item.type === deleteInfo.type;
          });

          if (idx !== -1) {
            store.cusTabs.splice(idx, 1);
          }
          /// 同步tabs数据
          const tabIdx = store.tabs.findIndex((item: any) => {
            return item.type === deleteInfo.type;
          });

          if (tabIdx !== -1) {
            store.tabs.splice(tabIdx, 1);
          }

          let deleteId = deleteInfo.id;

          deleteInfo.verified = false;
          deleteInfo.id = '';
          configList.value.push(deleteInfo);

          if (deleteId !== '') {
            const deleteReq = {
              id: deleteId.toString(),
            };
            await DeleteCustomerTranslateTab(deleteReq);
          }
        },
        onCancel: () => {},
      });
    };

    return {
      modelList,
      configList,
      store,
      titles,
      current_idx,
      menuChoose,
      itemClick,
      deleteAction,
      menuitemClick,
      initConfigList,
    };
  },
});
</script>

<style lang="less" scoped>
.config-containter {
  display: flex;
  flex-wrap: wrap;
  padding: 10px 0px;

  .config-dropdown-wrap {
    width: 100%;
    color: #999;

    .config-dropdown {
      width: 100%;
    }
  }

  .config-button {
    font-size: 12px;
    color: #4e5969;
    width: 150px;
    flex: 1;
    border: 0px;
    margin-top: 12px;
    padding-left: 16px;
    padding-right: 16px;
    background: #f0f2f5;
  }

  .config-item {
    width: 100%;
    height: 40px;
    border-bottom: 1px solid #e5e6eb;
  }

  .config-item-selected {
    width: 100%;
    height: 40px;
    background: #e8f5ff;
    border-bottom: 1px solid #e5e6eb;
  }

  .config-title {
    font-size: 14px;
    font-weight: bold;
    color: #262625;
    padding: 16px;
    font-family:
      PingFangSC-Regular,
      PingFang SC;
  }

  .config-desc {
    font-size: 12px;
    color: #86919c;
    padding-left: 16px;
    padding-bottom: 0px;
    font-family:
      PingFangSC-Regular,
      PingFang SC;
  }

  .config-content {
    display: flex;
    align-items: center;
    height: 40px;
    justify-content: space-between;
    padding: 0px;
  }

  .config-button-item {
    flex: 1;
    border: 0px;
    text-align: left;
    background: #f7f8fa;

    color: #262625;
  }

  .config-button-item-selected {
    flex: 1;
    border: 0px;
    text-align: left;
    background: #e8f5ff;

    color: #262625;
  }

  .config-delete-button {
    color: #86919c;
    border: none;
    background: none;
    padding: 0;
  }
}
</style>
