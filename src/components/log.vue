<template>
  <div>
    <a-empty :description="cuid">在做了在做了...</a-empty>
  </div>
</template>
<script>
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
    },
    logs(newVal) {
      this.logs = newVal;
    }
  },
  methods: {},
};
</script>
