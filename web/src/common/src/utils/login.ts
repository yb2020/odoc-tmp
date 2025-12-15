/*
 * Created Date: June 10th 2021, 5:24:07 pm
 * Author: zhoupengcheng
 * -----
 * Last Modified: June 23rd 2021, 4:52:01 pm
 */
import { Base64 } from 'js-base64';
import { getDomainOrigin } from './env';

export const gotoLogin = () => {
  window.location.href = `${getDomainOrigin()}/login?redirect_url=${encodeURIComponent(
    window.location.href
  )}`;
};

export function setAuthorization(clientId: string, clientSecret: string) {
  return 'Basic ' + Base64.encode(clientId + ':' + clientSecret);
}
