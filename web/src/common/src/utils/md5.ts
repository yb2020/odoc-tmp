import SparkMD5 from 'spark-md5';

export const calculateFileMD5 = (file: File): Promise<string> => {
  return new Promise((resolve) => {
    const blobSlice =
      File.prototype.slice ||
      (File.prototype as any).mozSlice ||
      (File.prototype as any).webkitSlice;
    const chunkSize = 2097152; // Read in chunks of 2MB
    const chunks = Math.ceil(file.size / chunkSize);
    let currentChunk = 0;
    const spark = new SparkMD5.ArrayBuffer();
    const fileReader = new FileReader();

    fileReader.onload = function (e) {
      console.log('read chunk nr', currentChunk + 1, 'of', chunks);
      spark.append(e.target!.result as ArrayBuffer); // Append array buffer
      currentChunk++;

      if (currentChunk < chunks) {
        loadNext();
      } else {
        console.log('finished loading');
        resolve(spark.end());
      }
    };

    fileReader.onerror = function () {
      console.warn('oops, something went wrong.');
      resolve('');
    };

    function loadNext() {
      const start = currentChunk * chunkSize;
      const end =
        start + chunkSize >= file.size ? file.size : start + chunkSize;

      fileReader.readAsArrayBuffer(blobSlice.call(file, start, end));
    }

    loadNext();
  });
};
