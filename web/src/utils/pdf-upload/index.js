import { v4 as uuidv4 } from 'uuid';
import SparkMD5 from 'spark-md5';
import { getUserDocCreateStatus, getUploadToken,getS3UploadToken, handleFileFastUpload } from '@/api/document';
// 导入计算文件SHA256的js库
import CryptoJS from 'crypto-js';
// 导入PDF.js
import * as pdfjsLib from '@idea/pdfjs-dist';

import {
  UserDocParsedStatusEnum,
} from 'go-sea-proto/gen/ts/doc/UserDocParsedStatus';
import { isErrorStatus, getStatusProgress } from './statusMapper.js';

// 初始化PDF.js worker
pdfjsLib.GlobalWorkerOptions.workerSrc = new URL('@idea/pdfjs-dist/build/pdf.worker.min.js', import.meta.url);

// 检查Web Crypto API是否可用
const isCryptoSubtleAvailable = () => {
  return window.crypto && window.crypto.subtle && typeof window.crypto.subtle.digest === 'function';
};

// 导出新的状态枚举供外部使用
export { UserDocParsedStatusEnum };

// 导出状态判断函数
export { isErrorStatus };

// 创建上传项
export const createItem = (file, extra) => {
  return {
    id: uuidv4(),
    file,
    status: UserDocParsedStatusEnum.READY,
    progress: 0,
    extra
  };
};

// 开始轮询文档创建状态
const startPollingStatus = async (item, token, onUpdate) => {
  try {
    // 轮询间隔
    const pollInterval = 2000; // 轮询间隔2秒
    
    // 更新状态的辅助函数
    const updateStatus = (status, progress) => {
      item.status = status;
      item.progress = progress;
      onUpdate();
    };
    
    // 继续轮询的辅助函数
    const continuePolling = () => {
      setTimeout(pollStatus, pollInterval);
    };
    
    const pollStatus = async () => {
      try {
        const statusResponse = await getUserDocCreateStatus({ token });
        console.log('statusResponse', statusResponse);
        
        // 根据接口返回的数据结构，使用data.status作为状态码
        const createStatus = statusResponse.data?.status; 
        
        //如果状态为空 则不更新状态，继续轮询
        if (createStatus === 0) {
          continuePolling();
          return;
        }
        
        // 如果状态没有变化，继续轮询
        if (createStatus === item.status) {
          continuePolling();
          return;
        }
        
        // 直接使用API返回的状态枚举值
        if (createStatus >= UserDocParsedStatusEnum.HEADER_DATA_PARSED) {
          updateStatus(createStatus, getStatusProgress(createStatus));
          // 确保将文档信息设置到 docInfo 属性中
          if (statusResponse.data && statusResponse.data.docInfo) {
            item.docInfo = statusResponse.data.docInfo;
          }
          
          // 添加延时让用户看到最终状态
          await new Promise((resolve) => setTimeout(resolve, 2000));
          
          // 触发上传完成事件
          if (typeof window !== 'undefined') {
            window.dispatchEvent(new CustomEvent('uploadFinished'));
          }
          // 状态已完成，停止轮询
          return;
        } else if (isErrorStatus(createStatus)) {
          // 处理错误状态
          updateStatus(createStatus, getStatusProgress(createStatus));
          return;
        } else {
          // 更新到中间状态，添加延时让用户看到
          updateStatus(createStatus, getStatusProgress(createStatus));
          
          // 添加延时让用户看到中间状态变化
          await new Promise((resolve) => setTimeout(resolve, 2000));
          
          // 继续轮询其他状态
          continuePolling();
        }
      } catch (error) {
        console.error("轮询过程中发生错误", error);
        // 轮询过程中的错误不应该导致上传状态变为错误
        // 只记录错误并继续轮询
        continuePolling();
      }
    };
    
    pollStatus();
  } catch (error) {
    // 即使轮询失败，也不中断上传流程
  }
};

// 使用 PDF.js 获取 PDF 页数
const getPdfPageCount = async (file) => {
  return new Promise((resolve, reject) => {
    processPdfFile(file, resolve, reject);
  });
};

