<template>
  <div id="app">
    <div v-if="access">
      <vue-particles
        color="#dedede"
        style="position: fixed; width:100%; height:100%"
        :linesDistance="150"
        :hoverEffect="false"
        :clickEffect="false"
        :particleOpacity="0.7"
        :particlesNumber="50"
      >
      </vue-particles>
      <router-view> </router-view>
    </div>
    <div v-else>
      <auth v-on:token_valid="onTokenValid"/>
    </div>
  </div>
</template>
 
<script>
import auth from "./components/auth.vue";

export default {
  components: { auth },
  name: "App",
  data() {
    return {
      access: true,
    };
  },
  methods: {
    onTokenValid() {
      this.access = true
    }
  }
};
</script>

<style>
</style>
