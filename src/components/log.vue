<template>
  <div>
    <a-empty :description="cuid">在做了在做了...</a-empty>
  </div>
</template>
<script>
import axios from "axios";
export default {
  props: ["uid"],
  data() {
    return {
      cuid: this.uid,
      logs: []
    };
  },
  watch: {
    uid(newVal) {
      this.cuid = newVal;
      axios
        .get("https://www.scubot.com/.netlify/functions/log", {
          params: {
            uid: newVal,
          },
        })
        .then((resp) => {
            this.logs = resp.data['message'];
            this.$message.success('日志获取成功');
            console.log(this.logs);
        })
        .catch((err) => {
            this.$message.error("日志获取失败");
            console.log(err);
        })
    },
  },
  methods: {},
};
</script>