// 处理 PDF 文件并获取页数
const processPdfFile = (file, resolve, reject) => {
  const fileReader = new FileReader();
  fileReader.onload = function() {
    const typedarray = new Uint8Array(this.result);
    
    pdfjsLib.getDocument(typedarray).promise.then(function(pdf) {
      resolve(pdf.numPages);
    }).catch(function(error) {
      console.error('Error getting PDF page count:', error);
      resolve(0); // 出错时返回默认值
    });
  };
  
  fileReader.onerror = function(error) {
    console.error('Error reading file:', error);
    resolve(0); // 出错时返回默认值
  };
  
  fileReader.readAsArrayBuffer(file);
};

// 上传文件
export const uploadItem = async (item, options, fileSHA256, onUpdate, tokenResponse) => {
  try {
    // 更新状态
    item.status = UserDocParsedStatusEnum.UPLOADING;
    item.progress = getStatusProgress(UserDocParsedStatusEnum.UPLOADING);
    onUpdate();

    //这里增加睡眠时间五秒
    // await new Promise((resolve) => setTimeout(resolve, 5000));

    //需要对其他状态的错误码进行处理
    if (!tokenResponse.data.needUpload && !tokenResponse.data.needParsed) {
      try {
        const handleFileFastUploadResponse = await handleFileFastUpload(tokenResponse.data);
        // 检查秒传接口的返回结果
        if (handleFileFastUploadResponse.data && handleFileFastUploadResponse.status === 1) {
          // 秒传成功后，更新状态为完成
          item.status = UserDocParsedStatusEnum.HEADER_DATA_PARSED;
          item.progress = getStatusProgress(UserDocParsedStatusEnum.HEADER_DATA_PARSED);
          // 如果有文档信息，设置到 docInfo 属性中
          if (handleFileFastUploadResponse.data?.docInfo) {
            item.docInfo = handleFileFastUploadResponse.data.docInfo;
          }
          onUpdate();
          // 触发上传完成事件，通知应用刷新列表
          if (typeof window !== 'undefined') {
            window.dispatchEvent(new CustomEvent('uploadFinished'));
          }
        } else {
          // 秒传失败，更新状态为错误
          item.status = UserDocParsedStatusEnum.UPLOAD_FAILED;
          item.exception = new Error(handleFileFastUploadResponse.data?.message || 'upload failed');
          item.message = handleFileFastUploadResponse.data?.message || 'upload failed';
          item.progress = getStatusProgress(UserDocParsedStatusEnum.UPLOAD_FAILED);
          onUpdate();
        }
      } catch (error) {
        // 处理秒传过程中的异常
        item.status = UserDocParsedStatusEnum.UPLOAD_FAILED;
        item.exception = error;
        item.message = error.message || 'upload failed';
        item.progress = getStatusProgress(UserDocParsedStatusEnum.UPLOAD_FAILED);
        onUpdate();
      }
      
      return;
    }
    
    // 处理needUpload为false但needParsed为true的情况
    if (!tokenResponse.data.needUpload && tokenResponse.data.needParsed) {
      // 直接更新状态为解析头部数据中
      item.status = UserDocParsedStatusEnum.PARSING_HEADER_DATA;
      item.progress = getStatusProgress(UserDocParsedStatusEnum.PARSING_HEADER_DATA);
      onUpdate();

      //这里增加睡眠时间五秒
      // await new Promise((resolve) => setTimeout(resolve, 5000));

      // 开始轮询文档创建状态
      startPollingStatus(item, tokenResponse.data.token, onUpdate);
      return; // 直接返回，不执行后续的上传逻辑
    }
    
    // 如果needUpload为true，使用返回的信息上传文件
    if (tokenResponse.data.uploadInfo && tokenResponse.data.uploadInfo.url && tokenResponse.data.uploadInfo.headers) {
      // 保存上传信息供后续使用
      item.uploadInfo = tokenResponse.data.uploadInfo;
    }
    // 检查是否有上传URL
    if (!item.uploadInfo || !item.uploadInfo.url) {
      console.error('[pdf-upload] 缺少上传URL');
      // 更新状态为错误
      item.status = UserDocParsedStatusEnum.UPLOAD_FAILED;
      item.message = '上传已取消';
      item.progress = getStatusProgress(UserDocParsedStatusEnum.UPLOAD_FAILED);
      onUpdate();
      return;
    }

    try {

      // 执行文件上传到预签名URL
      await executeFileUploadToPutUrl(item.uploadInfo, item.file);

    } catch (error) {
      // 捕获在上传过程中发生的任何网络错误或服务器错误。
      console.error('An error occurred during the upload process:', error);
      
      // 更新 UI 状态，向用户显示上传失败。
      item.status = UserDocParsedStatusEnum.UPLOAD_FAILED;
      item.exception = error;
      item.message = error.message || '上传失败，请检查网络连接或联系管理员。';
      item.progress = getStatusProgress(UserDocParsedStatusEnum.UPLOAD_FAILED); // 无论成功失败，进度都应结束
      onUpdate();
      
      // 可以在这里停止后续操作，因为上传已经失败。
      return; 
    }

    // 文件上传成功后，先更新状态为上传完成
    item.status = UserDocParsedStatusEnum.UPLOADED;
    item.progress = getStatusProgress(UserDocParsedStatusEnum.UPLOADED);
    onUpdate();

    //这里增加睡眠时间五秒
    // await new Promise((resolve) => setTimeout(resolve, 5000));
  
    // 开始轮询文档创建状态
    startPollingStatus(item, tokenResponse.data.token, onUpdate);

  } catch (error) {
    console.error('[上传流程异常] 文件:', item.file?.name, error);
    item.status = UserDocParsedStatusEnum.UPLOAD_FAILED;
    item.exception = error;
    item.message = error.message || '上传失败';
    item.progress = getStatusProgress(UserDocParsedStatusEnum.UPLOAD_FAILED);
  } finally {
    onUpdate();
  }
};


