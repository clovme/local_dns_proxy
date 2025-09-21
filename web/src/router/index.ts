import { createRouter, createWebHashHistory } from 'vue-router'
import { routeConfigs } from './config'
import NProgress from 'nprogress'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: routeConfigs
})

router.beforeEach((to, __, next) => {
  NProgress.start()
  if (!to.name) {
    next({ name: 'DnsList' }) // 没有路由名，跳到首页
  } else {
    next() // 有路由名，正常放行
  }
})

router.afterEach(() => {
  // 更新标题
  NProgress.done()
})

export default router
