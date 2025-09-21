import { RouteRecordRaw } from 'vue-router'

import UserLayout from '@/views/layout/UserLayout.vue'

export const routeConfigs: Array<RouteRecordRaw> = [
  {
    path: '/',
    component: UserLayout,
    children: [
      {
        path: 'list',
        name: 'DnsList',
        component: () => import('../views/dns/DnsList.vue'),
        meta: {
          title: 'DNS列表'
        }
      }
    ]
  },
  {
    path: '/404',
    name: 'PageError404',
    component: () => import('../views/error/PageError404.vue'),
    meta: {
      title: '404 找不到页面'
    }
  },
  {
    path: '/403',
    name: 'PageError403',
    component: () => import('../views/error/PageError403.vue'),
    meta: {
      title: '403 无权限访问'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: {
      name: 'PageError404'
    }
  }
]