// =============================================抽象的公共方法===========================================

/**
 * 使用 crypto-js 计算文件的 SHA-256 哈希值。
 * 这是一个备用函数，用于在 Web Crypto API 不可用或计算失败时调用。
 * 
 * @param {ArrayBuffer} fileBuffer - 从 FileReader 读取到的文件 ArrayBuffer 数据。
 * @returns {string} - 文件的 SHA-256 哈希值（16进制字符串）。
 * @throws {Error} 如果 CryptoJS 库未加载，则会抛出错误。
 */
const calculateWithCryptoJS = (fileBuffer) => {
  // 1. 确保 CryptoJS 库已成功导入
  if (typeof CryptoJS === 'undefined') {
    throw new Error('CryptoJS library is not loaded. Please check your import.');
  }
  
  // 2. 将 ArrayBuffer 转换为 CryptoJS 可识别的 WordArray 格式
  const wordArray = CryptoJS.lib.WordArray.create(fileBuffer);
  
  // 3. 使用 CryptoJS 计算 SHA-256 哈希
  const hash = CryptoJS.SHA256(wordArray);
  
  // 4. 将哈希结果转换为 16 进制字符串并返回
  return hash.toString(CryptoJS.enc.Hex);
};


/**
 * 计算文件的 SHA-256 哈希值，优先使用浏览器内置的 Web Crypto API，
 * 并在 API 不可用或计算失败时，自动降级到使用 crypto-js 库。
 * 
 * @param {File} file - 需要计算哈希值的 File 对象。
 * @returns {Promise<string>} - 一个解析为文件 SHA-256 哈希值（16进制字符串）的 Promise。
 *                             如果两种方法都失败，Promise 将被 reject。
 */
