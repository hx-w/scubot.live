import Vue from 'vue' //引入 Vue
import VueRouter from 'vue-router' //引入 Vue 路由

Vue.use(VueRouter) //安装插件

const router = new VueRouter({
  mode: 'history',
  base: '/',
  routes: [
    { path: '/', name: 'scu', component: () => import('@/views/auto-checkin') },
  ]
})


export default router
