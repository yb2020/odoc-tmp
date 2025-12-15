export enum ImageMimeType {
  PNG = 'image/png',
  JPEG = 'image/jpeg',
  WEBP = 'image/webp',
}

export const compressImageBase64 = (
  canvas: HTMLCanvasElement,
  ratio = 1,
  mime = ImageMimeType.PNG
) => {
  let base64 = canvas.toDataURL(mime, ratio);

  if (
    mime === ImageMimeType.WEBP &&
    base64.indexOf(ImageMimeType.WEBP) === -1
  ) {
    base64 = canvas.toDataURL(ImageMimeType.JPEG, ratio);
  }

  return base64;
};

export class Defer<T> {
  public promise: Promise<T>;
  public resolve!: (value: T) => void;
  public reject!: (reason?: unknown) => void;

  public constructor() {
    this.promise = new Promise<T>((resolve, reject) => {
      this.resolve = resolve;
      this.reject = reject;
    });
  }
}
