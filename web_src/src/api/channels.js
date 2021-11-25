import request from './request'

export function list(token, id) {
  return request({
    url: `/api/v1/devices/${id}/channels?token=${token}`,
    method: 'get'
  })
}
