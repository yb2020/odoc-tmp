import isMobile from 'is-mobile';
import { getWebHost } from '../utils/env';

const getAgentMetaContentURL = (s?: string) => {
  const o = s?.split(';').reduce(
    (acc, cur) => {
      const [key, value] = cur.split('=');
      acc[key] = value;
      return acc;
    },
    {} as Record<string, string>
  );

  return o?.url.trim();
};

try {
  const isH5 = isMobile();
  let href;
  if (!isH5) {
    href =
      document.querySelector('link[rel="canonical"]')?.getAttribute('href') ||
      '';
  } else {
    href =
      document.querySelector('link[rel="alternate"]')?.getAttribute('href') ??
      getAgentMetaContentURL(
        document.querySelector<HTMLMetaElement>('meta[name="mobile-agent"]')
          ?.content
      );
  }

  if (href) {
    const parsed = new URL(href);

    if (window.location.search.slice(1)) {
      const query = new URLSearchParams(window.location.search);
      const params = new URLSearchParams({
        ...Object.fromEntries(parsed.searchParams),
        ...Object.fromEntries(query),
      });
      parsed.search = params.toString();
    }

    if (['dev', 'uat'].includes(import.meta.env.VITE_API_ENV)) {
      parsed.hostname = getWebHost();
    }

    window.location.replace(parsed.toString());
  }
} catch (e) {}