const calculateFileSHA256 = (file) => {
  return new Promise((resolve, reject) => {
    // 1. 创建 FileReader 实例以异步读取文件
    const reader = new FileReader();

    // 2. 定义文件读取成功后的回调函数
    reader.onload = async (event) => {
      try {
        // 获取文件的 ArrayBuffer 数据
        const fileBuffer = event.target.result;

        // 3. **优先尝试使用 Web Crypto API**
        if (window.crypto && window.crypto.subtle) {
          console.log("Attempting to calculate SHA-256 using Web Crypto API...");
          try {
            // a. 使用 crypto.subtle.digest 计算哈希
            const hashBuffer = await window.crypto.subtle.digest('SHA-256', fileBuffer);
            
            // b. 将哈希值从 ArrayBuffer 转换为 16 进制字符串
            const hashArray = Array.from(new Uint8Array(hashBuffer));
            const hashHex = hashArray
              .map(byte => byte.toString(16).padStart(2, '0'))
              .join('');
            
            console.log("Web Crypto API calculation successful.");
            // c. 计算成功，返回结果，函数提前结束
            return resolve(hashHex);
          } catch (cryptoError) {
            // d. 如果 Web Crypto API 计算失败，在控制台打印错误，然后继续执行下面的备用方案
            console.error("Web Crypto API failed, falling back to crypto-js.", cryptoError);
          }
        }
        
        // 4. **备用方案：使用 crypto-js**
        // (此代码块会在 Web Crypto API 不存在或执行失败时运行)
        console.warn("Using crypto-js as primary or fallback method.");
        const hashHex = calculateWithCryptoJS(fileBuffer);
        console.log("crypto-js calculation successful.");
        resolve(hashHex);

      } catch (error) {
        // 5. 最终错误处理 (例如 crypto-js 计算失败)
        console.error("A critical error occurred during hash calculation with fallback.", error);
        reject(error);
      }
    };

    // 6. 定义文件读取失败的回调函数
    reader.onerror = () => {
      const errorMessage = 'Error reading the file.';
      console.error(errorMessage);
      reject(new Error(errorMessage));
    };

    // 7. 启动文件读取过程
    reader.readAsArrayBuffer(file);
  });
};

/**
 * 公共方法：安全地计算文件SHA-256哈希值
 * @param {File} file - 需要计算哈希值的文件
 * @returns {Promise<string>} - 返回SHA-256哈希值，如果计算失败返回空字符串
 */
export const calculateFileHashSHA256 = async (file) => {
  let fileSHA256 = '';
  try {
    fileSHA256 = await calculateFileSHA256(file);
  } catch (sha256Error) {
    console.error('calculate sha256 error:', sha256Error);
  }
  return fileSHA256;
};
/**
 * 公共方法：执行文件上传到预签名URL
 * @param {Object} uploadInfo - 上传信息对象，包含 url, method, headers
 * @param {File} file - 要上传的文件对象
 * @returns {Promise<Response>} - 返回fetch响应对象
 */
export const executeFileUploadToPutUrl = async (uploadInfo, file) => {
  // 1. 从后端获取上传信息
  const { url, method, headers } = uploadInfo;

  // 2. 验证上传方法是否为 PUT
  if (method.toUpperCase() !== 'PUT') {
    throw new Error(`Unsupported upload method: ${method}. Expected 'PUT'.`);
  }
  // 确保状态更新被渲染
  // await new Promise(resolve => setTimeout(resolve, 100));
  
  // 3. 使用 fetch 发送 PUT 请求。
  const response = await fetch(url, {
    method: 'PUT',
    // 4. 设置从后端获取的必需请求头。
    //    对于PUT上传，这通常只包含 'Content-Type'。
    headers: headers,
    // 5. 将文件本身作为请求体。
    body: file,
  });

  // 6. 检查响应状态。
  //    S3/MinIO 在 PUT 上传成功后会返回 200 OK。
  if (!response.ok) {
    // 如果上传失败，尝试从响应体中读取详细的XML错误信息。
    const errorText = await response.text();
    console.error('Upload failed with status:', response.status, 'Response Body:', errorText);
    throw new Error(`Upload failed. Server responded with status ${response.status}.`);
  }

  return response;
};

/**
 * 公共方法：从base64字符串完成文件上传流程
 * @param {string} base64String - base64字符串，可以包含或不包含data:前缀
 * @param {string} fileExtension - 文件扩展名，如 'pdf', 'png', 'jpg' 等（不需要包含点号）
 * @param {number} bucketEnum - 存储桶枚举值  参考OSSBucketEnum  路径go-sea-proto/gen/ts/oss/OSS
 * @param {number} keyPolicy - 密钥策略   参考OSSKeyPolicyEnum  路径go-sea-proto/gen/ts/oss/OSS
 * @returns {Promise<Object>} - 返回uploadInfo对象，上传失败时抛出错误
 */
