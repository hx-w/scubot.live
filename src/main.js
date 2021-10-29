import Vue from 'vue'
import JSONView from 'vue-json-viewer'

import App from './App.vue'
import router from './router' //引入路由配置
import axios from 'axios'
import './assets/index.css'
import './plugins/element'
import VueAmap from '@amap/amap-vue'
import Antd from 'ant-design-vue';
import { VueJsonp } from 'vue-jsonp';
import VideoPlayer from 'vue-video-player'
import 'video.js/dist/video-js.css'  // css 一定要引入
import 'ant-design-vue/dist/antd.css';

VueAmap.config.key = "dbf27129f0765f4a00dc2b1de0690568"
Vue.use(JSONView)
Vue.use(VueAmap)
Vue.use(Antd)
Vue.use(VueJsonp)
Vue.use(VideoPlayer)


Vue.prototype.$http = axios
Vue.config.productionTip = false


new Vue({
    render: h => h(App),
    router, //使用路由配置
}).$mount('#app')