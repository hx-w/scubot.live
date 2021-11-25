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

          <el-form-item label="微服务Cookies">
            <el-input
              v-model="scu.cookies"
              placeholder="
使用Chrome或Edge登录https://wfw.scu.edu.cn/ncov/wap/default/index（健康每日报网页版）
F12，切换至Network选项卡，刷新页面
将index请求的Cookies复制到此处
"
              :rows="5"
              type="textarea"
              @focus="showCookiesNotify"
              @blur="closeCookiesNotify"
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
          <el-form-item label="每日打卡地点">
            <el-input v-model="scu.area" @click.native="onMapSelect">
              <i slot="prefix" class="el-input__icon el-icon-map-location"></i>
            </el-input>
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
            <el-dialog
              title="选择定位"
              :visible.sync="dialogMapVisible"
              center
              top="15px"
            >
              <scumap ref="scumap" />
              <span slot="footer" class="dialog-footer">
                <el-button @click="dialogMapVisible = false">取 消</el-button>
                <el-button type="primary" @click="confirmLoc">确 定</el-button>
              </span>
            </el-dialog>
            <el-dialog
              title="如何获取Cookies"
              :visible.sync="dialogVideoVisible"
              center
              top="10px"
              width="80%"
            >
              <!-- <cookies-video /> -->
              <bili-video
                ihtml="https://player.bilibili.com/player.html?aid=633770199&bvid=BV1Rb4y1h7TQ&cid=432741794&page=1&danmaku=0&high_quality=1"
              />
              <!-- <bili-video ihtml="https://player.bilibili.com/player.html?aid=90992146&cid=155380603&page=1&danmaku=0&high_quality=1" /> -->
            </el-dialog>
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
              expanded
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
import BiliVideo from "./bili-video.vue";
const ePattern = /eai-sess=(.+?)(?:;|$| )/g
const uPattern = /UUkey=(.+?)(?:;|$| )/g

function sleep(time) {
  return new Promise((resolve) => setTimeout(resolve, time));
}

export default {
  components: { scumap, BiliVideo },
  data() {
    return {
      scu: {
        cookies: "",
        userAgent: "你的浏览器UA",
        triggerTime: "00:10",
        qqid: "",
        area: "四川省 成都市 武侯区",
        uuid: "",
      },
      accessToken: "",
      preview: "请点击'预览'",
      dialogMapVisible: false,
      dialogVideoVisible: false,
      cacheLocation: {
        lat: 30.630839301216,
        lng: 104.079966362848,
      },
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
  },
  methods: {
    onPreview() {
      if (this.scu.cookies == null) {
        this.$message.error("<微服务Cookies>不能为空")
        return;
      }
      const eGot = ePattern.exec(this.scu.cookies)
      const uGot = uPattern.exec(this.scu.cookies)
      var cookies = {};
      if (eGot != null && eGot.length > 1 && uGot != null && uGot.length > 1) {
        cookies = {
          "eai-sess": eGot[1],
          "UUkey": uGot[1],
        };
        axios
          .get("https://www.scubot.live/.netlify/functions/uuid", {
            params: {
              'UUkey': cookies["UUkey"],
              'eai-sess': cookies["eai-sess"]
            }
          })
          .then((resp) => {
            console.log(resp.data)
            this.scu.uuid = JSON.parse(resp.data)["uid"];
            this.$message.success("cookies验证成功");

            this.preview = JSON.parse(JSON.stringify(this.scu));
            delete this.preview.area;
            this.preview.cookies = cookies;
            this.preview.location = this.cacheLocation;
          })
          .catch((error) => {
            console.log(error);
            this.$message.error("出现了一些错误");
          });
      } else {
        console.log(eGot, uGot);
        this.$message.error("<微服务Cookies>格式错误");
        return;
      }
    },
    onSubmit() {
      if (this.preview == "请点击'预览'") {
        this.$message.warning("请先点击预览，再点击提交");
        return;
      }
      var postData = JSON.parse(JSON.stringify(this.preview));
      postData.accessToken = this.accessToken;
      axios
        // .post("https://ci.scubot.live:12121/set_checkin", postData)
        .post("https://www.scubot.live/.netlify/functions/post", postData)
        .then((res) => {
          if (res.status == 200) {
            this.$message.success(res.data["message"]);
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
    onMapSelect() {
      this.dialogMapVisible = true;
    },
    confirmLoc() {
      this.dialogMapVisible = false;
      this.cacheLocation["lat"] = this.$refs.scumap.position[1];
      this.cacheLocation["lng"] = this.$refs.scumap.position[0];
      // get new area info
      this.$jsonp("https://api.map.baidu.com/reverse_geocoding/v3/", {
        output: "json",
        coordtype: "wgs84ll",
        ak: "0hYGiH3Ob5ZhV0eWzrGVXCD3bEdBCi6L",
        location: `${this.cacheLocation["lat"]},${this.cacheLocation["lng"]}`,
      })
        .then((response) => {
          if (response.status != 0) {
            console.log(response);
            this.$message.error("获取地理位置信息失败");
          } else {
            console.log(response.result);
            var addr = response.result.addressComponent;
            this.scu.area = `${addr.province} ${addr.city} ${addr.district}`;
            this.$message.success("地理位置获取：" + this.scu.area);
          }
        })
        .catch((error) => {
          console.log(error.response);
          this.$message.error("获取地理位置信息失败");
        });
    },
    showCookiesNotify() {
      this.cookies_not_inst = this.$notify({
        title: "如何获取cookies",
        dangerouslyUseHTMLString: true,
        duration: 0,
        onClick: this.cookiesNotifyClicked,
        showClose: false,
        message: "点击这里，可观看教程视频<i class='el-icon-video-play'/>",
      });
    },
    closeCookiesNotify() {
      this.cookies_not_inst.close();
    },
    cookiesNotifyClicked() {
      this.dialogVideoVisible = true;
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
  height: 100%;
}
</style>