export const uploadBase64File = async (base64String, fileExtension, bucketEnum, keyPolicy) => {
  try {
    // 1. 通过调用base64ToFile转成file对象
    const file = base64ToFile(base64String, fileExtension);
    
    // 2. 获取文件大小和计算SHA256
    const fileSize = file.size;
    const fileSHA256 = await calculateFileHashSHA256(file);
    
    // 3. 调用getUploadToken接口获取上传信息
    const uploadTokenReq = {
      fileName: file.name,
      fileSHA256: fileSHA256,
      folderId: 0, // 默认文件夹ID
      fileSize: fileSize,
      bucketEnum: bucketEnum,
      keyPolicy: keyPolicy,
    };
    const tokenResponse = await getS3UploadToken(uploadTokenReq);
    
    // 检查API响应状态
    if (tokenResponse.status !== 1) {
      throw new Error(tokenResponse.message || 'Failed to get upload token');
    }
    
    // 检查是否需要上传
    if (!tokenResponse.data || !tokenResponse.data.needUpload) {
      throw new Error('Upload not required or invalid response data');
    }
    
    // 获取uploadInfo
    const uploadInfo = tokenResponse.data.uploadInfo;
    if (!uploadInfo) {
      throw new Error('Missing uploadInfo in response');
    }
    
    // 4. 调用executeFileUploadToPutUrl执行上传
    await executeFileUploadToPutUrl(uploadInfo, file);
    
    // 上传成功，返回uploadInfo
    return tokenResponse.data;
    
  } catch (error) {
    // 上传失败，抛出错误
    console.error('Upload base64 file failed:', error);
    throw error;
  }
};

/**
 * 公共方法：从base64字符串完成文件上传流程
 * @param {File} file - 文件对象
 * @param {number} bucketEnum - 存储桶枚举值  参考OSSBucketEnum  路径go-sea-proto/gen/ts/oss/OSS
 * @param {number} keyPolicy - 密钥策略   参考OSSKeyPolicyEnum  路径go-sea-proto/gen/ts/oss/OSS
 * @returns {Promise<Object>} - 返回uploadInfo对象，上传失败时抛出错误
 */
export const uploadBaseFile = async (file, bucketEnum, keyPolicy) => {
  try {
    // 1. 获取文件大小和计算SHA256
    const fileSize = file.size;
    const fileSHA256 = await calculateFileHashSHA256(file);
    
    // 2. 调用getUploadToken接口获取上传信息
    const uploadTokenReq = {
      fileName: file.name,
      fileSHA256: fileSHA256,
      folderId: 0, // 默认文件夹ID
      fileSize: fileSize,
      bucketEnum: bucketEnum,
      keyPolicy: keyPolicy,
    };
    const tokenResponse = await getS3UploadToken(uploadTokenReq);
    
    // 检查API响应状态
    if (tokenResponse.status !== 1) {
      throw new Error(tokenResponse.message || 'Failed to get upload token');
    }
    
    // 检查是否需要上传
    if (!tokenResponse.data || !tokenResponse.data.needUpload) {
      throw new Error('Upload not required or invalid response data');
    }
    
    // 获取uploadInfo
    const uploadInfo = tokenResponse.data.uploadInfo;
    if (!uploadInfo) {
      throw new Error('Missing uploadInfo in response');
    }
    
    // 3. 调用executeFileUploadToPutUrl执行上传
    await executeFileUploadToPutUrl(uploadInfo, file);
    
    // 上传成功，返回uploadInfo
    return tokenResponse.data;
    
  } catch (error) {
    // 上传失败，抛出错误
    console.error('Upload base64 file failed:', error);
    throw error;
  }
};

/**
 * 公共方法：将base64字符串转换为File对象
 * @param {string} base64String - base64字符串，可以包含或不包含data:前缀
 * @param {string} fileExtension - 文件扩展名，如 'pdf', 'png', 'jpg' 等（不需要包含点号）
 * @returns {File} - 返回File对象
 */
