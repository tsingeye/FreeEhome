import request from './request'

export function list(token) {
  return request({
    url: `/api/v1/devices?token=${token}`,
    method: 'get'
  })
}
