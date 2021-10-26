<template>
  <el-col :span="22">
    <el-container style="solid #eee">
      <el-header><h1 align="center">SCU健康每日报自动打卡工作流</h1></el-header>
      <br />
      <el-form ref="scu" :model="scu" label-width="180px" size="large">
        <el-form-item label="常用的浏览器UA">
          <el-input v-model="scu.userAgent"></el-input>
        </el-form-item>

        <el-form-item label="微服务网页源代码">
          <el-input
            v-model="scu.uuid"
            placeholder="
使用Chrome或Edge登录https://wfw.scu.edu.cn/ncov/wap/default/index（健康每日报网页版）
右键网页空白处，选择'显示网页源代码' (或CTRL+U)
将网页中所有源代码复制粘贴到这里即可
"
            :rows="5"
            type="textarea"
          ></el-input>
          <el-link
            type="primary"
            href="https://wfw.scu.edu.cn/ncov/wap/default/index"
            target="_blank"
            >健康每日报网页链接</el-link
          >
        </el-form-item>

        <el-form-item label="每日打卡时间">
          <el-time-select
            v-model="scu.triggerTime"
            :picker-options="{
              start: '00:10',
              step: '00:30',
              end: '13:30',
            }"
            placeholder="选择时间"
          >
          </el-time-select>
        </el-form-item>
        <el-form-item label="接受消息的QQ">
          <el-input
            v-model="scu.qqid"
            placeholder="填写自己的QQ用于接受打卡回执，留空为不接受回执。"
          ></el-input>
        </el-form-item>
        <el-form-item label="Access Token">
          <el-input v-model="accessToken" show-password placeholder="联系开发者获取"></el-input>
        </el-form-item>
        <br />
        <el-form-item align="left">
          <el-button type="warning" @click="onPreview">预览</el-button>
          <el-button type="primary" @click="onSubmit">提交</el-button>
        </el-form-item>
        <el-form-item>
          <json-viewer
            :value="preview"
            :expand-depth="5"
            copyable
            boxed
            sort
          ></json-viewer>
        </el-form-item>
      </el-form>
    </el-container>
  </el-col>
</template>


<script>
import axios from "axios";
var pattern = /"uid":"(.+?)"/g;

export default {
  data() {
    return {
      scu: {
        uuid: "",
        userAgent: "你的浏览器UA",
        triggerTime: "00:10",
        qqid: "",
      },
      accessToken: "",
      preview: "请点击'预览'",
    };
  },
  created() {
    this.scu.userAgent = navigator.userAgent;
    this.$notify.success({
      title: "提示",
      message:
        "如果你需要每日打卡消息回执，请主动添加QQ(机器人)：3583618673，验证消息填：7355608",
      duration: 0,
    });
  },
  methods: {
    onPreview() {
      if (this.scu.uuid == null) {
        this.$message.error("<微服务网页源代码>不能为空");
        return;
      }
      var got = pattern.exec(this.scu.uuid);
      var uuid = "";
      if (got != null && got.length > 1) {
        uuid = got[1];
      } else {
        this.$message.error("<微服务网页源代码>格式错误");
        return;
      }
      this.scu.uuid = uuid;
      this.preview = this.scu;
    },
    onSubmit() {
      if (isNaN(this.scu.uuid)) {
        this.$message.error("请先点击预览，再点击提交");
        return;
      }
      var postData = this.scu;
      postData["accessToken"] = this.accessToken;
      axios
        .post("https://ci.csgowiki.top:12121/set_checkin", postData)
        .then((res) => {
          if (res.status == 200) {
            this.$message({
              message: res.data["message"],
              type: "success",
            });
          } else if (res.status == 403) {
            this.$message.error("错误：" + res.data["detail"]);
          } else {
            this.$message.error("出现了一些错误，请检查表单是否填写正确");
          }
        })
        .catch((res) => {
          console.log(res)
          if (res.response.status == 403) {
            this.$message.error("错误：" + res.response.data["detail"]);
            return;
          }
          this.$message.error("出现了一些错误，请检查表单是否填写正确");
        });
    },
  },
};
</script>
<style scoped>
/* @import url(https://unpkg.com/element-ui/lib/theme-chalk/index.css); */
@import url(https://gist.githubusercontent.com/hx-w/9b06a0027584661278da9f49c886858f/raw/5edef0c6de6560e2672bc1bbee9b26b7809c5e11/index.css);
.el-container {
  padding: 20px;
  font-family: Helvetica;
  margin-top: 20px;
}
</style>
