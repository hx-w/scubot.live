<template>
  <div id="app">
    <!-- 路由的出入口，路由的内容将被显示在这里 -->
    <vue-particles
      color="#dedede"
      style="position: fixed"
      :linesDistance="150"
      :hoverEffect="false"
      :clickEffect="false"
      :particleOpacity="0.7"
      :particlesNumber="80"
    >
    </vue-particles>
    <router-view> </router-view>
  </div>
</template>
 
<script>
export default {
  name: "App",
};
</script>

<style>
</style>
