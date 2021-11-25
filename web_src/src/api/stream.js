import request from './request'

export function stream(token, id) {
  return request({
    url: `/api/v1/channels/${id}/stream?token=${token}`,
    method: 'get'
  })
}
