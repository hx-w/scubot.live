import Vue from 'vue'
import JSONView from 'vue-json-viewer'

import App from './App.vue'
import router from './router' //引入路由配置
import axios from 'axios'

import './plugins/element'

Vue.use(JSONView)

Vue.prototype.$http=axios
Vue.config.productionTip=false


new Vue({
    render: h => h(App),
    router, //使用路由配置
}).$mount('#app')