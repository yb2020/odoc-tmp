import Seo from '@oothkoo/seo-js';

const seo = new Seo({
  debug: false,
});

export const setSeo = (title: string) => {
  seo.use(seo.tag('title'));

  seo.use(seo.link('icon'));

  seo.update({
    title,
    icon: '/favicon.ico',
  });
};
