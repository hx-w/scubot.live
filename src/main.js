import Vue from 'vue'
import JSONView from 'vue-json-viewer'

import App from './App.vue'
import router from './router' //引入路由配置
import axios from 'axios'
import './assets/index.css'
import './plugins/element'
import VueAmap from 'vue-amap'

Vue.use(JSONView)
Vue.use(VueAmap)

VueAmap.initAMapApiLoader({
    key: "dbf27129f0765f4a00dc2b1de0690568",
    plugin: ["AMap.Autocomplete", "AMap.Geocoder", "AMap.Geolocation"],
    v: "1.4.15",
    uiVersion: "1.1"
});

Vue.prototype.$http = axios
Vue.config.productionTip = false


new Vue({
    render: h => h(App),
    router, //使用路由配置
}).$mount('#app')