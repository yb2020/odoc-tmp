// import '~common/assets/vendors/particles.min.js'
import json from './particles.json';

declare global {
  interface Window {
    particlesJS?: any;
  }
}

export const initWelcomeAnimate = (domId: string) => {
  window.particlesJS?.(domId, json);
};