export const base64ToFile = (base64String, fileExtension) => {
  // 处理包含data:前缀的base64字符串
  let base64Data = base64String;
  let detectedMimeType = '';
  
  if (base64String.startsWith('data:')) {
    // 提取MIME类型和base64数据
    const matches = base64String.match(/^data:([^;]+);base64,(.+)$/);
    if (matches) {
      detectedMimeType = matches[1];
      base64Data = matches[2];
    } else {
      throw new Error('Invalid base64 data URL format');
    }
  }
  
  // 如果没有从base64中提取到MIME类型，根据文件扩展名推断
  if (!detectedMimeType && fileExtension) {
    // 确保扩展名是小写且不包含点号
    const extension = fileExtension.toLowerCase().replace(/^\./, '');
    const mimeTypeMap = {
      'pdf': 'application/pdf',
      'png': 'image/png',
      'jpg': 'image/jpeg',
      'jpeg': 'image/jpeg',
      'gif': 'image/gif',
      'txt': 'text/plain',
      'doc': 'application/msword',
      'docx': 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      'xls': 'application/vnd.ms-excel',
      'xlsx': 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    };
    detectedMimeType = mimeTypeMap[extension] || 'application/octet-stream';
  }
  
  // 如果仍然没有MIME类型，使用默认值
  if (!detectedMimeType) {
    detectedMimeType = 'application/octet-stream';
  }
  
  // 将base64转换为二进制数据
  try {
    const binaryString = atob(base64Data);
    const bytes = new Uint8Array(binaryString.length);
    
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }
    
    // 生成默认文件名
    const fileName = `file.${fileExtension || 'bin'}`;
    
    // 创建File对象
    return new File([bytes], fileName, { type: detectedMimeType });
  } catch (error) {
    throw new Error(`Failed to convert base64 to file: ${error.message}`);
  }
};

/**
 * Vue组合式API
 * @param {Object} options 
 */
export const useVueUpload = (options) => {
  const uploadList = options.ref([]);
  
  // 触发更新
  const triggerUpdate = () => {
    uploadList.value = [...uploadList.value];
  };
  
  // 添加上传
  // extra.localPath 可选，Tauri 环境下传入本地文件路径
  const addUpload = async (file, uploadOptions, extra) => {
    try {
      // 创建上传项
      const item = createItem(file, extra);
      
      // 计算文件SHA256
      const fileSHA256 = await calculateFileHashSHA256(file);
      
      // 获取PDF页数
      let pageCount = 0;
      if (file.type === 'application/pdf' || file.name.endsWith('.pdf')) {
        try {
          pageCount = await getPdfPageCount(file);
        } catch (pageError) {
          console.error('获取PDF页数出错:', pageError);
        }
      }
      
      // 打印本地文件路径（Tauri 环境）
      if (extra && extra.localPath) {
        console.log('[Tauri] 本地文件路径:', extra.localPath);
      }
      
      // 调用getUploadToken接口 TODO： 需要替换成协议
      let uploadTokenReq;
      if (extra && extra.localPath) {
        // Tauri 环境：只传递本地文件路径
        uploadTokenReq = {
          localFilePath: extra.localPath,
          fileName: file.name,
        };
      } else {
        // 浏览器环境：传递完整的文件信息
        uploadTokenReq = {
          fileName: file.name,
          fileSHA256: fileSHA256,
          folderId: uploadOptions && uploadOptions.folderId ? uploadOptions.folderId : 0,
          fileSize: file.size,
          filePage: pageCount,
        };
      }
      const tokenResponse = await getUploadToken(uploadTokenReq);
      
      // 检查是否是限制相关的响应
      if (tokenResponse.data && tokenResponse.data.limitHandled) {
        return null; // 不添加到uploadList，直接返回
      }
      
      if (tokenResponse.status == 0) {
        // API返回错误状态，不添加到uploadList
        console.error('获取上传Token失败:', tokenResponse.message);
        throw new Error(tokenResponse.message || 'Failed to get upload token');
      }
      
      // 接口调用成功后，才将文件添加到uploadList
      uploadList.value.push(item);
      triggerUpdate();
      
      // 开始上传（传递正确的参数）
      await uploadItem(item, uploadOptions, fileSHA256, triggerUpdate, tokenResponse);
      return item;
    } catch (error) {
      console.error('Add upload error:', error);
      throw error;
    }
  };
  
  // 清空上传列表
  const clearUploadList = () => {
    uploadList.value = [];
  };
  
  // 获取上传项
  const getUploadItem = (id) => {
    return uploadList.value.find(item => item.id === id);
  };
  
  // 获取上传项通过uid
  const getUploadItemByUid = (uid) => {
    return uploadList.value.find(item => item.extra.uid === uid);
  };
  
  options.onMounted(() => {
    console.log('useVueUpload mounted');
  });
  
  return {
    uploadList,
    addUpload,
    clearUploadList,
    getUploadItem,
    getUploadItemByUid
  };
};
