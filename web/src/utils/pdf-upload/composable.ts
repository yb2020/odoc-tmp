import { ref, Ref, onMounted } from 'vue';
import { RcFile } from 'ant-design-vue/lib/vc-upload/interface';
import { Item, UploadExtra, UploadOptions, VueUploadOptions } from './types';
import { createItem, uploadItem, cancelUpload, resolveConflict } from './core';

export const useVueUpload = (options: VueUploadOptions) => {
  const uploadList: Ref<Item[]> = options.ref([]);
  
  // 触发更新
  const triggerUpdate = () => {
    uploadList.value = [...uploadList.value];
  };
  
  // 添加上传
  const addUpload = async (file: RcFile, uploadOptions: UploadOptions, extra: UploadExtra) => {
    console.log('addUpload called - Timestamp:', new Date().toISOString());
    
    // 创建上传项
    const item = createItem(file, extra);
    uploadList.value.push(item);
    triggerUpdate();
    
    // 开始上传
    try {
      await uploadItem(item, uploadOptions, triggerUpdate);
      return item;
    } catch (error) {
      console.error('Add upload error:', error);
      throw error;
    }
  };
  
  // 取消上传
  const cancelUploadItem = async (item: Item) => {
    try {
      await cancelUpload(item, triggerUpdate);
      return item;
    } catch (error) {
      console.error('Cancel upload error:', error);
      throw error;
    }
  };
  
  // 解决冲突
  const resolveUploadConflict = async (item: Item, useOrigin: boolean) => {
    try {
      await resolveConflict(item, useOrigin, triggerUpdate);
      return item;
    } catch (error) {
      console.error('Resolve conflict error:', error);
      throw error;
    }
  };
  
  // 清空上传列表
  const clearUploadList = () => {
    uploadList.value = [];
  };
  
  // 获取上传项
  const getUploadItem = (id: string) => {
    return uploadList.value.find(item => item.id === id);
  };
  
  // 获取上传项通过uid
  const getUploadItemByUid = (uid: string) => {
    return uploadList.value.find(item => item.extra.uid === uid);
  };
  
  options.onMounted(() => {
    console.log('useVueUpload mounted');
  });
  
  return {
    uploadList,
    addUpload,
    cancelUpload: cancelUploadItem,
    resolveUploadConflict,
    clearUploadList,
    getUploadItem,
    getUploadItemByUid
  };
};
