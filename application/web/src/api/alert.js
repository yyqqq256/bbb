import request from '@/utils/request'

// 异常检测
export function detectAnomalies(data) {
  return request({
    url: '/detectAnomalies',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 创建异常报警
export function createAlert(data) {
  return request({
    url: '/createAlert',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 创建召回记录
export function createRecall(data) {
  return request({
    url: '/createRecall',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 更新报警状态
export function updateAlertStatus(data) {
  return request({
    url: '/updateAlertStatus',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 更新召回状态
export function updateRecallStatus(data) {
  return request({
    url: '/updateRecallStatus',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 获取产品的报警记录
export function getFruitAlerts(data) {
  return request({
    url: '/getFruitAlerts',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 获取产品的召回记录
export function getFruitRecalls(data) {
  return request({
    url: '/getFruitRecalls',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 获取所有待处理报警
export function getPendingAlerts(data) {
  return request({
    url: '/getPendingAlerts',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 获取所有进行中的召回
export function getActiveRecalls(data) {
  return request({
    url: '/getActiveRecalls',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}