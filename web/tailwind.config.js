/** @type {import('tailwindcss').Config} */

export default {
  darkMode: 'media', // or 'media' or 'class'
  content: [
    `./index.html`,
    `./src/**/*.{vue,js,ts,jsx,tsx}`,
    '../readpaper-ai/common/src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        'rp-blue-1': '#E8F5FF',
        'rp-blue-6': '#1F71E0',
        'rp-neutral-8': '#4e5969',
        'rp-neutral-5': '#A8AFBA',
        'rp-neutral-4': '#C9CDD4',
        'rp-neutral-3': '#e5e6eb',
        'rp-neutral-2': '#f0f2f5',
        'rp-neutral-1': '#F7F8FA',
        'rp-neutral-10': '#1D2229',
        'rp-neutral-6': '#86919C',
        'rp-62738c': '#62738c',
        'rp-red-6': '#E66045',
        'rp-dark-1': '#222426',
        'rp-dark-2': '#2F3337',
        'rp-dark-4': '#414548',
        'rp-dark-5': '#4D5256',
        'rp-dark-8': '#5B6167',
        'rp-white-6': '#FFFFFF73',
        'rp-white-8': '#FFFFFFA6',
        'rp-white-10': '#FFFFFFD9',
        'rp-darkblue-7': '#4387D9',
        'rp-grass-8': '#4A8C03',
        'rp-green-6': '#52C41A',
      },
      boxShadow: {
        'rp-icon':
          '2px 3px 6px 3px rgba(0, 0, 0, 0.08), 1px 1px 3px rgba(0, 0, 0, 0.12)',
      },
    },
    container: {},
  },
  plugins: [],
};
