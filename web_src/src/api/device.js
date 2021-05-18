import * as config from './config';
import request from '@/utils/request'

export const deviceList = (data) => request.get(config.deviceList,data);
export const channelList = (data) => request.get(config.channelList,data);

export const startStream = (data) => request.get(config.startStream,data);
export const stopStream = (data) => request.get(config.stopStream,data);