import { requestAjax } from './http'

export interface DnsVO {
  id: number
  protocol: string
  domain: string
  ip: string
  status: string
  port: string
  updatedAt: string
  createdAt: string
}

export function getDnsListPage (params?: any) {
  return requestAjax({
    url: '/list',
    method: 'get',
    params
  })
}

export function getNetworkInterfaces () {
  return requestAjax({
    url: '/network/interfaces',
    method: 'get'
  })
}

export function getCopyright () {
  return requestAjax({
    url: '/copyright',
    method: 'get'
  })
}

export function postDnsSaveBatch (data?: any) {
  return requestAjax({
    url: '/save',
    method: 'post',
    data
  })
}

export function deleteDnsDelete (data?: any) {
  return requestAjax({
    url: '/delete',
    method: 'delete',
    data
  })
}

export function postDnsService (action: string, iface: string) {
  // @ts-ignore
  const first = window.isFirstLoading
  return requestAjax({
    url: `/service/${action}/${first}/${iface}`,
    method: 'post'
  })
}
