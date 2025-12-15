const glob = require('glob');
const assetsFolder = `${process.cwd()}/dist/assets`;

const files = glob.sync(`${assetsFolder}/**/*.*`);

const start = async () => {
  const { cdn } = await import('@idea/aiknowledge-lib-qiniu-cdn');
  cdn(files, {
    checkBefore: true,
    fileKeyPrefix: 'pdf-annotate/2.0/assets'
  })
}

start()
