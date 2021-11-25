import request from './request'

export function login(data) {
  return request({
    url: '/api/v1/system/login',
    method: 'post',
    data
  })
}

export function logout(token) {
  return request({
    url: `/api/v1/system/logout?token=${token}`,
    method: 'get'
  })
}
