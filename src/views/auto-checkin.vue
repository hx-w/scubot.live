<template>
  <el-container direction="vertical">
    <el-row>
      <el-header>
        <h1 align="center">SCU健康每日报自动打卡工作流</h1>
        <div align="center">
          <img
            align="center"
            alt="runtime"
            hspace="5"
            src="https://img.shields.io/endpoint?url=https%3A%2F%2Fci.scubot.live%3A12121%2Fget_badge%3Ftype_%3D2"
          />
          <img
            align="center"
            alt="total"
            hspace="5"
            src="https://img.shields.io/endpoint?url=https%3A%2F%2Fci.scubot.live%3A12121%2Fget_badge%3Ftype_%3D1"
          />
          <img
            align="center"
            alt="valid workflow"
            hspace="5"
            src="https://img.shields.io/endpoint?url=https%3A%2F%2Fci.scubot.live%3A12121%2Fget_badge%3Ftype_%3D0"
          />
        </div>
        <br />
      </el-header>
    </el-row>
    <el-row justify="center" type="flex">
      <el-col :span="20">
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
            <el-input
              v-model="accessToken"
              show-password
              placeholder="联系开发者获取"
            ></el-input>
          </el-form-item>
          <div class="scumap" align="center">
          <el-card class="box-card" style="width: 85%; left:90px;position: relative;">
            <amap>
              <scumap />
            </amap>
          </el-card>
          </div>
          <br />
          <el-form-item>
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
      </el-col>
    </el-row>
  </el-container>
</template>


<script>
import axios from "axios";
import scumap from "./scumap.vue";
var pattern = /"uid":"(.+?)"/g;

function sleep(time) {
  return new Promise((resolve) => setTimeout(resolve, time));
}

export default {
  components: { scumap },
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
      dangerouslyUseHTMLString: true,
      message:
        "如果你需要每日打卡消息回执，请主动添加QQ(机器人)：<strong><a href='https://wpa.qq.com/msgrd?v=3&uin=3583618673&site=scubot&from=scubot&menu=yes'>3583618673</a></strong>，验证消息填：<strong><span id='verifyCode' onclick='copy();'>7355608</span></strong>",
      duration: 0,
    });
    this.$notify.info({
      title: "功能相关",
      message: "地图定位功能还未接入，目前打卡定位地点为[四川大学望江校区研究生院]",
      duration: 0,
      offset: 120,
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
      if (isNaN(this.scu.uuid) || this.scu.uuid.length == 0) {
        this.$message({ message: "请先点击预览，再点击提交", type: "warning" });
        return;
      }
      var postData = this.scu;
      postData["accessToken"] = this.accessToken;
      axios
        .post("https://ci.scubot.live:12121/set_checkin", postData)
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
          this.$message.error("出现了一些错误，请检查表单是否填写正确");
          console.log(res);
          if (res.response != undefined && res.response.status == 403) {
            sleep(500).then(() => {
              this.$message.error("错误：" + res.response.data["detail"]);
            });
          }
        });
    },
  },
};
</script>
<style scoped>
/* @import url(https://unpkg.com/element-ui/lib/theme-chalk/index.css); */
.el-container {
  padding: 10px;
  font-family: Helvetica;
  margin-top: 10px;
}

.el-form {
  margin-top: 50px;
}

.scumap {
  width: 100%;
  height: 20%;
}
</style>
