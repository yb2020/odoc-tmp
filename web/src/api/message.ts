import axios from 'axios';
import { SuccessResponse } from './type';

export const getGlogalMessage = async () => {
  return null;
  // const res = await axios.get<
  //   SuccessResponse<{
  //     show: boolean;
  //     message: {
  //       zh: string;
  //       en: string;
  //     };
  //     lsKey: string;
  //   }>
  // >(`https://nuxt.cdn.readpaper.com/config/global_msg.json?ts=${Date.now()}`);
  // return res.data.data;
};
