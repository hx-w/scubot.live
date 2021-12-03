<template>
  <a-row type="flex" justify="center">
    <vue-particles
      class="wallpaper"
      color="#dedede"
      style="position: absolute; width: 100%; height: 100%"
      :linesDistance="100"
      :hoverEffect="false"
      :clickEffect="false"
      :particleOpacity="0.7"
      :particlesNumber="80"
    >
    </vue-particles>
    <a-input-password
      placeholder="输入一些应该输入的"
      style="width: 40%; margin-top: 300px"
      size="large"
      @pressEnter="onSubmit"
    />
  </a-row>
</template>

<script>
import axios from "axios";

export default {
  name: "auth",
  data() {
    return {
      token: "",
    };
  },
  methods: {
    onSubmit(event) {
      axios
        .get("https://www.scubot.com/.netlify/functions/token", {
          params: {
            "simple_verify": this.token,
          },
        })
        .then((resp) => {
          this.$emit("token_valid");
        })
        .catch((error) => {
          this.$message.error(
            "不对哦",
            3
          );
        });
    },
  },
};
</script>

<style>
.wallpaper {
  background-image: url("../assets/img/wallpaper.png");
  background-size: auto;
}
</style>